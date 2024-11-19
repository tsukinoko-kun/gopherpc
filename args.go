package gopherpc

import (
	"fmt"
	"reflect"
	"strconv"
)

// Parse args from GopheRPC registered function call into struct
func Unmarshal(args []any, v any) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("v must be a pointer")
	}
	for i := 0; i < reflect.ValueOf(v).Elem().NumField(); i++ {
		field := reflect.ValueOf(v).Elem().Field(i)
		tag := reflect.TypeOf(v).Elem().Field(i).Tag.Get("index")
		var index int
		if tag == "-" {
			continue
		} else if tag == "" {
			index = i
		} else {
			var err error
			index, err = strconv.Atoi(tag)
			if err != nil {
				return err
			}
		}
		if index >= len(args) {
			return fmt.Errorf("index out of range")
		}

		// is right type?
		if reflect.TypeOf(args[index]) == field.Type() {
			field.Set(reflect.ValueOf(args[index]))
		} else {
			// convert to right type
			switch field.Kind() {
			case reflect.String:
				field.SetString(fmt.Sprintf("%v", args[index]))
			case reflect.Int:
				field.SetInt(int64(args[index].(float64)))
			case reflect.Float64:
				field.SetFloat(args[index].(float64))
			case reflect.Float32:
				field.SetFloat(float64(args[index].(float32)))
			case reflect.Bool:
				field.SetBool(args[index].(bool))
			default:
				return fmt.Errorf("type not match")
			}
		}
	}
	return nil
}
