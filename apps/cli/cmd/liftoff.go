package cmd

import (
	"fmt"
	"os"

	"github.com/salfatigroup/nopeus/cli/util"
	"github.com/salfatigroup/nopeus/config"
	"github.com/salfatigroup/nopeus/core"
	"github.com/spf13/cobra"
)

// the config path as defined by the users flag
var configPath string
var dryRun bool

func init() {
    // init command after user argument is defined
    cobra.OnInitialize(initConfig)

    // define the liftoff flags
    liftoffCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to config file. Defaults to $( pwd )/nopeus.yaml")
    liftoffCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Dry run. Don't actually deploy to the cloud")

    // register new command
    rootCmd.AddCommand(liftoffCmd)
}

// define the command that build the configuration and
// deploy the application to the cloud
var liftoffCmd = &cobra.Command{
    Use:   "liftoff",
    Short: "Deploys your application layer to the cloud",
    Run: liftoff,
}

// This command parses the configuration file and
// deploys the application to the cloud
func liftoff(cmd *cobra.Command, args []string) {
    fmt.Println(
        "🔥",
        util.GradientText("[NOPEUS::STARTUP]", "#db2777", "#f9a8d4"),
        "- preparing your application for deployment to the cloud",
    )
    cfg := config.GetNopeusConfig()

    // deploy the application
    fmt.Println(
        "🚀",
        util.GradientText("[NOPEUS::LIFTOFF]", "#db2777", "#f9a8d4"),
        "- deploying your application to the cloud",
    )
    if err := core.Deploy(cfg); err != nil {
        fmt.Println(
            "💥",
            util.GradientText("[NOPEUS::TERMINATE]", "#db2777", "#f9a8d4"),
            " - failed to deploy your application to the cloud \n",
            err,
        )
        os.Exit(1)
    }

    fmt.Println(
        "🛰",
        util.GradientText("[NOPEUS::MECO]", "#db2777", "#f9a8d4"),
        "- your application is securely deployed to the cloud",
    )
}

// apply the provided user argument to the configs
func initConfig() {
    cfg := config.GetNopeusConfig()

    // apply the config path to the config
    if configPath != "" {
        cfg.SetConfigPath(configPath)
    }

    cfg.SetDryRun(dryRun)

    // initialize configs
    if err := cfg.Init(); err != nil {
        fmt.Println(
            "💥",
            util.GradientText("[NOPEUS::TERMINATE]", "#db2777", "#f9a8d4"),
            "- failed to initialize nopeus config \n",
            err,
        )
        os.Exit(1)
    }
}
