package configs

import (
	"embed"
	"io/fs"
)

//go:embed *.env
var embeddedEnv embed.FS

var EnvFS fs.ReadFileFS = embeddedEnv
