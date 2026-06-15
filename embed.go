package vortego

import "embed"

//go:embed .env
var EnvFile embed.FS

//go:embed internal/config/config.default.yaml
var ConfigFile embed.FS
