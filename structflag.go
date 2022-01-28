package structflag

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
	"unsafe"

	flag "github.com/spf13/pflag"
)

type HasHelpText interface {
	HelpText() string
}

// TODO: nested
// TODO: map
// TODO: embed

type Builder struct {
	Name         string
	HandlingMode flag.ErrorHandling

	EnvvarSupport bool
	EnvPrefix     string
	EnvNameFunc   func(string) string

	FlagnameTag  string
	FlagNameFunc func(string) string

	ShorthandTag string
	HelpTextTag  string
}

func NewBuilder() *Builder {
	name := os.Args[0]
	b := &Builder{
		Name:          name,
		FlagnameTag:   "flag",
		ShorthandTag:  "short",
		HelpTextTag:   "help",
		EnvvarSupport: true,
		HandlingMode:  flag.ExitOnError,
	}
	if v := os.Getenv("ENV_PREFIX"); v != "" {
		b.EnvPrefix = v
	}
	b.EnvNameFunc = func(name string) string {
		return b.EnvPrefix + strings.ReplaceAll(strings.ToUpper(name), "-", "_")
	}
	b.FlagNameFunc = func(v string) string {
		if strings.Contains(v, ",") {
			return strings.TrimSpace(strings.SplitN(v, ",", 2)[0]) // e.g. json's omitempty
		}
		return v
	}
	return b
}

type FlagSet struct {
	*flag.FlagSet
	builder *Builder
}

var (
	rTimeDurationType reflect.Type
	rFlagValueType    reflect.Type
)

func init() {
	rTimeDurationType = reflect.TypeOf(time.Second)
	rFlagValueType = reflect.TypeOf(func() flag.Value { return nil }).Out(0)
}

func (b *Builder) Build(o interface{}) *FlagSet {
	rt := reflect.TypeOf(o)
	rv := reflect.ValueOf(o)

	if rt.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("%v is not pointer of struct", rt)) // for canAddr
	}
	rt = rt.Elem()
	rv = rv.Elem()

	name := b.Name
	if name == "" {
		name = rt.Name()
	}
	fs := flag.NewFlagSet(name, b.HandlingMode)
	b.walk(fs, rt, rv, "")
	return &FlagSet{FlagSet: fs, builder: b}
}

func (b *Builder) walk(fs *flag.FlagSet, rt reflect.Type, rv reflect.Value, prefix string) {
	for i := 0; i < rt.NumField(); i++ {
		rf := rt.Field(i)
		fv := rv.Field(i)

		fieldname := rf.Name
		if v, ok := rf.Tag.Lookup(b.FlagnameTag); ok {
			if v == "-" {
				continue
			}
			fieldname = b.FlagNameFunc(prefix + v)
		} else {
			if !rf.IsExported() {
				continue
			}
			fieldname = b.FlagNameFunc(prefix + fieldname)
		}

		helpText := "-"
		if v, ok := rf.Tag.Lookup(b.HelpTextTag); ok {
			helpText = v
		} else {
			// for enum, for custom help message
			if fv.CanInterface() {
				impl, ok := fv.Interface().(HasHelpText)
				if ok {
					helpText = impl.HelpText()
				}
			}
		}
		if b.EnvvarSupport {
			helpText = fmt.Sprintf("ENV: %s\t", b.EnvNameFunc(fieldname)) + helpText
		}

		shorthand := ""
		if v, ok := rf.Tag.Lookup(b.ShorthandTag); ok {
			if prefix == "" {
				shorthand = v
			}
		}

		b.walkField(fs, rf.Type, fv, fieldcontext{
			fieldname: fieldname,
			helpText:  helpText,
			shorthand: shorthand,
			prefix:    prefix,
		})
	}
}

type fieldcontext struct {
	fieldname string
	helpText  string
	shorthand string

	prefix string
}

