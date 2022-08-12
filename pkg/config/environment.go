package config

import "github.com/hashicorp/terraform-exec/tfexec"

// define the environment configs
type EnvironmentConfig struct {
    outputs map[string]tfexec.OutputMeta
}

func NewEnvironmentConfig() *EnvironmentConfig {
    return &EnvironmentConfig{}
}

// set the terraform outputs to the infrastructure config
func (i *EnvironmentConfig) SetOutputs(outputs map[string]tfexec.OutputMeta) {
    i.outputs = outputs
}

// returns the terraform outputs
func (i *EnvironmentConfig) GetOutputs() map[string]tfexec.OutputMeta {
    return i.outputs
}
