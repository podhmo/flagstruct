package flagstruct_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/podhmo/flagstruct"
	"github.com/spf13/pflag"
)

func TestBuilder_Build(t *testing.T) {
	newBuilder := func() *flagstruct.Builder {
		b := flagstruct.NewBuilder()
		b.Name = "-"
		b.FlagnameTags = []string{"flag"}
		b.ShorthandTag = "short"
		b.EnvvarSupport = false
		b.HandlingMode = pflag.ContinueOnError
		return b
	}

	normalize := func(t *testing.T, ob interface{}) string {
		t.Helper()
		b, err := json.Marshal(ob)
		if err != nil {
			t.Fatalf("error %+v (encode)", err)
		}
		return string(b)
	}
	normalizeString := func(t *testing.T, s string) string {
		t.Helper()
		var ob interface{}
		err := json.Unmarshal([]byte(s), &ob)
		if err != nil {
			t.Fatalf("error %+v (encodeString unmarshal)", err)
		}
		b, err := json.Marshal(ob)
		if err != nil {
			t.Fatalf("error %+v (encodeString marshal)", err)
		}
		return string(b)
	}

	tests := []struct {
		name   string
		args   []string
		want   string
		create func() (*flagstruct.Builder, interface{})

		errorString string
	}{
		{
			name: "types--string",
			args: []string{"--name", "foo"},
			want: `{"Name":"foo"}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					Name string `flag:"name"`
				}
				return newBuilder(), &Options{}
			},
		},
		{
			name: "types--string-pointer",
			args: []string{"--name", "foo"},
			want: `{"Name":"foo"}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					Name *string `flag:"name"`
				}
				return newBuilder(), &Options{}
			},
		},
		{
			name: "types--string,default",
			args: []string{"--name2", "bar"},
			want: `{"Name":"foo", "Name2":"bar"}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					Name  string `flag:"name"`
					Name2 string `flag:"name2"`
				}
				return newBuilder(), &Options{Name: "foo", Name2: "foo"} // default value
			},
		},
		{
			name: "types--int",
			args: []string{"--age", "20"},
			want: `{"Age":20}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					Age int `flag:"age"`
				}
				return newBuilder(), &Options{}
			},
		},
		{
			name: "types--int-slice",
			args: []string{"-n", "20", "-n", "30"},
			want: `{"Nums": [20, 30]}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					Nums []int `flag:"nums" short:"n"`
				}
				return newBuilder(), &Options{}
			},
		},
		{
			name: "options--long",
			args: []string{"--verbose"},
			want: `{"Verbose":true}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					Verbose bool `flag:"verbose"`
				}
				return newBuilder(), &Options{}
			},
		},
		{
			name: "options--short",
			args: []string{"-v"},
			want: `{"Verbose":true}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					Verbose bool `flag:"verbose" short:"v"`
				}
				return newBuilder(), &Options{}
			},
		},
		{
			name: "options--short-only",
			args: []string{"-v"},
			want: `{"Verbose":true}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					Verbose bool `short:"v"`
				}
				return newBuilder(), &Options{}
			},
		},
		{
			name: "options--nothing",
			args: []string{"--Verbose"},
			want: `{"Verbose":true}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					Verbose bool
				}
				return newBuilder(), &Options{}
			},
		},
		{
			name: "skip--exported,tag",
			args: []string{"--name", "foo"},
			want: `{"name":"foo"}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					Name string `flag:"-"`
				}
				return newBuilder(), &Options{}
			},
			errorString: "unknown flag: --name",
		},
		{
			name: "skip--unexported",
			args: []string{"--name", "foo"},
			want: `{"name":"foo"}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					name string
				}
				return newBuilder(), &Options{}
			},
			errorString: "unknown flag: --name",
		},
		{
			name: "skip--unexported,tag,not-skipped",
			args: []string{"--name", "foo"},
			want: `{}`, // normalized by json, so unexported field is not included.
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					name string `flag:"name"`
				}
				return newBuilder(), &Options{}
			},
		},
		{
			name: "skip--pointer",
			args: []string{"--name", "foo"},
			want: `{"Name":"foo"}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					Name *string
				}
				return newBuilder(), &Options{}
			},
			errorString: "unknown flag: --name",
		},
		{
			name: "skip--pointer,tag,not-skipped",
			args: []string{"--name", "foo"},
			want: `{"Name":"foo"}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					Name *string `flag:"name"`
				}
				return newBuilder(), &Options{}
			},
		},
		{
			name: "lookup--tag--json",
			args: []string{"--verbose"},
			want: `{"verbose":true}`, // serialized by encoding/json
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					Verbose bool `json:"verbose"` // not flag
				}
				b := newBuilder()
				b.FlagnameTags = append(b.FlagnameTags, "json")
				return b, &Options{}
			},
		},
		{
			name: "lookup--tag--json,omitempty",
			args: []string{"--verbose"},
			want: `{"verbose":true}`, // serialized by encoding/json
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					Verbose bool `json:"verbose,omitempty"` // not flag
				}
				b := newBuilder()
				b.FlagnameTags = append(b.FlagnameTags, "json")
				return b, &Options{}
			},
		},
		{
			name: "customize--enum",
			args: []string{"--log-level", "info"},
			want: `{"LogLevel":"INFO", "LogLevelDefault": "WARN", "LogLevelPointer": "WARN"}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Options struct {
					LogLevel        LogLevel  `flag:"log-level"`
					LogLevelDefault LogLevel  `flag:"log-level-default"`
					LogLevelPointer *LogLevel `flag:"log-level-pointer"`
				}
				b := newBuilder()
				logDefault := LogLevelWarning
				return b, &Options{LogLevel: logDefault, LogLevelDefault: logDefault, LogLevelPointer: &logDefault}
			},
		},
		{
			name: "nested",
			args: []string{"--father.name", "foo"},
			want: `{"Father": {"Name": "foo"}, "Mother": {"Name": "moo"}, "Zero": {"Name": ""}}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Person struct {
					Name string `flag:"name"`
				}
				type Options struct {
					Father Person `flag:"father"`
					Mother Person `flag:"mother"`
					Zero   Person `flag:"zero"`
				}
				b := newBuilder()
				return b, &Options{Mother: Person{Name: "moo"}}
			},
		},
		{
			name: "nested,embedded",
			args: []string{"--name", "foo"},
			want: `{"Name": "foo"}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Person struct {
					Name string `flag:"name"`
				}
				type Options struct {
					Person
				}
				b := newBuilder()
				return b, &Options{}
			},
		},
		{
			name: "nested,pointer",
			args: []string{"--father.name", "foo"},
			want: `{"Father": {"Name": "foo"}, "Mother": {"Name": "moo"}, "Zero": {"Name": ""}}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Person struct {
					Name string `flag:"name"`
				}
				type Options struct {
					Father *Person `flag:"father"`
					Mother *Person `flag:"mother"`
					Zero   *Person `flag:"zero"`
				}
				b := newBuilder()
				return b, &Options{Mother: &Person{Name: "moo"}}
			},
		},
		{
			name: "nested,pointer,with-json",
			args: []string{"--zero.name", "foo"},
			want: `{"zero": {"name": "foo"}}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Person struct {
					Name string `json:"name"`
				}
				type Options struct {
					Zero *Person `json:"zero"`
				}
				b := newBuilder()
				b.FlagnameTags = append(b.FlagnameTags, "json")
				return b, &Options{}
			},
		},
		{
			name: "nested,pointer,with-json,override",
			args: []string{"--zero.name", "foo"},
			want: `{}`,
			create: func() (*flagstruct.Builder, interface{}) {
				type Person struct {
					Name string `json:"name"`
				}
				type Options struct {
					Zero *Person `json:"zero" flag:"-"`
				}
				b := newBuilder()
				b.FlagnameTags = append(b.FlagnameTags, "json")
				return b, &Options{}
			},
			errorString: "unknown flag: --zero.name",
		},

		// MEMO: []struct[T] is impossible. maybe.
		// {
		// 	name: "nested,slice",
		// 	args: []string{"--people.name", "foo"},
		// 	want: `{"People": [{"Name": "moo"}, {"Name": "foo"}]}`,
		// 	create: func() (*flagstruct.Builder, interface{}) {
		// 		type Person struct {
		// 			Name string `flag:"name"`
		// 		}
		// 		type Options struct {
		// 			People []Person `flag:"people"`
		// 		}
		// 		b := newBuilder()
		// 		return b, &Options{People: []Person{{Name: "moo"}}}
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, options := tt.create()
			fs := b.Build(options)

			err := fs.Parse(tt.args)
			if tt.errorString == "" {
				if err != nil {
					t.Fatalf("unexpected error: %+v with (%v)", err, tt.args) // TODO: help message
				}
			} else {
				if err == nil {
					t.Fatalf("must be error, but nil")
				}
				if tt.errorString != err.Error() {
					t.Fatalf("unexpected error: %+q\n\tbut expected message is %q", err.Error(), tt.errorString)
				}
				return
			}

			got := normalize(t, options)
			want := normalizeString(t, tt.want)
			t.Logf("parsed: %+#v", options)

			if got != want {
				t.Errorf("Builder.Build() = %v, want %v\nargs = %s", got, want, tt.args)
			}
		})
	}
}

// test for enum

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

// for flagstruct.HasHelpText
func (v LogLevel) HelpText() string {
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
