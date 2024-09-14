package root

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/jsumners/ghm/internal/app"
	"github.com/spf13/cobra"
)

type InitFn func() error

var longDesc = heredoc.Doc(`
	A tool for performing administrative and analysis tasks with GitHub
	organizations and repositories.
`)

func New(app *app.App, configFile *string, debugEnabled *bool, initFns ...InitFn) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ghm <command> <subcommand> [flags]",
		Short: "GitHub Management Tool",
		Long:  longDesc,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			for _, initFn := range initFns {
				err := initFn()
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(
		configFile,
		"config-file",
		"c",
		"",
		"Set the file from which configuration will be loaded.",
	)

	cmd.PersistentFlags().BoolVarP(
		debugEnabled,
		"verbose",
		"v",
		false,
		"Enable verbose/debug logging. All logs are written to stderr.",
	)

	return cmd
}
