package templates

import "github.com/salfatigroup/nopeus/config"

// define the main infrastructure config interface
type EnvironmentConfigInterface interface {
    GetRendererValues() *TerraformRendererValues
}

// define the renderer values that will be used to render the terraform templates
type TerraformRendererValues struct {
    // the environment name
    Environment string
    // deployment name
    Name string
}

func getTFValues(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) *TerraformRendererValues {
    return &TerraformRendererValues{
        Environment: envName,
        Name: cfg.CAL.GetName(),
    }
}
