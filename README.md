# flagstruct

A Go library that automatically builds pflag.FlagSet from struct definitions using reflection. Define your CLI options as a struct, and flagstruct handles all the flag registration, environment variable binding, and parsing for you.

## features

- Automatic flag generation - Builds pflag.FlagSet by struct definition with minimal boilerplate
- Type-safe - Full type safety with Go structs
- Nested structure support - Organize flags hierarchically ([example](./examples/03nested/), [shared common option](./examples/08shared-common-option/))
- Environment variable support - Automatic binding with customizable prefix
- Rich type support - All basic types, slices, time.Duration, custom types via flag.Value and encoding.TextUnmarshaler
- Short flags - Define shorthand flags via struct tags
- Required flags - Mark flags as required with validation
- Cobra integration - Works seamlessly with Cobra commands ([example](./examples/07use-with-cobra/))
- Custom help text - Customize help messages via struct tags or interfaces

## install

```console
$ go get -v github.com/podhmo/flagstruct
```

## how to use

### Basic usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/podhmo/flagstruct"
)

type Options struct {
	Name    string `flag:"name" help:"name of greeting"`
	Verbose bool   `flag:"verbose" short:"v"`
}

func main() {
	options := &Options{Name: "foo"} // default value

	flagstruct.Parse(options, func(b *flagstruct.Builder) {
		b.Name = "hello"
		b.EnvPrefix = "X_"
	})

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
parsed: &main.Options{Name:"bar", Verbose:true}

# envvar support
$ X_NAME=bar hello -v
parsed: &main.Options{Name:"bar", Verbose:true}
```

## Supported types

flagstruct supports the following types out of the box:

Basic types:
- `string`
- `bool`
- `int`, `int64`
- `uint`, `uint64`
- `float64`
- `time.Duration`

Slice types:
- `[]string`
- `[]bool`
- `[]int`, `[]int64`
- `[]uint`
- `[]float64`
- `[]time.Duration`

Custom types:
- Types implementing `flag.Value` interface (Set, String, Type methods)
- Types implementing `encoding.TextUnmarshaler` and `encoding.TextMarshaler` interfaces

## Struct tags

### `flag` tag
Defines the flag name. If omitted, unexported fields are ignored.

```go
type Options struct {
    Name string `flag:"name"`        // --name
    Age  int    `flag:"age"`         // --age
    Skip string `flag:"-"`           // ignored
    Port int                         // ignored (no flag tag, unexported not supported)
}
```

### `short` tag
Defines a short flag alias. Only works for top-level flags (not nested).

```go
type Options struct {
    Verbose bool `flag:"verbose" short:"v"`  // -v, --verbose
    Debug   bool `flag:"debug" short:"d"`    // -d, --debug
}
```

### `help` tag
Provides help text for the flag.

```go
type Options struct {
    Name string `flag:"name" help:"Your name"`
    Port int    `flag:"port" help:"Server port number"`
}
```

### `required` tag
Marks a flag as required. Validation fails if not provided.

```go
type Options struct {
    APIKey string `flag:"api-key" required:"true" help:"API key for authentication"`
}
```

## Advanced usage

### Nested structures

Organize related flags into nested structures. Flag names are prefixed with the struct field name.

```go
type DBConfig struct {
    URI   string `flag:"uri"`
    Debug bool   `flag:"debug"`
}

type Options struct {
    DB        DBConfig `flag:"db"`         // --db.uri, --db.debug
    AnotherDB DBConfig `flag:"another-db"` // --another-db.uri, --another-db.debug
}

func main() {
    options := &Options{}
    flagstruct.Parse(options)
}
```

### Shared common options

Share common configuration across multiple nested structs using embedded pointers.

```go
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
    BaseConfig  // top-level embedded
    A AConfig `json:"a"`
    B BConfig `json:"b"`
}

func main() {
    config := &Config{}
    flagstruct.Parse(config, flagstruct.WithMoreFlagnameTags("json"))
    // --debug flag is shared between BaseConfig, AConfig.BaseConfig, and BConfig.BaseConfig
}
```

### Custom types with flag.Value

Implement the `flag.Value` interface for custom types:

```go
type LogLevel string

