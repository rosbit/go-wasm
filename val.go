package wasm

import (
	"syscall/js"
	"reflect"
	"encoding/json"
)

func toValue(v interface{}) interface{} {
	if v == nil {
		return undefined
	}

	switch vv := v.(type) {
	case int,int8,int16,int32,int64,
		uint,uint8,uint16,uint32,uint64,
		float32,float64,
		string,bool:
		return v
	case []byte:
		return string(vv)
	case js.Value:
		return vv
	default:
		v2 := reflect.ValueOf(v)
		switch v2.Kind() {
		case reflect.Slice, reflect.Array:
			r := make([]interface{}, v2.Len())
			for i:=0; i<v2.Len(); i++ {
				r[i] = toValue(v2.Index(i).Interface())
			}
			return r
		case reflect.Map:
			r := make(map[interface{}]interface{})
			iter := v2.MapRange()
			for iter.Next() {
				k, v1 := iter.Key(), iter.Value()
				r[toValue(k.Interface())] = toValue(v1.Interface())
			}
			return r
		case reflect.Struct:
			return bindGoStruct(v2)
		case reflect.Ptr:
			e := v2.Elem()
			if e.Kind() == reflect.Struct {
				return bindGoStruct(v2)
			}
			return toValue(e.Interface())
		case reflect.Func:
			if f, err := bindGoFunc(v); err == nil {
				return f
			}
			return undefined
		case reflect.Interface:
			return bindGoInterface(v2)
		default:
			return null
		}
	}
}

func fromValue(v js.Value) (interface{}) {
	switch v.Type() {
	case js.TypeUndefined:
		return nil
	case js.TypeNull:
		return nil
	case js.TypeBoolean:
		return v.Bool()
	case js.TypeNumber:
		return v.Float()
	case js.TypeString:
		return v.String()
	case js.TypeSymbol:
		return v.String()
	case js.TypeObject:
		if v.InstanceOf(array) {
			l := v.Length()
			res := make([]interface{}, l)
			for i:=0; i<l; i++ {
				res[i] = fromValue(v.Index(i))
			}
			return res
		}
		obj := jsJSON.Call("stringify", v).String()
		res := make(map[string]interface{})
		json.Unmarshal([]byte(obj), &res)
		return res
	case js.TypeFunction:
		return v
	default:
		return nil
	}
}

