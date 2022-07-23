package cmd

import (
	"fmt"
	"os"

	"github.com/salfatigroup/nopeus/cli/util"
	"github.com/salfatigroup/nopeus/config"
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
    PersistentPostRun: func(cmd *cobra.Command, args []string) {
        cfg := config.GetNopeusConfig()

        // delete the tmp directory if KeepExecutionFiles is enabled
        if !cfg.Runtime.KeepExecutionFiles {
            if err := os.RemoveAll(cfg.Runtime.TmpFileLocation); err != nil {
                fmt.Printf("Error deleting tmp nopeus directory: %s\n", err)
            }
        }
    },
}

// add root command flags
func init() {
    cfg := config.GetNopeusConfig()

    // private command to keep the tmp directory (keep-execution-files)
    rootCmd.PersistentFlags().BoolVar(&cfg.Runtime.KeepExecutionFiles, "keep-execution-files", false, "Keep the execution files in the tmp directory")

    // mark the flag as hidden
    rootCmd.Flags().MarkHidden("keep-execution-files")
}


// Execute the CLI root command
func Execute() error {
    if err := rootCmd.Execute(); err != nil {
        return err
    }

    return nil
}