const (
    LogLevelDebug LogLevel = "DEBUG"
    LogLevelInfo  LogLevel = "INFO"
    LogLevelWarn  LogLevel = "WARN"
    LogLevelError LogLevel = "ERROR"
)

func (v *LogLevel) String() string {
    return string(*v)
}

func (v *LogLevel) Set(value string) error {
    *v = LogLevel(strings.ToUpper(value))
    switch *v {
    case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
        return nil
    default:
        return fmt.Errorf("invalid log level: %s", value)
    }
}

func (v *LogLevel) Type() string {
    return "LogLevel"
}

// Optional: implement HasHelpText for custom help message
func (v LogLevel) HelpText() string {
    return "log level {DEBUG, INFO, WARN, ERROR}"
}

type Options struct {
    LogLevel LogLevel `flag:"log-level"`
}
```

### Using with Cobra

flagstruct works seamlessly with Cobra commands:

```go
import (
    "github.com/podhmo/flagstruct"
    "github.com/spf13/cobra"
)

type Options struct {
    Name    string `flag:"name" help:"name of greeting"`
    Verbose bool   `flag:"verbose" short:"v"`
}

func main() {
    options := &Options{Name: "default"}
    
    cmd := &cobra.Command{
        Use:   "hello",
        Short: "A brief description",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("Hello %s\n", options.Name)
        },
    }
    
    fs := flagstruct.Build(options)
    cmd.Flags().AddFlagSet(fs.FlagSet)
    
    cmd.Execute()
}
```

### Error handling

Use `Build` and `ParseArgs` for custom error handling:

```go
options := &Options{}
fs := flagstruct.Build(options, func(b *flagstruct.Builder) {
    b.EnvPrefix = "MYAPP_"
})

if err := fs.Parse(os.Args[1:]); err != nil {
    // Handle error
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}
```

### Environment variables

Environment variables are automatically supported with customizable prefix:

```go
flagstruct.Parse(options, func(b *flagstruct.Builder) {
    b.EnvPrefix = "MYAPP_"  // flags will check MYAPP_<FLAGNAME> env vars
})
```

Flag names are converted to environment variable names by:
1. Converting to uppercase
2. Replacing `-` and `.` with `_`
3. Adding the prefix

Example: `--db.uri` becomes `MYAPP_DB_URI`

### Additional flag name tags

Use multiple tag names (e.g., combine with `json` tags):

```go
type Options struct {
    Name string `json:"name" help:"Your name"`
}

flagstruct.Parse(options, flagstruct.WithMoreFlagnameTags("json"))
// This creates --name flag from the json tag
```

## API reference

### `Parse[T any](o *T, options ...func(*Builder))`
Parse command line arguments with automatic error handling and exit on error.

### `ParseArgs[T any](o *T, args []string, options ...func(*Builder))`
Parse specific arguments with automatic error handling.

### `Build[T any](o *T, options ...func(*Builder)) *FlagSet`
Build a FlagSet without parsing. Useful for Cobra integration or custom error handling.

### Builder options

```go
func(b *flagstruct.Builder) {
    b.Name = "myapp"              // Program name for help
    b.EnvPrefix = "MYAPP_"        // Environment variable prefix
    b.HandlingMode = flag.ExitOnError  // Error handling mode
}
```

### `WithMoreFlagnameTags(tags ...string) func(b *Builder)`
Add additional struct tags to use for flag names (e.g., "json", "yaml").

## Examples

See [./examples](./examples) for more complete examples:

- [00hello](./examples/00hello/) - Basic usage
- [01json-tag](./examples/01json-tag/) - Using json tags
- [02enum](./examples/02enum/) - Custom enum types
- [03nested](./examples/03nested/) - Nested structures
- [04textvar](./examples/04textvar/) - TextUnmarshaler support
- [05return-error](./examples/05return-error/) - Error handling
- [06use-with-pflag](./examples/06use-with-pflag/) - Direct pflag usage
- [07use-with-cobra](./examples/07use-with-cobra/) - Cobra integration
- [08shared-common-option](./examples/08shared-common-option/) - Shared options
- [10cobra-sub-command](./examples/10cobra-sub-command/) - Cobra subcommands

## inspired by

- https://github.com/heetch/confita
- https://github.com/AdamSLevy/flagbind
