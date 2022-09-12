package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/salfatigroup/nopeus/cli/util"
	"github.com/salfatigroup/nopeus/helm"
	helmrepo "helm.sh/helm/v3/pkg/repo"
)

// store that helm runtime data
type HelmRuntime struct {
	// the service template data that will be used to render the helm charts
	ServiceTemplateData []ServiceTemplateData
}

// define the nopeus runtime config
type RuntimeConfig struct {
	// the nopeus config location
	ConfigPath string

	// has the config been initialized yet or not
	HasBeenInitialized bool

	// the root nopeus directory
	RootNopeusDir string

	// the location of the tmp directory
	TmpFileLocation string

	// a hidden command for debug purposes that stores
	// all the execution files in the tmp directory
	RemoveKeepExecutionFiles bool

	// a path to terraform binary
	TerraformExecutablePath string

	// dry run mode - will not apply any changes to the cloud
	DryRun bool

	// the runtime services data that will be used to render
	// the final helm charts
	HelmRuntime *HelmRuntime

	// default helm repos to load on init
	HelmRepos []*helmrepo.Entry

	// set the default namespace for the main deployments
	DefaultNamespace string

	// nopeus cloud token
	NopeusCloudToken string
}

// create a new instance of the runtime config with all the required default values
func NewRuntimeConfig() *RuntimeConfig {
	// get the ~/.nopeus directory
	currentDir, _ := os.Getwd()
	rootNopeusDir := filepath.Join(currentDir, ".nopeus")
	terraformPath, err := exec.LookPath("terraform")
	if err != nil {
		fmt.Println(
			"💥 ",
			util.GradientText("[NOPEUS::TERMINATE]", "#db2777", "#f9a8d4"),
			" - terraform not found in PATH",
		)
		os.Exit(1)
	}

	// return configs
	runtime := &RuntimeConfig{
		// lookup the default config path at $( pwd )/nopeus.yaml
		ConfigPath: GetDefaultConfigPath(),

		// not initialized until the config is loaded
		HasBeenInitialized: false,

		// by default use the root nopeus directory
		RootNopeusDir: rootNopeusDir,

		// point temp file location to tmp dir
		TmpFileLocation: filepath.Join(rootNopeusDir, "session"),

		// a hiddent command for debug purposes that stores
		// all the execution files in the tmp directory
		RemoveKeepExecutionFiles: false,

		// a path to terraform binary
		TerraformExecutablePath: terraformPath,

		// by default ignore the dry run mode
		DryRun: false,

		// the helm runtime data
		HelmRuntime: &HelmRuntime{},

		// default helm repos to load on init
		HelmRepos: []*helmrepo.Entry{
			{
				Name: "salfatigroup",
				URL:  "https://charts.salfati.group",
			},
			{
				Name: "kong",
				URL:  "https://charts.konghq.com",
			},
			{
				Name: "bitnami",
				URL:  "https://charts.bitnami.com/bitnami",
			},
		},

		// default namespace for the main deployments
		DefaultNamespace: "nopeus-app",

		// nopeus cloud token
		NopeusCloudToken: "",
	}

	return runtime
}

// Once the HasBeenInitialized flag is set to true,
// it means the noepsu config has been loaded and is ready to be used
func (c *NopeusConfig) Init() error {
	// find and parse user nopeus config file
	if err := c.parseConfig(); err != nil {
		return err
	}

	// initialize the root nopeus directory and temp directory
	// create the root nopeus directory if it doesn't exist
	if _, err := os.Stat(c.Runtime.RootNopeusDir); os.IsNotExist(err) {
		if err := os.MkdirAll(c.Runtime.RootNopeusDir, 0o755); err != nil {
			return err
		}
	}

	// create the tmp directory if it doesn't exist
	if _, err := os.Stat(c.Runtime.TmpFileLocation); os.IsNotExist(err) {
		if err := os.MkdirAll(c.Runtime.TmpFileLocation, 0o755); err != nil {
			return err
		}
	}

	// load helm repos
	for _, repo := range c.Runtime.HelmRepos {
		if err := helm.AddChartRepo(*repo); err != nil {
			return err
		}
	}

	// mark the config as initialized
	c.Runtime.HasBeenInitialized = true
	return nil
}

// define the nopeus config
func (c *NopeusConfig) SetConfigPath(path string) {
	c.Runtime.ConfigPath = path
	// side effect - update the root nopeus directory
	rootNopeus := filepath.Join(filepath.Dir(path), ".nopeus")
	c.Runtime.RootNopeusDir = rootNopeus
	c.Runtime.TmpFileLocation = filepath.Join(rootNopeus, "session")
}

// Return the default nopeus config path
func GetDefaultConfigPath() string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, "nopeus.yaml")
}

// define the dry run mode in the global runtime config
func (c *NopeusConfig) SetDryRun(dryRun bool) {
	c.Runtime.DryRun = dryRun
}
