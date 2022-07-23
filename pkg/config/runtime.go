package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// define the nopeus runtime config
type RuntimeConfig struct {
    ConfigPath string
    HasBeenInitialized bool
    Environments []string
    RootNopeusDir string
    TmpFileLocation string
    KeepExecutionFiles bool
    TerraformExecutablePath string
    DryRun bool
}

// create a new instance of the runtime config with all the required default values
func NewRuntimeConfig() *RuntimeConfig {
    // get the ~/.nopeus directory
    homeDir, _ := os.UserHomeDir()
    rootNopeusDir := filepath.Join(homeDir, ".nopeus")
    terraformPath, err := exec.LookPath("terraform")
    if err != nil {
        fmt.Println("terraform not found in PATH. Please install terraform to use nopeus.")
        os.Exit(1)
    }

    // return configs
    return &RuntimeConfig{
        // lookup the default config path at $( pwd )/nopeus.yaml
        ConfigPath: GetDefaultConfigPath(),

        // not initialized until the config is loaded
        HasBeenInitialized: false,

        // by default use only one environment - production
        Environments: []string{ "prod" },

        // by default use the root nopeus directory
        RootNopeusDir: rootNopeusDir,

        // point temp file location to tmp dir
        TmpFileLocation: filepath.Join(rootNopeusDir, "tmp"),

        // a hiddent command for debug purposes that stores
        // all the execution files in the tmp directory
        KeepExecutionFiles: false,

        // a path to terraform binary
        TerraformExecutablePath: terraformPath,

        // by default ignore the dry run mode
        DryRun: false,
    }
}

// Once the HasBeenInitialized flag is set to true,
// it means the noepsu config has been loaded and is ready to be used
func (c *NopeusConfig) Init() error {
    // initialize the root nopeus directory and temp directory
    // create the root nopeus directory if it doesn't exist
    if _, err := os.Stat(c.Runtime.RootNopeusDir); os.IsNotExist(err) {
        if err := os.MkdirAll(c.Runtime.RootNopeusDir, 0755); err != nil {
            return err
        }
    }

    // create the tmp directory if it doesn't exist
    if _, err := os.Stat(c.Runtime.TmpFileLocation); os.IsNotExist(err) {
        if err := os.MkdirAll(c.Runtime.TmpFileLocation, 0755); err != nil {
            return err
        }
    }

    // find and parse user nopeus config file
    if err := c.parseConfig(); err != nil {
        return err
    }

    // mark the config as initialized
    c.Runtime.HasBeenInitialized = true
    return nil
}

// define the nopeus config
func (c *NopeusConfig) SetConfigPath(path string) {
    c.Runtime.ConfigPath = path
}

// Return the default nopeus config path
func GetDefaultConfigPath() string {
    pwd, _ := os.Getwd()
    return filepath.Join(pwd, "nopeus.yaml")
}

// define the dry run mode in the global runtime config
func (c *NopeusConfig) SetDryRun(dryRun bool) {
    c.Runtime.DryRun = dryRun
}
