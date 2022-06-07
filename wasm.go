package wasm

import (
	"syscall/js"
	"fmt"
	"reflect"
)

var (
	global js.Value
	array  js.Value
	object js.Value
	jsJSON js.Value
	eval   js.Value
	null   js.Value
	undefined js.Value
)

func init() {
	global = js.Global()
	array  = global.Get("Array")
	object = global.Get("Object")
	jsJSON = global.Get("JSON")
	eval   = global.Get("eval")
	null   = js.Null()
	undefined = js.Undefined()
}

func Eval(script string, vars map[string]interface{}) (res interface{}) {
	SetGlobals(vars)
	r := eval.Invoke(script)
	res = fromValue(r)
	return
}

// make a golang pointer of sturct instance as a JS object.
// @param structVarPtr  pointer of struct instance is recommended.
func BindObject(jsObjName string, structVarPtr interface{}) (err error) {
	if structVarPtr == nil {
		err = fmt.Errorf("structVarPtr must ba non-nil pointer of struct")
		return
	}
	v := reflect.ValueOf(structVarPtr)
	if v.Kind() == reflect.Struct || (v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct) {
		goStruct := bindGoStruct(v)
		global.Set(jsObjName, goStruct)
		return
	}
	err = fmt.Errorf("structVarPtr must be struct or pointer of strcut")
	return
}

// bind a var of golang func with a JS function name, so calling JS function
// is just calling the related golang func.
// @param funcVarPtr  in format `var funcVar func(....) ...; funcVarPtr = &funcVar`
func BindFunc(jsFnName string, funcVarPtr interface{}) (err error) {
	if funcVarPtr == nil {
		err = fmt.Errorf("funcVarPtr must be a non-nil poiter of func")
		return
	}
	t := reflect.TypeOf(funcVarPtr)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Func {
		err = fmt.Errorf("funcVarPtr expected to be a pointer of func")
		return
	}

	v := global.Get(jsFnName)
	if v.Type() != js.TypeFunction {
		err = fmt.Errorf("var %s is not with type function", jsFnName)
		return
	}
	bindFunc(&v, funcVarPtr)
	return
}

// make a golang func as a built-in JS function, so the function can be called in JS script.
func MakeJsFunc(fnName string, funcVar interface{}) (err error) {
	goFunc, e := bindGoFunc(funcVar)
	if e != nil {
		err = e
		return
	}
	global.Set(fnName, goFunc)
	return
}

func SetGlobals(vars map[string]interface{}) {
	for k, v := range vars {
		if len(k) == 0 {
			continue
		}
		global.Set(k, toValue(v))
	}
}

func GetGlobal(name string) (res interface{}) {
	v := global.Get(name)
	res = fromValue(v)
	return
}

func SetGlobalObject(varName string, name2vals map[string]interface{}) (err error) {
	if len(varName) == 0 {
		err = fmt.Errorf("varName expected")
		return
	}
	obj := convertMap(name2vals)
	global.Set(varName, obj)
	return
}

func JSONStringify(v interface{}) string {
	return jsJSON.Call("stringify", toValue(v)).String()
}

func JSONParse(val string) interface{} {
	return fromValue(jsJSON.Call("parse", val))
}

func CallObjectMethod(objName, objMethod string, args ...interface{}) js.Value {
	obj := global.Get(objName)
	if !obj.Truthy() {
		return undefined
	}
	return obj.Call(objMethod, args...)
}

func CallFunc(funcName string, args ...interface{}) js.Value {
	fn := global.Get(funcName)
	if fn.Type() != js.TypeFunction {
		return undefined
	}
	return fn.Invoke(args...)
}

func convertMap(vars map[string]interface{}) (map[string]interface{}) {
	if len(vars) == 0 {
		return nil
	}
	res := make(map[string]interface{})
	for k, v := range vars {
		if len(k) == 0 {
			continue
		}
		res[k] = toValue(v)
	}
	return res
}
