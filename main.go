package main

import (
	"context"
	"github.com/google/go-github/v63/github"
	"github.com/jsumners/ghm/internal/app"
	"github.com/jsumners/ghm/internal/commands/conf"
	"github.com/jsumners/ghm/internal/commands/contributors"
	"github.com/jsumners/ghm/internal/commands/refs"
	"github.com/jsumners/ghm/internal/commands/root"
	"github.com/jsumners/ghm/internal/config"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

var cliApp *app.App
var cfgFile string
var debugEnabled bool

func main() {
	cliApp = &app.App{
		Config: config.New(),
	}

	cmd := root.New(
		cliApp,
		&cfgFile,
		&debugEnabled,
		initConfig, initLogger, createClient, setContext,
	)
	cmd.AddCommand(conf.New(cliApp))
	cmd.AddCommand(contributors.New(cliApp))
	cmd.AddCommand(refs.New(cliApp))

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() error {
	return cliApp.Config.ReadConfig(cfgFile)
}

func initLogger() error {
	var LEVEL slog.Leveler
	switch strings.ToLower(cliApp.Config.GetString("log_level")) {
	case "info":
		LEVEL = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: LEVEL,
	}))
	cliApp.Logger = logger

	return nil
}

func createClient() error {
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost:     5,
			MaxIdleConnsPerHost: 2,
		},
	}
	cliApp.Client = github.NewClient(httpClient)

	if cliApp.Config.AuthToken != "" {
		cliApp.Client = cliApp.Client.WithAuthToken(cliApp.Config.AuthToken)
	} else {
		cliApp.Logger.Warn("Authentication token not found. Expect to be rate limited.")
	}

	return nil
}

func setContext() error {
	cliApp.Context = context.WithValue(
		context.Background(),
		github.SleepUntilPrimaryRateLimitResetWhenRateLimited,
		true,
	)
	return nil
}
