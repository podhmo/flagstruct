package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/podhmo/flagstruct"
)

type Options struct {
	LogLevel  slog.Level  `flag:"log-level" help:"log level  {DEBUG, INFO, WARN, ERROR}"`
	LogLevel2 *slog.Level `flag:"log-level2" help:"log level  {DEBUG, INFO, WARN, ERROR}"`
}

func main() {
	options := &Options{LogLevel: slog.LevelInfo} // default value

	flagstruct.Parse(options, func(b *flagstruct.Builder) {
		b.Name = "textvar"
		b.EnvPrefix = "X_"
	})

	fmt.Printf("parsed: %#+v\n", options)
	fmt.Println("----------------------------------------")
	fmt.Println("json:")
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(options)
}
