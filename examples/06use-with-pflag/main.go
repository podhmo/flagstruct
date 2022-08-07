package main

import (
	"encoding/json"
	"os"

	"github.com/podhmo/flagstruct"
	"github.com/spf13/pflag"
)

type Config struct {
	Name string `flag:"name"`
}

var debug bool

func main() {
	fs := pflag.NewFlagSet("hello", pflag.ExitOnError)
	fs.BoolVarP(&debug, "debug", "", debug, "debug option")

	config := &Config{Name: "anonymous"}

	binder := &flagstruct.Binder{Config: flagstruct.DefaultConfig()}
	setenv := binder.Bind(fs, config)
	if err := setenv(fs); err != nil {
		panic(err)
	}

	fs.Parse(os.Args[1:])

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(map[string]interface{}{
		"debug":  debug,
		"config": config,
	})
}
