package cmd

import (
	"fmt"
	"os"

	"github.com/salfatigroup/gologsnag"
	"github.com/salfatigroup/nopeus/cli/util"
	"github.com/salfatigroup/nopeus/config"
	"github.com/salfatigroup/nopeus/core"
	"github.com/salfatigroup/nopeus/logger"
	"github.com/salfatigroup/nopeus/remote"
	"github.com/spf13/cobra"
)

// the config path as defined by the users flag
var configPath string

var (
	environmentsOverwrite []string
	overwriteVersion      string
)

func init() {
	// get global config
	cfg := config.GetNopeusConfig()

	// define the liftoff flags
	liftoffCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to config file. Defaults to $( pwd )/nopeus.yaml")
	liftoffCmd.Flags().BoolVar(&cfg.Runtime.DryRun, "dry-run", false, "Dry run. Don't actually deploy to the cloud")
	liftoffCmd.Flags().StringVarP(&cfg.Runtime.NopeusCloudToken, "token", "t", "", "Token to use for authentication")
	liftoffCmd.Flags().StringSliceVarP(&environmentsOverwrite, "env", "e", []string{}, "Deploy only specific environments out of the environments list in the nopeus.yaml configurations. Values passed to this flag must exists in the nopeus.yaml e.g., --env prod")
	liftoffCmd.Flags().StringVarP(&overwriteVersion, "version", "v", "", "Overwrite all the images version to deploy")

	// register new command
	rootCmd.AddCommand(liftoffCmd)
}

// define the command that build the configuration and
// deploy the application to the cloud
var liftoffCmd = &cobra.Command{
	Use:   "liftoff",
	Short: "Deploys your application layer to the cloud",
	Run:   liftoff,
}

// This command parses the configuration file and
// deploys the application to the cloud
func liftoff(cmd *cobra.Command, args []string) {
	fmt.Println(
		"🔥",
		util.GradientText("[NOPEUS::STARTUP]", "#db2777", "#f9a8d4"),
		"- preparing your application for deployment to the cloud",
	)

	// init configs
	initConfig()
	cfg := config.GetNopeusConfig()

	// deploy the application
	fmt.Println(
		"🚀",
		util.GradientText("[NOPEUS::LIFTOFF]", "#db2777", "#f9a8d4"),
		"- deploying your application to the cloud",
	)
	if err := core.Deploy(cfg); err != nil {
		logger.Publish(&gologsnag.PublishOptions{Event: "error", Description: err.Error(), Tags: &gologsnag.Tags{"func": "liftoff"}})
		logger.Errorf("Failed to deploy application: %+v", err)
		fmt.Println(
			"💥",
			util.GradientText("[NOPEUS::TERMINATE]", "#db2777", "#f9a8d4"),
			" - failed to deploy your application to the cloud \n",
			err,
		)
		os.Exit(1)
	}

	fmt.Println(
		"🛰 ",
		util.GradientText("[NOPEUS::MECO]", "#db2777", "#f9a8d4"),
		"- your application is securely deployed to the cloud",
	)
	logger.Debug("Liftoff command finished")
	logger.Publish(&gologsnag.PublishOptions{Event: "liftoff-finished", Icon: "🎉", Notify: true})
}

// apply the provided user argument to the configs
func initConfig() {
	logger.Publish(&gologsnag.PublishOptions{Event: "liftoff"})
	logger.Debug("Liftoff command called")
	cfg := config.GetNopeusConfig()

	// apply the config path to the config
	if configPath != "" {
		cfg.SetConfigPath(configPath)
	}

	// create remote session if token is provided
	if cfg.Runtime.NopeusCloudToken != "" {
		logger.Debug("Found nopeus token, creating remote session client")
		// verify the token now to prevent any errors later
		_, err := remote.NewRemoteSession(cfg.Runtime.NopeusCloudToken)
		if err != nil {
			logger.Publish(&gologsnag.PublishOptions{Event: "error", Description: err.Error(), Tags: &gologsnag.Tags{"func": "initConfig"}})
			logger.Errorf("Failed to create remote session: %+v", err)
			fmt.Println(
				"💥",
				util.GradientText("[NOPEUS::TERMINATE]", "#db2777", "#f9a8d4"),
				" - failed to create remote session \n",
				err,
			)
			os.Exit(1)
		}
	}

	// initialize configs
	if err := cfg.Init(); err != nil {
		logger.Publish(&gologsnag.PublishOptions{Event: "error", Description: err.Error(), Tags: &gologsnag.Tags{"func": "initConfig"}})
		logger.Errorf("Failed to initialize configs: %+v", err)
		fmt.Println(
			"💥",
			util.GradientText("[NOPEUS::TERMINATE]", "#db2777", "#f9a8d4"),
			"- failed to initialize nopeus config \n",
			err,
		)
		os.Exit(1)
	}

	logger.Publish(&gologsnag.PublishOptions{Event: "liftoff-config-initialized"})
	logger.Debug("Configs initialized")
}
