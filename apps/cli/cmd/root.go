package cmd

import (
	"github.com/salfatigroup/nopeus/cli/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "nopeus",
    Short: util.ShortDescription(),
    Long: util.LongDescription(),
    Version: "0.0.1",
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            cmd.Help()
        }
    },
}

// Execute the CLI root command
func Execute() error {
    if err := rootCmd.Execute(); err != nil {
        return err
    }

    return nil
}

