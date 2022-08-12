package core

import (
	"github.com/salfatigroup/nopeus/config"
	"github.com/salfatigroup/nopeus/templates"
)

// generateTerraformFiles generates the terraform files
// based on the provided configurations
func generateTerraformFiles(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) error {
    // iterate over the infrastructure configs map[string]InfrastructureConfig
    if err := templates.GenerateTerraformEnvironment(cfg, envName, envData); err != nil {
        return err
    }

    return nil
}
