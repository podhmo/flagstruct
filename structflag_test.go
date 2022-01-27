package structflag_test

import (
	"encoding/json"
	"testing"

	"github.com/podhmo/structflag"
	"github.com/spf13/pflag"
)

func TestBuilder_Build(t *testing.T) {
	newBuilder := func() *structflag.Builder {
		b := structflag.NewBuilder()
		b.Name = "-"
		b.FlagnameTag = "flag"
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
		create func() (*structflag.Builder, interface{})
	}{
		{
			name: "string",
			args: []string{"--name", "foo"},
			want: `{"Name":"foo"}`,
			create: func() (*structflag.Builder, interface{}) {
				type Options struct {
					Name string `flag:"name"`
				}
				return newBuilder(), &Options{}
			},
		},
		{
			name: "options",
			args: []string{"--long", "foo", "-s", "bar"},
			want: `{"Long":"foo", "Short": "bar"}`,
			create: func() (*structflag.Builder, interface{}) {
				type Options struct {
					Long  string `flag:"long"`
					Short string `flag:"short" short:"s"`
				}
				return newBuilder(), &Options{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, options := tt.create()
			fs := b.Build(options)

			if err := fs.Parse(tt.args); err != nil {
				t.Fatalf("parse error: %+v with (%v)", err, tt.args) // TODO: help message
			}

			got := normalize(t, options)
			want := normalizeString(t, tt.want)

			if got != want {
				t.Errorf("Builder.Build() = %v, want %v\nargs = %s", got, want, tt.args)
			}
		})
	}
}
