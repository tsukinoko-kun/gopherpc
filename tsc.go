package gopherpc

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

var (
	timeType = reflect.TypeOf(time.Time{})
)

func goTypeToTypescript(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Invalid:
		return "never"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Uint64, reflect.Int64:
		return "number"
	case reflect.Array, reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Uint8, reflect.Int8:
			return "string"
		}
		return fmt.Sprintf("Array<%s>", goTypeToTypescript(t.Elem()))
	case reflect.Chan, reflect.Func:
		return "unknown"
	case reflect.Interface:
		return "any"
	case reflect.Map:
		return fmt.Sprintf("Record<%s, %s>", goTypeToTypescript(t.Key()), goTypeToTypescript(t.Elem()))
	case reflect.Pointer, reflect.UnsafePointer:
		return fmt.Sprintf("(null|%s)", goTypeToTypescript(t.Elem()))
	case reflect.String:
		return "string"
	case reflect.Struct:
		if t == timeType {
			return "string"
		}
		s := strings.Builder{}
		s.WriteString("{ ")
		numFields := t.NumField()
		for i := 0; i < numFields; i++ {
			if i != 0 {
				s.WriteString("; ")
			}
			field := t.Field(i)
			if jsonTag, ok := field.Tag.Lookup("json"); ok {
				if jsonTag == "-" || jsonTag == "" {
					s.WriteString(fmt.Sprintf("_%d: any", i))
					continue
				}
				s.WriteString(fmt.Sprintf("%q: ", strings.Split(jsonTag, ",")[0]))
			} else {
				s.WriteString(fmt.Sprintf("%q: ", field.Name))
			}
			s.WriteString(goTypeToTypescript(field.Type))
		}
		s.WriteString(" }")
		return s.String()
	default:
		return "any"
	}
}
