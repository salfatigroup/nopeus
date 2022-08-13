package cache

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/salfatigroup/nopeus/config"
)

// define the nopeus state object
type NopeusState struct {
    EnvironmentName string `json:"environment"`
    CloudVendor string `json:"cloud_vendor"`
    TerraformState string `json:"terraform_state"`
    DeployedServices []string `json:"deployed_services"`
}

// create a new nopeus state file
func NewNopeusState(envName string, envData *config.EnvironmentConfig, cfg *config.NopeusConfig) (*NopeusState, error) {
    // get the tfstate file location
    tfstateLocation := filepath.Join(
        cfg.Runtime.TmpFileLocation,
        cfg.CAL.CloudVendor,
        envName,
        "terraform.tfstate",
    )

    // read the tfstate file
    tfstate, err := readTfstate(tfstateLocation)
    if err != nil {
        return nil, err
    }

    // get services
    services, err := cfg.CAL.GetServices()
    if err != nil {
        return nil, err
    }

    // get all the keys from the deployed services
    deployedServices := getKeys(services)

    // create the nopeus state object
    nopeusState := &NopeusState{
        EnvironmentName: envName,
        CloudVendor: cfg.CAL.CloudVendor,
        TerraformState: tfstate,
        DeployedServices: deployedServices,
    }

    return nopeusState, nil
}

// write the list of nopeus states to the given location
func (s *NopeusState) WriteNopeusState(location string) error {
    // marshal the nopeus state object to json
    json, err := json.Marshal(s)
    if err != nil {
        return err
    }

    // ensure the directory exists
    dirName := filepath.Dir(location)
    // check if the directory exists
    if _, err := os.Stat(dirName); os.IsNotExist(err) {
        // create the directory
        if err := os.MkdirAll(dirName, 0755); err != nil {
            return err
        }
    }

    // write the json to the given location
    return ioutil.WriteFile(location, json, 0644)
}
