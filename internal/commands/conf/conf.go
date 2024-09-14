package conf

import (
	"fmt"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/jsumners/ghm/internal/app"
	"github.com/jsumners/ghm/internal/config"
	"github.com/spf13/cobra"
)

func New(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Work with CLI configuration.",
		Long: heredoc.Doc(`
			Provides commands to view the configuration,
			generate a basic configuration file, and view
			the loaded configuration.
		`),
	}

	dumpConfigCommand := &cobra.Command{
		Use:   "dump",
		Short: "Write found configuration to stdout.",
		Long: heredoc.Doc(`
			Write the configuration, as the application has read it
			from the configuration file and environment, to stdout.
		`),
		RunE: func(*cobra.Command, []string) error {
			return dumpConfig(app.Config)
		},
	}
	cmd.AddCommand(dumpConfigCommand)

	generateConfigCommand := &cobra.Command{
		Use:   "generate",
		Short: "Write default configuration to stdout.",
		RunE: func(*cobra.Command, []string) error {
			return generateConfig(app.Config)
		},
	}
	cmd.AddCommand(generateConfigCommand)

	return cmd
}

func dumpConfig(cfg *config.Config) error {
	yml, err := cfg.GenerateCurrentYaml()
	if err != nil {
		return err
	}
	fmt.Println(yml)
	return nil
}

func generateConfig(cfg *config.Config) error {
	yml, err := cfg.GenerateDefaultYaml()
	if err != nil {
		return err
	}
	fmt.Println(yml)
	return nil
}
