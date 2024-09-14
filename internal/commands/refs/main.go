package refs

import (
	"github.com/jsumners/ghm/internal/app"
	"github.com/jsumners/ghm/internal/commands/refs/recentreleases"
	"github.com/spf13/cobra"
)

func New(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refs",
		Short: "Discover and work directly with Git refs.",
	}

	cmd.AddCommand(recentreleases.New(app))

	return cmd
}
