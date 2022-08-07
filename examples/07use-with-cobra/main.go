package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/podhmo/flagstruct"
	"github.com/spf13/cobra"
)

type Config struct {
	Name string `json:"name" flag:"name" required:"true" help:"name for the user"`
}

var debug bool

func main() {
	if err := Execute(); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "hello",
	Short: "hello world",
	Long: `hello world
	hello world`,

	Run: func(cmd *cobra.Command, args []string) {
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(map[string]interface{}{
			"debug":  debug,
			"config": config,
		})
	},
}

var config = &Config{}

func init() {
	rootCmd.Flags().BoolVarP(&debug, "debug", "", false, "debug option")

	binder := &flagstruct.Binder{Config: flagstruct.DefaultConfig()}
	binder.EnvPrefix = "X_"
	fs := rootCmd.Flags()

	setenv := binder.Bind(fs, config)
	if err := setenv(fs); err != nil {
		panic(err)
	}

	for _, requiredName := range binder.AllRequiredFlagNames() {
		cobra.MarkFlagRequired(fs, requiredName) // e.g. required flag(s) "name" not set
	}
}

func Execute() error {
	return rootCmd.Execute()
}
