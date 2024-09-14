package contributors

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/jsumners/ghm/internal/app"
	"github.com/jsumners/ghm/internal/commands/contributors/listall"
	"github.com/spf13/cobra"
)

func New(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contributors",
		Short: "Interrogate contributor information.",
		Long: heredoc.Doc(`
			Query organizations and repositories for contributor information.
		`),
	}

	cmd.AddCommand(listall.New(app))

	return cmd
}
