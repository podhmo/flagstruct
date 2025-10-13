# flagstruct

A Go library that automatically builds pflag.FlagSet from struct definitions using reflection. Define your CLI options as a struct, and flagstruct handles all the flag registration, environment variable binding, and parsing for you.

## features

- Automatic flag generation - Builds pflag.FlagSet by struct definition with minimal boilerplate
- Nested structure support - Organize flags hierarchically ([example](./examples/03nested/), [shared common option](./examples/08shared-common-option/))
- Environment variable support - Automatic binding with customizable prefix
- Short flags - Define shorthand flags via struct tags
- Required flags - Mark flags as required with validation
- Custom help text - Customize help messages via struct tags or interfaces
- Interspersed arguments - Supports interspersed flags and non-flag arguments by default

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

Map types:
- `map[string]string`
- `map[string]bool`
- `map[string]int`, `map[string]int64`
- `map[string]uint`, `map[string]uint64`
- `map[string]float64`
- `map[string]time.Duration`

Map values are passed as a single comma-separated string of `key=value` pairs. For example: `--metadata project=alpha,owner=john.doe`

Custom types:
- Types implementing `flag.Value` interface (Set, String, Type methods)
- Types implementing `encoding.TextUnmarshaler` and `encoding.TextMarshaler` interfaces

In addition to individual custom types, slices and maps of types that implement `encoding.TextUnmarshaler` are also supported. For example, you can use `[]net.IP` or `map[string]net.IP` directly in your configuration struct.

```go
type Config struct {
    IPAddresses []net.IP `flag:"ips"`
    IPMapping   map[string]net.IP `flag:"ip-map"`
}
```

You can then provide these flags on the command line:

```sh
$ go run main.go --ips "192.168.1.1,10.0.0.1" --ip-map "private=192.168.1.1,public=8.8.8.8"
```

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

### `envname` tag
Specifies a custom environment variable name for a flag.

```go
type Options struct {
    // This flag will be bound to the `MY_APP_PORT` environment variable
    Port int `flag:"port" envname:"MY_APP_PORT" help:"Server port number"`
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

### Embedded fields

Embedded structs allow you to compose configuration from multiple sources. Fields from embedded structs are promoted to the parent level.

```go
type DatabaseConfig struct {
    Host string `flag:"host"`
    Port int    `flag:"port"`
}

type Options struct {
    DatabaseConfig  // embedded struct
}

func main() {
    options := &Options{}
    flagstruct.Parse(options)
    // flags: --host, --port (fields are promoted to parent level)
}
```

The `--host` and `--port` flags from the embedded `DatabaseConfig` struct are available directly at the top level, not as `--databaseconfig.host` or `--databaseconfig.port`.

#### `inline` option in `flag` tag

Alternatively, you can achieve the same result using the `inline` option without embedding:

```go
type DatabaseConfig struct {
    Host string `flag:"host"`
    Port int    `flag:"port"`
}

type Options struct {
    DB DatabaseConfig `flag:",inline"`  // flags: --host, --port (not --db.host, --db.port)
}
```

Both approaches flatten nested struct fields to the parent level without adding a prefix. Use embedded structs when you want field promotion in Go code, or use `inline` when you want to keep the nested structure in code but flatten flags.

### Shared common options

Share common configuration across multiple nested structs using embedded pointers. This pattern is useful when you want a single flag to control the same setting across multiple components.

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
    // When you set --debug=true, all three embedded BaseConfig instances point to the same value
    fmt.Printf("Debug: %v, A.Debug: %v, B.Debug: %v\n", 
        config.Debug, config.A.Debug, config.B.Debug)
    // Output: Debug: true, A.Debug: true, B.Debug: true
}
```

### Interspersed Arguments

By default, `flagstruct` allows flags and non-flag arguments to be interspersed. This means you can mix flags and other arguments on the command line.

For example, with the following struct:

```go
type Options struct {
    Name    string `flag:"name"`
    Verbose bool   `flag:"verbose" short:"v"`
}
```

You can run your program with arguments like this:

```console
$ my-app --name foo arg1 --verbose arg2
```

`flagstruct` will correctly parse `--name` as "foo" and `--verbose` as `true`, while the remaining non-flag arguments (`arg1`, `arg2`) can be retrieved from the `*pflag.FlagSet`.

```go
options := &Options{}
fs := flagstruct.Build(options)
fs.Parse(os.Args[1:])

fmt.Printf("options: %#v\n", options)
fmt.Printf("args: %v\n", fs.Args())
```

Output:
```
options: &main.Options{Name:"foo", Verbose:true}
args: [arg1 arg2]
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

Environment variables are automatically supported with a customizable prefix. Command-line arguments take precedence over environment variables.

You can customize the environment variable name for a specific field using the `envname` tag.

```go
flagstruct.Parse(options, func(b *flagstruct.Builder) {
    b.EnvPrefix = "MYAPP_"  // flags will check MYAPP_<FLAGNAME> env vars
})
```

If `envname` tag is not provided, flag names are converted to environment variable names by:
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
