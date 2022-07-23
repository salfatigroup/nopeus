package templates

import "embed"

//go:embed terraform
var StaticTerraformTemplates embed.FS
