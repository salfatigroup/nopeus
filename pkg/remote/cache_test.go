package remote

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/salfatigroup/nopeus/config"
)

func TestMain(m *testing.M) {
    setup()
    code := m.Run()
    // shutdown()
    os.Exit(code)
}

func setup() {
    // load .env file
    err := godotenv.Load("../../.env")
    if err != nil {
        panic(err)
    }
}


// TestRemoteCache tests the remote cache functionality
// by uploading a file to the nopeus cloud server
func TestSetRemoteCache(t *testing.T) {
    // get the configs and token required for the test
    cfg := config.GetNopeusConfig()
    token := os.Getenv("NOPEUS_TOKEN")

    // create a new remote session
    session, err := NewRemoteSession(token)
    if err != nil {
        t.Errorf("error creating remote session: %s", err)
        return
    }

    // call SetRemoteCache to upload a file to the nopeus cloud server
    err = session.SetRemoteCache(cfg)
    if err != nil {
        t.Errorf("error setting remote cache: %s", err)
        return
    }
}

