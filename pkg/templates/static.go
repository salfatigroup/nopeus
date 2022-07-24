package templates

import "embed"

//go:embed terraform
var StaticTerraformTemplates embed.FS

//go:embed helm
var StaticHelmTemplates embed.FS
