package wasm

import (
	elutils "github.com/rosbit/go-embedding-utils"
	"syscall/js"
	"reflect"
	"strings"
	// "fmt"
)

func bindGoStruct(structVar reflect.Value) (goStruct map[string]interface{}) {
	var structE reflect.Value
	if structVar.Kind() == reflect.Ptr {
		structE = structVar.Elem()
	} else {
		structE = structVar
	}
	structT := structE.Type()

	if structE == structVar {
		// struct is unaddressable, so make a copy of struct to an Elem of struct-pointer.
		// NOTE: changes of the copied struct cannot effect the original one. it is recommended to use the pointer of struct.
		structVar = reflect.New(structT) // make a struct pointer
		structVar.Elem().Set(structE)    // copy the old struct
		structE = structVar.Elem()       // structE is the copied struct
	}

	goStruct = wrapGoStruct(structVar, structE, structT)
	return
}

func lowerFirst(name string) string {
	return strings.ToLower(name[:1]) + name[1:]
}

func wrapGoStruct(structVar, structE reflect.Value, structT reflect.Type) map[string]interface{} {
	r := make(map[string]interface{})
	for i:=0; i<structT.NumField(); i++ {
		strField := structT.Field(i)
		fv := structE.Field(i)
		if !fv.CanInterface() {
			continue
		}
		name := strField.Name
		name = lowerFirst(name)
		r[name] = toValue(fv.Interface())
	}

	// receiver is the struct
	bindGoMethod(structE, structT, r)

	// reciver is the pointer of struct
	t := structVar.Type()
	bindGoMethod(structVar, t, r)
	return r
}

func bindGoMethod(structV reflect.Value, structT reflect.Type, r map[string]interface{}) {
	for i := 0; i<structV.NumMethod(); i+=1 {
		m := structT.Method(i)
		name := lowerFirst(m.Name)
		mV := structV.Method(i)
		mT := mV.Type()
		r[name] = js.FuncOf(wrapGoFunc(elutils.NewGolangFuncHelperDiretly(mV, mT)))
	}
}
