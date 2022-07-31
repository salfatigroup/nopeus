package core

import (
	"github.com/salfatigroup/nopeus/config"
	"github.com/salfatigroup/nopeus/templates"
)

// generateTerraformFiles generates the terraform files
// based on the provided configurations
func generateTerraformFiles(cfg *config.NopeusConfig) error {
    // iterate over the infrastructure configs map[string]InfrastructureConfig
    for env, infra := range cfg.Runtime.Infrastructure {
        if err := templates.GenerateTerraformEnvironment(cfg, env, infra); err != nil {
            return err
        }
    }

    return nil
}
