package structflag

import (
	"fmt"
	"os"
	"reflect"
	"time"
	"unsafe"

	flag "github.com/spf13/pflag"
)

// TODO: nested
// TODO: map
// TODO: embed

type Builder struct {
	Name         string
	HandlingMode flag.ErrorHandling

	FlagnameTag  string
	ShorthandTag string
	HelpTextTag  string
}

func NewBuilder() *Builder {
	name := os.Args[0]
	b := &Builder{
		Name:         name,
		FlagnameTag:  "json",
		ShorthandTag: "short",
		HelpTextTag:  "help",
		HandlingMode: flag.ExitOnError,
	}
	return b
}

type FlagSet struct {
	*flag.FlagSet
	builder *Builder
}

var (
	rTimeDuration = reflect.TypeOf(time.Second)
)

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

	for i := 0; i < rt.NumField(); i++ {
		rf := rt.Field(i)
		if !rf.IsExported() {
			continue
		}

		fieldname := rf.Name
		if v, ok := rf.Tag.Lookup(b.FlagnameTag); ok {
			fieldname = v
		}
		helpText := "-"
		if v, ok := rf.Tag.Lookup(b.HelpTextTag); ok {
			helpText = v
		}
		shorthand := ""
		if v, ok := rf.Tag.Lookup(b.ShorthandTag); ok {
			shorthand = v
		}

		fv := rv.Field(i)

		switch rf.Type.Kind() {
		case reflect.Bool:
			ref := (*bool)(unsafe.Pointer(fv.UnsafeAddr()))
			fs.BoolVarP(ref, fieldname, shorthand, fv.Bool(), helpText)
		case reflect.Float64:
			ref := (*float64)(unsafe.Pointer(fv.UnsafeAddr()))
			fs.Float64VarP(ref, fieldname, shorthand, fv.Float(), helpText)
		case reflect.Int64:
			switch rf.Type {
			case rTimeDuration:
				ref := (*time.Duration)(unsafe.Pointer(fv.UnsafeAddr()))
				fs.DurationVarP(ref, fieldname, shorthand, time.Duration(fv.Int()), helpText)
			default:
				ref := (*int64)(unsafe.Pointer(fv.UnsafeAddr()))
				fs.Int64VarP(ref, fieldname, shorthand, fv.Int(), helpText)
			}
		case reflect.Int:
			ref := (*int)(unsafe.Pointer(fv.UnsafeAddr()))
			fs.IntVarP(ref, fieldname, shorthand, int(fv.Int()), helpText)
		case reflect.String:
			ref := (*string)(unsafe.Pointer(fv.UnsafeAddr()))
			fs.StringVarP(ref, fieldname, shorthand, fv.String(), helpText)
		case reflect.Uint64:
			ref := (*uint64)(unsafe.Pointer(fv.UnsafeAddr()))
			fs.Uint64VarP(ref, fieldname, shorthand, fv.Uint(), helpText)
		case reflect.Uint:
			ref := (*uint)(unsafe.Pointer(fv.UnsafeAddr()))
			fs.UintVarP(ref, fieldname, shorthand, uint(fv.Uint()), helpText)
		default:
			// TODO: map
			panic(fmt.Sprintf("unsupported type %v", rf.Type))
		}
	}
	return &FlagSet{FlagSet: fs, builder: b}
}
