package wasm

import (
	elutils "github.com/rosbit/go-embedding-utils"
	"syscall/js"
)

func bindGoFunc(funcVar interface{}) (goFunc js.Func, err error) {
	helper, e := elutils.NewGolangFuncHelper("noname", funcVar)
	if e != nil {
		err = e
		return
	}
	goFunc = js.FuncOf(wrapGoFunc(helper))
	return
}

func wrapGoFunc(helper *elutils.GolangFuncHelper) func(this js.Value, args []js.Value)interface{} {
	return func(this js.Value, args []js.Value) (val interface{}) {
		getArgs := func(i int) interface{} {
			return fromValue(args[i])
		}

		v, e := helper.CallGolangFunc(len(args), "gofunc", getArgs)
		if e != nil || v == nil {
			return
		}

		val = v
		return
	}
}
