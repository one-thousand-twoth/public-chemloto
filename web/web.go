package web

import "embed"

//go:embed pages
var IndexHTML embed.FS

//go:embed items static
var StaticFS embed.FS
