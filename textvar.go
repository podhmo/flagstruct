package flagstruct

import (
	"encoding"
	"fmt"
	"reflect"
)

// copy from go1.19's flag package

// func (f *FlagSet) TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string) {
// 	f.Var(newTextValue(value, p), name, usage)
// }

type textValue struct {
	p     encoding.TextUnmarshaler
	isPtr bool
}

func newTextValue(val encoding.TextMarshaler, p encoding.TextUnmarshaler, isPtr bool) textValue {
	ptrVal := reflect.ValueOf(p)
	if ptrVal.Kind() != reflect.Ptr {
		panic("variable value type must be a pointer")
	}
	defVal := reflect.ValueOf(val)
	if defVal.Kind() == reflect.Ptr {
		defVal = defVal.Elem()
	}
	if defVal.Type() != ptrVal.Type().Elem() {
		panic(fmt.Sprintf("default type does not match variable type: %v != %v", defVal.Type(), ptrVal.Type().Elem()))
	}
	ptrVal.Elem().Set(defVal)
	return textValue{p, isPtr}
}

func (v textValue) Set(s string) error {
	return v.p.UnmarshalText([]byte(s))
}

func (v textValue) Get() interface{} {
	return v.p
}

func (v textValue) String() string {
	if m, ok := v.p.(encoding.TextMarshaler); ok {
		if b, err := m.MarshalText(); err == nil {
			return string(b)
		}
	}
	return ""
}

// for pflag.Value
func (v textValue) Type() string {
	if !v.isPtr {
		return fmt.Sprintf("%T", v.p)[1:]
	}
	return fmt.Sprintf("%T", v.p)
}
