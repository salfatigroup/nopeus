package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/salfatigroup/nopeus/cache"
	"github.com/salfatigroup/nopeus/config"
)

const (
    NOPEUS_CLOUD_ARTIFACTS_URI = "/api/artifacts/v1/state"
)

// upload each file in FILES_TO_CACHE from the tmp runtime
// directory to the nopeus cloud server.
func (s *RemoteSession) SetRemoteCache(cfg *config.NopeusConfig, newstate *cache.NopeusState) error {
    // check if token has been verified and authorized in the client
    // to reduce http requests
    if !s.tokenVerified {
        return fmt.Errorf("token not verified")
    }

    // upload the new state object
    if err := s.uploadFile(cfg, newstate); err != nil {
        return err
    }

    return nil
}

// get the remote cache from the nopeus cloud server
func (s *RemoteSession) GetRemoteCache(envName string) error {
    // check if token has been verified and authorized in the client
    // to reduce http requests
    if !s.tokenVerified {
        return fmt.Errorf("token not verified")
    }

    // get the config
    cfg := config.GetNopeusConfig()

    // get the remote state based on the stack name
    if err := s.downloadFile(cfg, envName); err != nil {
        return err
    }

    return nil
}

// mark the remote cache as in used to prevent
// terraform override between users
func (s *RemoteSession) LockRemoteState() error {
    return nil
}

// mark the remote cache as unused to allow
// nopeus operations
func (s *RemoteSession) UnlockRemoteState() error {
    return nil
}

// create a nopeus state object in salfati group cloud
func (s *RemoteSession) uploadFile(cfg *config.NopeusConfig, newstate *cache.NopeusState) error {
    endpoint := NOPEUSCLOUD_API_BASE_URL + NOPEUS_CLOUD_ARTIFACTS_URI

    // convert newstate to a json string
    body, err := json.Marshal(newstate)
    if err != nil {
        return err
    }

    // create the http request
    req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
    if err != nil {
        return err
    }

    // set the authorization header
    req.Header.Set("Authorization", "Token "+s.token)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")

    // send the request to the nopeus cloud server
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }

    // close the response body
    defer resp.Body.Close()

    // check the response status code
    if resp.StatusCode != 200 {
        return fmt.Errorf("failed to upload state to the remote nopeus cloud server")
    }

    return nil
}

// download a file from nopeus artifact storage
func (s *RemoteSession) downloadFile(cfg *config.NopeusConfig, envName string) error {
    // define the endpoint
    endpoint := NOPEUSCLOUD_API_BASE_URL + NOPEUS_CLOUD_ARTIFACTS_URI + "/" + cfg.CAL.GetName() + "-" + envName

    // create the http request
    req, err := http.NewRequest("GET", endpoint, nil)
    if err != nil {
        return err
    }

    // set the authorization header
    req.Header.Set("Authorization", "Token "+s.token)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")

    // send the request to the nopeus cloud server
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }

    // close the response body
    defer resp.Body.Close()

    // ignore 404s as state is not required to be found
    if resp.StatusCode == 404 {
        return nil
    }

    // check the response status code
    if resp.StatusCode != 200 {
        return fmt.Errorf("failed to download state from the remote nopeus cloud server")
    }

    // decode the response body
    var state cache.NopeusState
    if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
        return err
    }

    // write the state to the local file
    nopeusStateLocation := filepath.Join(cfg.Runtime.RootNopeusDir, "state", envName+".nopeus.state")
    if err := state.WriteNopeusState(nopeusStateLocation); err != nil {
        return err
    }

    return nil
}
