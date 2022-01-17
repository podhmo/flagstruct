# structflag

define flagset with struct and reflect

## features

- Builds pflag.FlagSet by struct definition
- Supports only a single use case (shorthand of flag package)
- (Default envvar support)

## install

(This package is currently under development (wip))

```console
$ go get -v github.com/podhmo/structflag
```

## how to use

```go
package main

import (
	"fmt"
	"os"

	"github.com/podhmo/structflag"
)

type Options struct {
	Name    string `flag:"name" help:"name of greeting"`
	Verbose bool   `flag:"verbose" short:"v"`
}

func main() {
	options := &Options{Name: "foo"} // default value

	b := structflag.NewBuilder()
	b.Name = "hello"
	b.EnvPrefix = "X_"

	fs := b.Build(options)
	fs.Parse(os.Args[1:])

	fmt.Printf("parsed: %#+v\n", options)

}
```

default help message.

```console
$ hello --help
Usage of hello:
      --name string   ENV: X_NAME	name of greeting (default "foo")
  -v, --verbose       ENV: X_VERBOSE	-
pflag: help requested
exit status 2
```

run it.

```console
$ hello --verbose
parsed: &main.Options{Name:"foo", Verbose:true}
$ hello -v --name bar
parsed: &main.Options{Name:"name", Verbose:true}

# envvar support
$ X_NAME=bar hello -v --name bar
parsed: &main.Options{Name:"bar", Verbose:true}
```

see also [./examples](./examples)