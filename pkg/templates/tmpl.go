package templates

import (
	"bytes"
	"os"
	"path/filepath"
	tmpl "text/template"

	"github.com/salfatigroup/nopeus/config"
)

// Generate a terraform environment file including all the relevant modules
// to a given destination path
func GenerateTerraformEnvironment(cfg *config.NopeusConfig, envName string, envData *config.EnvironmentConfig) error {
	// get the cloud vendor
	cloudVendor, err := cfg.CAL.GetCloudVendor()
	if err != nil {
		return err
	}

	// define the destination path of the modules
	destLocation, err := getTfDestLocation(cfg.Runtime.TmpFileLocation, cloudVendor, envName)
	if err != nil {
		return err
	}

	// for each file in the embedded templates directory (StaticTerraformTemplates) recursivly
	// generate a terraform file in the destination directory
	// after rendering the template
	if err := renderTerraformTemplates(cfg, destLocation, envName, envData, "terraform"); err != nil {
		return err
	}

	return nil
}

// return the destination directory to copy to
func getTfDestLocation(tmpfileLocation string, cloudVendor string, env string) (string, error) {
	location := filepath.Join(tmpfileLocation, cloudVendor, env)

	// if the location does not exist, create it
	if _, err := os.Stat(location); os.IsNotExist(err) {
		if err := os.MkdirAll(location, 0o755); err != nil {
			return "", err
		}
	}

	return location, nil
}

// generate the terraform files recursively
// from the embedded StaticTerraformTemplates
func renderTerraformTemplates(cfg *config.NopeusConfig, destLocation string, envName string, envData *config.EnvironmentConfig, dir string) error {
	dirs, err := StaticTerraformTemplates.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, d := range dirs {
		if d.IsDir() {
			if err := renderTerraformTemplates(cfg, destLocation, envName, envData, filepath.Join(dir, d.Name())); err != nil {
				return err
			}
		} else {
			// render the template
			if err := renderTerraformFile(cfg, destLocation, envName, envData, filepath.Join(dir, d.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}

// render a single terraform file template
func renderTerraformFile(cfg *config.NopeusConfig, destLocation string, envName string, envData *config.EnvironmentConfig, file string) error {
	// render the template
	rendered, err := renderTfTemplate(cfg, file, envName, envData)
	if err != nil {
		return err
	}

	// write the rendered template to the destination location
	filename := filepath.Base(file)
	destFile := filepath.Join(destLocation, filename)
	if err := writeFile(destFile, rendered); err != nil {
		return err
	}

	return nil
}

// render a helm values chart
func RenderHelmTemplateFile(runtimeServices config.ServiceTemplateData) (err error) {
	if runtimeServices.GetHelmValuesFile() != "" {
		templateContent := make([]byte, 0)
		valuesTemplate := runtimeServices.GetHelmValuesTemplate()

		if valuesTemplate != "" {
			// get the template file
			templateContent, err = StaticHelmTemplates.ReadFile(filepath.Join("helm", valuesTemplate))
			if err != nil {
				return err
			}
		}

		// render template
		tmpl, err := tmpl.New(valuesTemplate).
			Funcs(GetTempalteFuncs()).
			Parse(string(templateContent))
		if err != nil {
			return err
		}

		// create the buffer to write the rendered template to
		var renderedBuffer bytes.Buffer

		// render the template
		if err := tmpl.Execute(&renderedBuffer, runtimeServices.GetHelmValues()); err != nil {
			return err
		}

		return writeFile(runtimeServices.GetHelmValuesFile(), renderedBuffer.String())
	}

	return nil
}

// render a specific template
func renderTfTemplate(cfg *config.NopeusConfig, file string, envName string, envData *config.EnvironmentConfig) (string, error) {
	// read the template file
	templateContent, err := StaticTerraformTemplates.ReadFile(file)
	if err != nil {
		return "", err
	}

	// render the template
	tmpl, err := tmpl.New(file).
		Funcs(GetTempalteFuncs()).
		Parse(string(templateContent))
	if err != nil {
		return "", err
	}

	// create the buffer to write the rendered template to
	var renderedBuffer bytes.Buffer

	// render the template
	if err := tmpl.Execute(&renderedBuffer, getTFValues(envName, envData, cfg)); err != nil {
		return "", err
	}

	return renderedBuffer.String(), nil
}

// write a file to the given location
func writeFile(file string, content string) error {
	// create the file
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	// write the content to the file
	if _, err := f.WriteString(content); err != nil {
		return err
	}

	// close the file
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
