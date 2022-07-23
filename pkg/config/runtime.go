package config

import (
	"os"
	"path/filepath"
)

// define the nopeus runtime config
type RuntimeConfig struct {
    ConfigPath string
    HasBeenInitialized bool
}

// create a new instance of the runtime config with all the required default values
func NewRuntimeConfig() *RuntimeConfig {
    return &RuntimeConfig{
        ConfigPath: GetDefaultConfigPath(),
        HasBeenInitialized: false,
    }
}

// Once the HasBeenInitialized flag is set to true,
// it means the noepsu config has been loaded and is ready to be used
func (c *NopeusConfig) Init() error {
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
