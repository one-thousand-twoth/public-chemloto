//go:build production
// +build production

package web

import "embed"

//go:embed dist
var DIST embed.FS

//go:embed polymers.json
var Polymers embed.FS
