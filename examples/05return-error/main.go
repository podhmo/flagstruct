package main

import (
	"fmt"
	"log"
	"os"

	"github.com/podhmo/flagstruct"
	"github.com/spf13/pflag"
)

type Options struct {
	Value   int  `flag:"value" required:"true"`
	Verbose bool `flag:"verbose" short:"v"`
}

func main() {
	options := &Options{Value: 10} // default value
	fs := flagstruct.Build(options, func(b *flagstruct.Builder) {
		b.Name = "hello"

		if b.HandlingMode != pflag.ContinueOnError {
			panic("unexpected handling mode")
		}
	})

	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatalf("hmm: %+v", err)
	}

	fmt.Printf("parsed: %#+v\n", options)
}
