// +build native

package main

import (
	"encoding/json"
	"github.com/gontikr99/chutzparse/internal"
	"os"
)

func main() {
	json.NewEncoder(os.Stdout).Encode(internal.PackageJson)
}
