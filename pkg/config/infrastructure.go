package config

import "github.com/hashicorp/terraform-exec/tfexec"

// define the main infrastructure config interface
type InfrastructureConfigInterface interface {
    GetRendererValues() *TerraformRendererValues
}

// define the renderer values that will be used to render the terraform templates
type TerraformRendererValues struct {
    // the environment name
    Environment string
}

// define the infrastucture config
type InfrastructureConfig struct {
    environment string
    outputs map[string]tfexec.OutputMeta
}

// set the terraform outputs to the infrastructure config
func (i *InfrastructureConfig) SetOutputs(outputs map[string]tfexec.OutputMeta) {
    i.outputs = outputs
}

// returns the terraform outputs
func (i *InfrastructureConfig) GetOutputs() map[string]tfexec.OutputMeta {
    return i.outputs
}

// return the renderer values
func (i *InfrastructureConfig) GetRendererValues() *TerraformRendererValues {
    return &TerraformRendererValues{
        Environment: i.environment,
    }
}
