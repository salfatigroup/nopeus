package templates

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/salfatigroup/nopeus/config"
)

// Generate a terraform environment file including all the relevant modules
// to a given destination path
func GenerateTerraformEnvironment(cfg *config.NopeusConfig, env string) error {
    // define the destination path of the modules
    destLocation, err := getTfDestLocation(cfg.Runtime.TmpFileLocation, cfg.CAL.CloudVendor, env)
    if err != nil {
        return err
    }

    // for each file in the embedded templates directory (StaticTerraformTemplates) recursivly
    // generate a terraform file in the destination directory
    // after rendering the template
    if err := renderTerraformTemplates(cfg, destLocation, env, "terraform"); err != nil {
        return err
    }

    return nil
}

// return the destination directory to copy to
func getTfDestLocation(tmpfileLocation string, cloudVendor string, env string) (string, error) {
    location := filepath.Join(tmpfileLocation, cloudVendor, env)

    // if the location does not exist, create it
    if _, err := os.Stat(location); os.IsNotExist(err) {
        if err := os.MkdirAll(location, 0755); err != nil {
            return "", err
        }
    }

    return location, nil
}

// generate the terraform files recursively
// from the embedded StaticTerraformTemplates
func renderTerraformTemplates(cfg *config.NopeusConfig, destLocation string, env string, dir string) error {
    dirs, err := StaticTerraformTemplates.ReadDir(dir)
    if err != nil {
        return err
    }

    for _, d := range dirs {
        if d.IsDir() {
            if err := renderTerraformTemplates(cfg, destLocation, env, filepath.Join(dir, d.Name())); err != nil {
                return err
            }
        } else {
            // render the template
            if err := renderTerraformFile(cfg, destLocation, env, filepath.Join(dir, d.Name())); err != nil {
                return err
            }
        }
    }

    return nil
}

// render a single terraform file template
func renderTerraformFile(cfg *config.NopeusConfig, destLocation string, env string, file string) error {
    // render the template
    rendered, err := renderTemplate(cfg, file, env)
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

// render a specific template
func renderTemplate(cfg *config.NopeusConfig, file string, env string) (string, error) {
    // read the template file
    templateContent, err := StaticTerraformTemplates.ReadFile(file)
    if err != nil {
        return "", err
    }

    // render the template
    tmpl, err := template.New(file).Parse(string(templateContent))
    if err != nil {
        return "", err
    }

    // create the buffer to write the rendered template to
    var renderedBuffer bytes.Buffer

    // render the template
    if err := tmpl.Execute(&renderedBuffer, cfg); err != nil {
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
