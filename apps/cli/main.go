package main

import (
	"fmt"

	"github.com/salfatigroup/nopeus/cli/cmd"
)

// Nopeus adds an application layer to the cloud.
// Simply define your applications and let nopeus do the rest.
func main() {
    if err := cmd.Execute(); err != nil {
        fmt.Errorf("%v", err)
    }
}
