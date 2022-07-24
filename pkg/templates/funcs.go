package templates

import (
	"reflect"
	tmpl "text/template"
)

func GetTempalteFuncs() map[string]any {
    return tmpl.FuncMap{
        "isString": func(i interface{}) bool {
            v := reflect.ValueOf(i)
            switch v.Kind() {
            case reflect.String:
                return true
            default:
                return false
            }
        },
        "isSlice": func(i interface{}) bool {
            v := reflect.ValueOf(i)
            switch v.Kind() {
            case reflect.Slice:
                return true
            default:
                return false
            }
        },
        "isArray": func(i interface{}) bool {
            v := reflect.ValueOf(i)
            switch v.Kind() {
            case reflect.Array:
                return true
            default:
                return false
            }
        },
        "isMap": func(i interface{}) bool {
            v := reflect.ValueOf(i)
            switch v.Kind() {
            case reflect.Map:
                return true
            default:
                return false
            }
        },
        "isList": func (i interface{}) bool {
            v := reflect.ValueOf(i).Kind()
            return v == reflect.Array || v == reflect.Slice
        },
        "isNumber": func (i interface{}) bool {
            v := reflect.ValueOf(i).Kind()
            switch v {
            case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
                return true
            default:
                return false
            }
        },
        "isInt": func (i interface{}) bool {
            v := reflect.ValueOf(i).Kind()
            switch v {
            case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
                return true
            default:
                return false
            }
        },
        "isFloat": func (i interface{}) bool {
            v := reflect.ValueOf(i).Kind()
            return v == reflect.Float32 || v == reflect.Float64
        },
    }
}
