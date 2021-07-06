// +build wasm,electron

package fs

import "github.com/gontikr99/chutzparse/pkg/nodejs"

var fs = nodejs.Require("fs")
var fsPromises = fs.Get("promises")
