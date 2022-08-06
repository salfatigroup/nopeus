package remote

import "os"

var (
    NOPEUSCLOUD_API_BASE_URL = "https://api.nopeus.salfati.group"
)

func init() {
    goenv := os.Getenv("GO_ENV")
    if goenv == "development" {
        NOPEUSCLOUD_API_BASE_URL = "http://localhost:8000"
    }
}
