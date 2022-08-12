package config

import (
	"fmt"
	"os"
	"strings"

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

    // convert each environment variable from each service to
    // a matching value from environment variables when value
    // is in the following format ${ENV_VAR}
    if err := c.convertEnvVars(); err != nil {
        return err
    }

    return nil
}

func (c *NopeusConfig) convertEnvVars() error {
    // get the services
    services, err := c.CAL.GetServices()
    if err != nil {
        return err
    }

    // iterate through each service
    for _, service := range services {
        // iterate through each environment variable
        for key, value := range service.EnvironmentVariables {
            // check if the value is in the following format ${ENV_VAR}
            if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
                // get the environment variable name
                envVar := value[2 : len(value)-1]
                // get the environment variable value
                envValue := os.Getenv(envVar)
                // check if the environment variable value is empty
                if envValue == "" {
                    return fmt.Errorf("environment variable %s is not set", envVar)
                }
                // update the environment variable value
                service.EnvironmentVariables[key] = envValue
            }
        }
    }

    return nil
}
