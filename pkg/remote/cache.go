package remote

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/salfatigroup/nopeus/config"
)

const (
    NOPEUS_CLOUD_ARTIFACTS_URI = "/artifacts/v1"
)

var (
    FILES_TO_CACHE = []string{
        "terraform.tfstate",
    }
)

// upload each file in FILES_TO_CACHE from the tmp runtime
// directory to the nopeus cloud server.
func (s *RemoteSession) SetRemoteCache(cfg *config.NopeusConfig) error {
    // check if token has been verified and authorized in the client
    // to reduce http requests
    if !s.tokenVerified {
        return fmt.Errorf("token not verified")
    }

    // for each file in FILES_TO_CACHE, upload it to the nopeus cloud server
    // at the NOPEUS_CLOUD_ARTIFACTS_URI
    for _, file := range FILES_TO_CACHE {
        s.uploadFile(cfg, file)
    }

    return nil
}

// get the remote cache from the nopeus cloud server
func (s *RemoteSession) GetRemoteCache() error {
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

// upload a file to nopeus artifact storage
// each request to nopeus cloud should have:
// 1. The token in the Authorization header
// 2. Upload the multipart form data with the file
// 3. Contain the `type` and `name` of the file in the form data
func (s *RemoteSession) uploadFile(cfg *config.NopeusConfig, file string) error {
    // get the file from the tmp directory
    filePath := cfg.Runtime.TmpFileLocation + "/" + file
    fileBytes, err := ioutil.ReadFile(filePath)
    if err != nil {
        return err
    }

    // create the multipart form data
    var buf bytes.Buffer
    w := multipart.NewWriter(&buf)
    fw, err := w.CreateFormFile("file", file)
    if err != nil {
        return err
    }

    // update the file name and file type in the form data
    if err := w.WriteField("type", "tfstate"); err != nil {
        return err
    }

    if err := w.WriteField("name", "terraform.tfstate"); err != nil {
        return err
    }

    // write the file to the multipart form data
    _, err = fw.Write(fileBytes)
    if err != nil {
        return err
    }

    // close the multipart form data
    w.Close()

    // create the http request
    req, err := http.NewRequest("POST", NOPEUSCLOUD_API_BASE_URL+NOPEUS_CLOUD_ARTIFACTS_URI, &buf)
    if err != nil {
        return err
    }

    // set the authorization header
    req.Header.Set("Authorization", "Token "+s.token)
    req.Header.Set("Content-Type", w.FormDataContentType())

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
        return fmt.Errorf("unable to upload file to nopeus cloud")
    }

    return nil
}
