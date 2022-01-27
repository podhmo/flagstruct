package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/podhmo/structflag"
)

type LogLevel string

const (
	LogLevelDebug   LogLevel = "DEBUG"
	LogLevelInfo    LogLevel = "INFO"
	LogLevelWarning LogLevel = "WARN"
	LogLevelError   LogLevel = "ERROR"
)

func (v LogLevel) Validate() error {
	switch v {
	case "DEBUG", "INFO", "WARN", "ERROR":
		return nil
	default:
		return fmt.Errorf("%v is an invalid value for %v", v, reflect.TypeOf(v))
	}
}

// for structflag.HasHelpText
func (v *LogLevel) HelpText() string {
	return "log level {DEBUG, INFO, WARN, ERROR}"
}

// for pflag.Value
func (v *LogLevel) String() string {
	if v == nil {
		return "<nil>"
	}
	return string(*v)
}

// for pflag.Value
func (v *LogLevel) Set(value string) error {
	if v == nil {
		return fmt.Errorf("nil is invalid for %v", reflect.TypeOf(v))
	}
	*v = LogLevel(strings.ToUpper(value))
	return v.Validate()
}

// for pflag.Value
func (v *LogLevel) Type() string {
	return "LogLevel"
}

type Options struct {
	Name      string    `flag:"name" help:"name of greeting"`
	LogLevel  LogLevel  `flag:"log-level"`
	LogLevel2 *LogLevel `flag:"log-level2"`
}

func main() {
	options := &Options{Name: "foo", LogLevel: LogLevelInfo} // default value

	b := structflag.NewBuilder()
	b.Name = "hello"
	b.EnvPrefix = "X_"

	fs := b.Build(options)
	fs.Parse(os.Args[1:])

	fmt.Printf("parsed: %#+v\n", options)

}