func (b *Builder) walkField(fs *flag.FlagSet, rt reflect.Type, fv reflect.Value, c fieldcontext) {
	// for enum (TODO: skip check with cache)
	{
		fv := fv
		ft := fv.Type()

		// Set() is pointer receiver only
		if ft.Kind() != reflect.Ptr {
			fv = fv.Addr()
			ft = reflect.PtrTo(ft)
		}

		if ft.Implements(rFlagValueType) {
			rfn := reflect.ValueOf(fs.VarP)
			rfn.Call([]reflect.Value{
				fv,
				reflect.ValueOf(c.fieldname),
				reflect.ValueOf(c.shorthand),
				reflect.ValueOf(c.helpText),
			})
			return
		}
	}

	switch rt.Kind() {
	case reflect.Struct:
		b.walk(fs, rt, fv, c.prefix+c.fieldname+".")
	case reflect.Bool:
		ref := (*bool)(unsafe.Pointer(fv.UnsafeAddr()))
		fs.BoolVarP(ref, c.fieldname, c.shorthand, fv.Bool(), c.helpText)
	case reflect.Float64:
		ref := (*float64)(unsafe.Pointer(fv.UnsafeAddr()))
		fs.Float64VarP(ref, c.fieldname, c.shorthand, fv.Float(), c.helpText)
	case reflect.Int64:
		switch rt {
		case rTimeDurationType:
			ref := (*time.Duration)(unsafe.Pointer(fv.UnsafeAddr()))
			fs.DurationVarP(ref, c.fieldname, c.shorthand, time.Duration(fv.Int()), c.helpText)
		default:
			ref := (*int64)(unsafe.Pointer(fv.UnsafeAddr()))
			fs.Int64VarP(ref, c.fieldname, c.shorthand, fv.Int(), c.helpText)
		}
	case reflect.Int:
		ref := (*int)(unsafe.Pointer(fv.UnsafeAddr()))
		fs.IntVarP(ref, c.fieldname, c.shorthand, int(fv.Int()), c.helpText)
	case reflect.String:
		ref := (*string)(unsafe.Pointer(fv.UnsafeAddr()))
		fs.StringVarP(ref, c.fieldname, c.shorthand, fv.String(), c.helpText)
	case reflect.Uint64:
		ref := (*uint64)(unsafe.Pointer(fv.UnsafeAddr()))
		fs.Uint64VarP(ref, c.fieldname, c.shorthand, fv.Uint(), c.helpText)
	case reflect.Uint:
		ref := (*uint)(unsafe.Pointer(fv.UnsafeAddr()))
		fs.UintVarP(ref, c.fieldname, c.shorthand, uint(fv.Uint()), c.helpText)
	case reflect.Slice:
		switch rt.Elem().Kind() {
		case reflect.Bool:
			var defaultValue []bool
			for i := 0; i < fv.Len(); i++ {
				defaultValue = append(defaultValue, fv.Index(i).Bool())
			}
			ref := (*[]bool)(unsafe.Pointer(fv.UnsafeAddr()))
			fs.BoolSliceVarP(ref, c.fieldname, c.shorthand, defaultValue, c.helpText)
		case reflect.Float64:
			var defaultValue []float64
			for i := 0; i < fv.Len(); i++ {
				defaultValue = append(defaultValue, fv.Index(i).Float())
			}
			ref := (*[]float64)(unsafe.Pointer(fv.UnsafeAddr()))
			fs.Float64SliceVarP(ref, c.fieldname, c.shorthand, defaultValue, c.helpText)
		case reflect.Int64:
			switch rt.Elem() {
			case rTimeDurationType:
				ref := (*[]time.Duration)(unsafe.Pointer(fv.UnsafeAddr()))
				var defaultValue []time.Duration
				for i := 0; i < fv.Len(); i++ {
					defaultValue = append(defaultValue, time.Duration(fv.Index(i).Int()))
				}
				fs.DurationSliceVarP(ref, c.fieldname, c.shorthand, defaultValue, c.helpText)
			default:
				var defaultValue []int64
				for i := 0; i < fv.Len(); i++ {
					defaultValue = append(defaultValue, fv.Index(i).Int())
				}
				ref := (*[]int64)(unsafe.Pointer(fv.UnsafeAddr()))
				fs.Int64SliceVarP(ref, c.fieldname, c.shorthand, defaultValue, c.helpText)
			}
		case reflect.Int:
			var defaultValue []int
			for i := 0; i < fv.Len(); i++ {
				defaultValue = append(defaultValue, int(fv.Index(i).Int()))
			}
			ref := (*[]int)(unsafe.Pointer(fv.UnsafeAddr()))
			fs.IntSliceVarP(ref, c.fieldname, c.shorthand, defaultValue, c.helpText)
		case reflect.String:
			var defaultValue []string
			for i := 0; i < fv.Len(); i++ {
				defaultValue = append(defaultValue, fv.Index(i).String())
			}
			ref := (*[]string)(unsafe.Pointer(fv.UnsafeAddr()))
			fs.StringSliceVarP(ref, c.fieldname, c.shorthand, defaultValue, c.helpText)
		case reflect.Uint:
			var defaultValue []uint
			for i := 0; i < fv.Len(); i++ {
				defaultValue = append(defaultValue, uint(fv.Index(i).Uint()))
			}
			ref := (*[]uint)(unsafe.Pointer(fv.UnsafeAddr()))
			fs.UintSliceVarP(ref, c.fieldname, c.shorthand, defaultValue, c.helpText)
		// case reflect.Uint64:
		default:
			panic(fmt.Sprintf("unsupported slice type %v", rt))
		}
	default:
		// TODO: map
		panic(fmt.Sprintf("unsupported type %v", rt))
	}
}

func (fs *FlagSet) Parse(args []string) error {
	err := fs.FlagSet.Parse(args)
	if err != nil {
		return err
	}
	if fs.builder.EnvvarSupport {
		fs.FlagSet.VisitAll(func(f *flag.Flag) {
			envname := fs.builder.EnvNameFunc(f.Name)
			if v := os.Getenv(envname); v != "" {
				if err := fs.Set(f.Name, v); err != nil {
					panic(fmt.Sprintf("on envvar %s=%v, %+v", envname, v, err))
				}
			}
		})
	}
	return nil
}
