package main

import (
	"encoding/json"
	"os"

	"github.com/podhmo/flagstruct"
)

type BaseConfig struct {
	Debug bool `json:"debug"`
}

type AConfig struct {
	*BaseConfig
	Name string `json:"name"`
}
type BConfig struct {
	*BaseConfig
	Verbose bool `json:"verbose"`
}
type Config struct {
	BaseConfig

	A AConfig `json:"a"`
	B BConfig `json:"b"`
}

func main() {
	config := &Config{}
	flagstruct.Parse(config, func(b *flagstruct.Builder) {
		b.FlagnameTags = append(b.FlagnameTags, "json")
	})

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(config)
}
