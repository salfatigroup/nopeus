package config

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

// Define the config structure for nopeus
type NopeusConfig struct {
	// define the runtime config
	Runtime *RuntimeConfig

	// define the cloud application layer config
	CAL *CloudApplicationLayerConfig
}

// global singleton config
var ncfg *NopeusConfig

// create an instance of the config on startup
func init() {
	ncfg = NewNopeusConfig()
}

// create a new nopeus config instance
func NewNopeusConfig() *NopeusConfig {
	return &NopeusConfig{
		Runtime: NewRuntimeConfig(),
	}
}

// return the singleton config
func GetNopeusConfig() *NopeusConfig {
	return ncfg
}

// parse the nopeus config file on initialization
// uses the runtime config to parse the nopeus configs
func (c *NopeusConfig) parseConfig() error {
	// check if config path exists
	if _, err := os.Stat(c.Runtime.ConfigPath); os.IsNotExist(err) {
		return err
	}

	// read the config file
	file, err := os.Open(c.Runtime.ConfigPath)
	if err != nil {
		return err
	}

	// parse the nopeus yaml config and unmarshall it into the nopeus config
	if err := yaml.NewDecoder(file).Decode(&c.CAL); err != nil {
		return err
	}

	// init environment configs
	for envName, envData := range c.CAL.GetEnvironments() {
		if envData == nil {
			c.CAL.Environments[envName] = NewEnvironmentConfig()
		}
	}

	return nil
}
