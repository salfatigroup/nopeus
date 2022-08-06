package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// define the toke verification request
type TokenVerificationRequest struct {
    Token string `json:"token"`
}

// define the verification response
type TokenVerificationResponse struct {
    Status string `json:"status"`
}

// authenticate with the remote nopeus cloud server
func (session *RemoteSession) Authenticate() error {
    // create the request body
    reqBody, err := json.Marshal(TokenVerificationRequest{Token: session.token})
    if err != nil {
        return err
    }

    // create a new http client
    client := &http.Client{}

    // create a new request
    req, err := http.NewRequest(
        "POST",
        NOPEUSCLOUD_API_BASE_URL+"/licenses/v1/verify",
        bytes.NewBuffer(reqBody),
    )
    if err != nil {
        return err
    }

    // set the content type
    req.Header.Set("Content-Type", "application/json")

    // execute the request
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // decode the response body
    var resBody TokenVerificationResponse
    if err := json.NewDecoder(resp.Body).Decode(&resBody); err != nil {
        return err
    }

    // check the status code and body
    if resp.StatusCode != 200 || resBody.Status != "active" {
        return fmt.Errorf("failed to authenticate with the remote nopeus cloud server")
    }

    // set the token verified flag
    session.tokenVerified = true
    return nil
}

// check if the session is authenticated
func (session *RemoteSession) IsAuthenticated() bool {
    return session.tokenVerified
}
