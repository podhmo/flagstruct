package flagstruct

import (
	"encoding"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	flag "github.com/spf13/pflag"
)

type HasHelpText interface {
	HelpText() string
}

// TODO: map

type Config struct {
	HandlingMode flag.ErrorHandling

	EnvvarSupport bool
	EnvPrefix     string
	EnvNameFunc   func(string) string

	FlagnameTags []string
	FlagNameFunc func(string) string

	ShorthandTag string
	HelpTextTag  string
	RequiredTag  string
}

func DefaultConfig() *Config {
	c := &Config{
		FlagnameTags:  []string{"flag"},
		ShorthandTag:  "short",
		HelpTextTag:   "help",
		RequiredTag:   "required",
		EnvvarSupport: true,
		HandlingMode:  flag.ExitOnError,
	}
	if v := os.Getenv("ENV_PREFIX"); v != "" {
		c.EnvPrefix = v
	}
	c.EnvNameFunc = func(name string) string {
		return c.EnvPrefix + strings.ReplaceAll(strings.ReplaceAll(strings.ToUpper(name), "-", "_"), ".", "_")
	}
	c.FlagNameFunc = func(v string) string {
		if strings.Contains(v, ",") {
			return strings.TrimSpace(strings.SplitN(v, ",", 2)[0]) // e.g. json's omitempty
		}
		return v
	}
	return c
}

var (
	rTimeDurationType    reflect.Type
	rFlagValueType       reflect.Type
	rTextUnmarshalerType reflect.Type
)

func init() {
	rTimeDurationType = reflect.TypeOf(time.Second)
	rFlagValueType = reflect.TypeOf(func() flag.Value { return nil }).Out(0)
	rTextUnmarshalerType = reflect.TypeOf(func() encoding.TextUnmarshaler { return nil }).Out(0)
}

type Builder struct {
	Name string
	*Config
}

func NewBuilder() *Builder {
	name := os.Args[0]
	b := &Builder{Name: name, Config: DefaultConfig()}
	return b
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

	binder := &Binder{Config: b.Config}
	binder.State.toplevelStructMap = map[reflect.Type]reflect.Value{}
	binder.State.embeddedStructPointerMap = map[reflect.Type][]reflect.Value{}

	binder.walk(fs, rt, rv, "")

	return &FlagSet{FlagSet: fs, Binder: binder}
}

type Binder struct {
	*Config

	State struct {
		visitedFields []fieldcontext

		toplevelStructMap        map[reflect.Type]reflect.Value
		embeddedStructPointerMap map[reflect.Type][]reflect.Value
	}
}

func (b *Binder) Bind(fs *flag.FlagSet, o interface{}) func(*flag.FlagSet) error {
	rt := reflect.TypeOf(o)
	rv := reflect.ValueOf(o)

	if rt.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("%v is not pointer of struct", rt)) // for canAddr
	}
	rt = rt.Elem()
	rv = rv.Elem()

	b.State.toplevelStructMap = map[reflect.Type]reflect.Value{}
	b.State.embeddedStructPointerMap = map[reflect.Type][]reflect.Value{}

	b.walk(fs, rt, rv, "")

	// for shared common option
	if len(b.State.embeddedStructPointerMap) > 0 {
		if err := b.setSharedCommonEmbeddedStruct(); err != nil {
			return func(*flag.FlagSet) error { return err }
		}
	}
	return b.setByEnvvars
}

func (b *Binder) setByEnvvars(fs *flag.FlagSet) (retErr error) {
	fs.VisitAll(func(f *flag.Flag) {
		envname := b.EnvNameFunc(f.Name)
		if v, ok := os.LookupEnv(envname); ok {
			if err := fs.Set(f.Name, v); err != nil {
				retErr = fmt.Errorf("on envvar %s=%v, %+v", envname, v, err)
			}
		}
	})
	return retErr
}

func (b *Binder) setSharedCommonEmbeddedStruct() error {
	for ft, fvs := range b.State.embeddedStructPointerMap {
		base, ok := b.State.toplevelStructMap[ft.Elem()]
		if !ok {
			return fmt.Errorf("shared common option is not found (%v)", ft.Elem())
		}
		for _, fv := range fvs {
			fv.Set(base.Addr())
		}
	}
	return nil
}

func (b *Binder) AllRequiredFlagNames() []string {
	var required []string
	for _, fc := range b.State.visitedFields {
		if fc.required {
			required = append(required, fc.fieldname)
		}
	}
	return required
}

func (b *Binder) ValidateRequiredFlags(fs *flag.FlagSet) error {
	for _, requiredName := range b.AllRequiredFlagNames() {
		if !fs.Lookup(requiredName).Changed {
			return fmt.Errorf("required flag(s) %q not set", requiredName)
		}
	}
	return nil
}

