package configs

import (
	_ "embed"
)

var (
	//go:embed configs.yaml
	Config_yml []byte
)
