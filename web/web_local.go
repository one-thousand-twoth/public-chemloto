//go:build !production
// +build !production

package web

import (
	"embed"
	"io/fs"
	"os"
)

type localFS struct {
	root string
}

func (l localFS) Open(name string) (fs.File, error) {
	return os.Open(l.root + "/" + name)
}

var DIST fs.FS = localFS{root: "./dist"}

//go:embed polymers.json
var Polymers embed.FS
