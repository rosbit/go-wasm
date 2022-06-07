package wasm

import (
	elutils "github.com/rosbit/go-embedding-utils"
	"syscall/js"
	"reflect"
)

func bindFunc(fn *js.Value, funcVarPtr interface{}) (err error) {
	helper, e := elutils.NewEmbeddingFuncHelper(funcVarPtr)
	if e != nil {
		err = e
		return
	}
	helper.BindEmbeddingFunc(wrapFunc(fn, helper))
	return
}

func wrapFunc(fn *js.Value, helper *elutils.EmbeddingFuncHelper) elutils.FnGoFunc {
	return func(args []reflect.Value) (results []reflect.Value) {
		var jsArgs []interface{}

		// make js args
		itArgs := helper.MakeGoFuncArgs(args)
		for arg := range itArgs {
			jsArgs = append(jsArgs, arg)
		}

		// call js function
		res := fn.Invoke(jsArgs...)

		// convert result to golang
		results = helper.ToGolangResults(fromValue(res), res.Type() == js.TypeObject && res.InstanceOf(array), nil)
		return
	}
}
