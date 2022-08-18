package main

import (
	"context"

	"github.com/shipt/bubinga"
	"github.com/shipt/tempest-template/internal/webserver"
	"github.com/shipt/tempest-template/internal/worker"
	"github.com/shipt/tempest/app"
	"github.com/shipt/tempest/config"
	"github.com/spf13/cobra"
)

// cfg is the global app configuration store
// you can use it to unmarshal configuration values into
// go structs, e.g.:
//
// var webserverCfg webserver.Config
// if err := cfg.Unmarshal(&webserverCfg); err != nil {
//   return err
// }
var cfg = config.New()

var (
	// the *cobra.Command's below define the service's command line interface
	//
	// adding new commands to the CLI is as simple as adding a new *cobra.Command
	// object below and associating it with the rootCmd in the main() function.
	// for more information on cobra see: https://github.com/spf13/cobra
	rootCmd = &cobra.Command{
		Use:   "tempest-template",
		Short: "tempest-template is a template project for Go services at Shipt",
	}

	webserverCmd = &cobra.Command{
		Use:   "webserver",
		Short: "run the webserver",
		RunE: func(cmd *cobra.Command, args []string) error {
			var appCfg app.Config
			if err := cfg.Unmarshal(&appCfg); err != nil {
				return err
			}
			return app.
				New(appCfg).
				Run(context.Background(), runWebserver)
		},
	}

	workerCmd = &cobra.Command{
		Use:   "worker",
		Short: "run the worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			var appCfg app.Config
			if err := cfg.Unmarshal(&appCfg); err != nil {
				return err
			}
			return app.
				New(appCfg).
				Run(context.Background(), runWorker)
		},
	}
)

func runWebserver(ctx context.Context) error {
	var webserverCfg webserver.Config
	if err := cfg.Unmarshal(&webserverCfg); err != nil {
		return err
	}
	return webserver.Run(ctx, webserverCfg)
}

func runWorker(ctx context.Context) error {
	var workerCfg worker.Config
	if err := cfg.Unmarshal(&workerCfg); err != nil {
		return err
	}
	return worker.Run(ctx, workerCfg)
}

func main() {
	// setup commands ...
	rootCmd.AddCommand(webserverCmd)
	rootCmd.AddCommand(workerCmd)
	// to add a new command: rootCmd.AddCommand(myNewCmd)

	// rock 'n roll!
	if err := rootCmd.Execute(); err != nil {
		bubinga.Fatal(context.Background(), "crashing", err)
	}
}
