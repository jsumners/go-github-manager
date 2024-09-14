package app

import (
	"context"
	"github.com/google/go-github/v63/github"
	"github.com/jsumners/ghm/internal/config"
	"log/slog"
)

type App struct {
	Config  *config.Config
	Client  *github.Client
	Context context.Context
	Logger  *slog.Logger

	// Owner is the GitHub user or organization name containing the repository,
	// or set of repositories, that will be targeted.
	Owner string
}
