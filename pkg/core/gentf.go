package core

import (
	"github.com/salfatigroup/nopeus/config"
	"github.com/salfatigroup/nopeus/templates"
)

// generateTerraformFiles generates the terraform files
// based on the provided configurations
func generateTerraformFiles(cfg *config.NopeusConfig) error {
    for _, env := range cfg.Runtime.Environments {
        if err := templates.GenerateTerraformEnvironment(cfg, env); err != nil {
            return err
        }
    }

    return nil
}
