package wasm

import (
	"reflect"
)

func bindGoInterface(interfaceVar reflect.Value) (goInterface map[string]interface{}) {
	t := interfaceVar.Type()
	r := make(map[string]interface{})
	bindGoMethod(interfaceVar, t, r)
	return r
}

