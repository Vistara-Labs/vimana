package scaffold

import (
	"embed"
	_ "embed"
)

//go:embed templates/*.*tmpl
var createTemplateFS embed.FS
