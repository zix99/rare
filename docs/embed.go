package docs

import (
	"embed"
)

//go:embed *.md
var DocFS embed.FS
