package internal

import (
	"io"
	"os"

	"github.com/fatih/color"
)

type OutputConfig struct {
	NoColor    bool `flag:"no-color"`
	WithStderr bool `flag:"with-stderr"`
}

func (c *OutputConfig) New(fg color.Attribute) *Output {
	w := os.Stdout
	Debug("output with-stderr is %v", c.WithStderr)
	if c.WithStderr {
		w = os.Stderr
	}
	Debug("output no-color is %v", c.NoColor)
	return &Output{
		w:       w,
		NoColor: c.NoColor,
		fg:      fg,
	}
}

type Output struct {
	w       io.Writer
	fg      color.Attribute
	NoColor bool
}

func (o *Output) Echo(msg string) {
	c := color.New(color.Bold, o.fg)
	if o.NoColor {
		c.DisableColor()
	}
	c.Fprintln(o.w, msg)
}
