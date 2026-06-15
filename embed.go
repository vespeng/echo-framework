package vortego

import "embed"

//go:embed internal/config/config.default.yaml
var ConfigFile embed.FS