func (b *Binder) walk(fs *flag.FlagSet, rt reflect.Type, rv reflect.Value, prefix string) {
	for i := 0; i < rt.NumField(); i++ {
		rf := rt.Field(i)
		fv := rv.Field(i)

		fieldname := rf.Name
		hasFlagname := false

		{
			for j := len(b.FlagnameTags) - 1; j >= 0; j-- {
				if v, ok := rf.Tag.Lookup(b.FlagnameTags[j]); ok {
					fieldname = v
					hasFlagname = true
				}
			}
			if fieldname == "-" {
				continue
			}
			if !hasFlagname && !rf.IsExported() {
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

		required := false
		if ok, _ := strconv.ParseBool(rf.Tag.Get(b.RequiredTag)); ok {
			required = true
		}
		if required {
			helpText = helpText + " [required]"
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

		fc := fieldcontext{
			fieldname: fieldname,
			helpText:  helpText,
			required:  required,
			shorthand: shorthand,

			prefix:      prefix,
			hasFlagname: hasFlagname,
			field:       rf,
		}

		b.State.visitedFields = append(b.State.visitedFields, fc)
		b.walkField(fs, rf.Type, fv, fc)
	}
}

type fieldcontext struct {
	fieldname string
	helpText  string
	shorthand string
	required  bool

	prefix      string
	hasFlagname bool
	field       reflect.StructField
}

func (b *Binder) walkField(fs *flag.FlagSet, rt reflect.Type, fv reflect.Value, c fieldcontext) {
	// for enum (TODO: skip check with cache)
	{
		fv := fv
		ft := fv.Type()
		isPtr := ft.Kind() == reflect.Ptr

		// Set() is pointer receiver only
		if !isPtr {
			fv = fv.Addr()
			ft = reflect.PtrTo(ft)
		}

		if ft.Implements(rFlagValueType) {
			fs.VarP(fv.Interface().(flag.Value), c.fieldname, c.shorthand, c.helpText)
			return
		}

		if ft.Implements(rTextUnmarshalerType) {
			if fv.IsNil() && fv.CanAddr() {
				// flagname is not found, will be skipped (even if the field is a pointer, with field tag, it will be treated as a flag forcely).
				if !c.hasFlagname {
					return
				}
				fv.Set(reflect.New(rt.Elem()))
			}

			ref := fv.Interface().(interface {
				encoding.TextMarshaler
				encoding.TextUnmarshaler
			})
			fs.VarP(newTextValue(ref, ref, isPtr), c.fieldname, c.shorthand, c.helpText)
			return
		}
	}

	switch rt.Kind() {
	case reflect.Ptr:
		if fv.IsNil() && fv.CanAddr() {
			// flagname is not found, will be skipped (even if the field is a pointer, with field tag, it will be treated as a flag forcely).
			if !c.hasFlagname {
				if c.field.Anonymous {
					// for shared common option (child)
					typ := c.field.Type
					b.State.embeddedStructPointerMap[typ] = append(b.State.embeddedStructPointerMap[typ], fv)
				}
				return
			}
			fv.Set(reflect.New(rt.Elem()))
		}
		b.walkField(fs, rt.Elem(), fv.Elem(), c)
	case reflect.Struct:
		if c.prefix == "" {
			// for shared common option (parent)
			b.State.toplevelStructMap[fv.Type()] = fv
		}

		if c.field.Anonymous {
			b.walk(fs, rt, fv, c.prefix)
			return
		}
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

type FlagSet struct {
	*flag.FlagSet
	Binder *Binder
}

func (fs *FlagSet) Parse(args []string) error {
	if err := fs.FlagSet.Parse(args); err != nil {
		return err
	}

	// for envar
	if fs.Binder.EnvvarSupport {
		if err := fs.Binder.setByEnvvars(fs.FlagSet); err != nil {
			return err
		}
	}

	// for shared common option
	if len(fs.Binder.State.embeddedStructPointerMap) > 0 {
		if err := fs.Binder.setSharedCommonEmbeddedStruct(); err != nil {
			return err
		}
	}

	return fs.Binder.ValidateRequiredFlags(fs.FlagSet)
}

func Build[T any](o *T, options ...func(*Builder)) *FlagSet {
	b := NewBuilder()
	b.HandlingMode = flag.ContinueOnError
	for _, opt := range options {
		opt(b)
	}
	return b.Build(o)
}

func ParseArgs[T any](o *T, args []string, options ...func(*Builder)) {
	b := NewBuilder()
	b.HandlingMode = flag.ExitOnError
	for _, opt := range options {
		opt(b)
	}
	fs := b.Build(o)
	if err := fs.Parse(args); err != nil { // never error in fs.FlagSet.Parse(), because default handling mode is not ContinueOnError
		PrintHelpAndExitIfError(fs.FlagSet, err, 2)
	}
}

func Parse[T any](o *T, options ...func(*Builder)) {
	args := os.Args[1:]
	ParseArgs(o, args, options...)
}

func WithMoreFlagnameTags(tags ...string) func(b *Builder) {
	return func(b *Builder) {
		b.FlagnameTags = append(b.FlagnameTags, tags...)
	}
}

func PrintHelpAndExitIfError(fs *flag.FlagSet, err error, code int) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n%+v\n", err)
		os.Exit(code)
	}
}
