package docs

import (
	"embed"
)

const BasePath = "usage"

//go:embed usage/*.md
var DocFS embed.FS
