package core

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/salfatigroup/nopeus/config"
)

// run and deploy terraform files per environment
func runTerraform(cfg *config.NopeusConfig) error {
    for _, env := range cfg.Runtime.Environments {
        workingTfDir := fmt.Sprintf("%s/%s/%s", cfg.Runtime.TmpFileLocation, cfg.CAL.CloudVendor, env)

        if err := runTerraformFile(cfg, workingTfDir); err != nil {
            return err
        }
    }

    return nil
}

// run and deploy terraform file
func runTerraformFile(cfg *config.NopeusConfig, workingTfDir string) error {
    tf, err := tfexec.NewTerraform(workingTfDir, cfg.Runtime.TerraformExecutablePath)
    if err != nil {
        return err
    }

    // initialize terraform
    fmt.Println("Initializing your cloud deployment...")
    if err := tf.Init(context.Background(), tfexec.Upgrade(true)); err != nil {
        return err
    }

    // plan the terraform file and output the plan file
    fmt.Println("Planning your cloud infrastructure...")
    newChanges, err := tf.Plan(context.Background())
    if err != nil {
        return err
    }

    // apply the plan in dry run mode file if new changes are found
    if newChanges {
        fmt.Println("Upading your cloud infrastructure... This can take a while, going to grab some coffee ☕️...")
        if !cfg.Runtime.DryRun {
            if err := tf.Apply(context.Background()); err != nil {
                return err
            }

            fmt.Println("Your cloud infrastructure has been updated.")
        }
    } else {
        fmt.Println("No new changes found in terraform plan 🤷‍♂️")
    }

    return nil
}
