# go-wasm, makes creating WebAssembly with golang easily

`go-wasm` is a package extending the golang package `syscall/js` and making creating WebAssembly with golang easily

### Usage

The package is fully go-getable, so, just type

  `go get github.com/rosbit/go-wasm`

to install.

All the following Wasm result should be built with
  `GOOS=js GOARCH=wasm go build`

#### 1. Evaluate expressions

```go
package main

import (
  "github.com/rosbit/go-wasm"
  "fmt"
)

func main() {
  res, _ := wasm.Eval("1 + 2", nil)
  fmt.Println("result is:", res)
}
```

#### 2. Go calls JS function

Suppose there's a JS file named `a.js` like this:

```js
function add(a, b) {
    return a+b
}
```

one can call the JS function `add()` in Go code like the following:

```go
package main

import (
  "github.com/rosbit/go-wasm"
  "fmt"
)

var add func(int, int)int

func main() {
  if err := wasm.BindFunc("add", &add); err != nil {
     fmt.Printf("%v\n", err)
     return
  }

  res := add(1, 2)
  fmt.Println("result is:", res)
}
```

#### 3. JS calls Go function

JS calling Go function is also easy. In the Go code, make a Golang function
as JS built-in function by calling `MakeJsFunc("funcname", function)`. There's the example:

```go
package main

import "github.com/rosbit/go-wasm"

// function to be called by JS
func adder(a1 float64, a2 float64) float64 {
    return a1 + a2
}

func main() {
  wasm.MakeJsFunc("adder", adder) // now "adder" is a global JS function
  done := make(chan struct{})
  <-done
}
```

In JS code, one can call the registered function directly. There's the example `b.js`.

```js
r = adder(1, 100)   # the function "adder" is implemented in Go
console.log(r)
```

#### 4. Make Go struct instance as a JS object

This package provides a function `SetModule` which will convert a Go struct instance into
a JS object. There's the example `c.js`, `m` is the object var provided by Go code:

```js
m.incAge(10)
print(m)

console.log('m.name', m.name)
console.log('m.age', m.age)
```

The Go code is like this:

```go
package main

import "github.com/rosbit/go-wasm"

type M struct {
   Name string
   Age int
}
func (m *M) IncAge(a int) {
   m.Age += a
}

func main() {
  wasm.BindObject("m", &M{Name:"rosbit", Age: 1}) // "m" is the object var name

  done := make(chan struct{})
  <-done
}
```

#### 5. Set many built-in functions and objects at one time

If there're a lot of functions and objects to be registered, a map could be constructed with function `SetGlobals` or put as an
argument for function `Eval`.

```go
package main

import "github.com/rosbit/go-wasm"
import "fmt"

type M struct {
   Name string
   Age int
}
func (m *M) IncAge(a int) {
   m.Age += a
}

func adder(a1 float64, a2 float64) float64 {
    return a1 + a2
}

func main() {
  vars := map[string]interface{}{
     "m": &M{Name:"rosbit", Age:1}, // to JS object
     "adder": adder,                // to JS built-in function
     "a": []int{1,2,3},             // to JS array
  }

  wasm.SetGlobals(vars)
  res := wasm.GetGlobal("a") // get the value of var named "a". Any variables in script could be get by GetGlobal
  fmt.Printf("res:", res)
}
```

#### 6. Wrap anything as JS global object

This package also provides a function `SetGlobalObject` which will create a JS variable integrating any
Go values/functions as an object. There's the example `d.js` which will use object `tm` provided by Go code:

```js
a = tm.newA("rosbit", 10)
a.incAge(10)
console.log(a)

tm.printf('a.name: %s\n', a.name)
tm.printf('a.age: %d\n', a.age)
```

The Go code is like this:

```go
package main

import (
  "github.com/rosbit/go-wasm"
  "fmt"
)

type A struct {
   Name string
   Age int
}
func (m *A) IncAge(a int) {
   m.Age += a
}
func newA(name string, age int) *A {
   return &A{Name: name, Age: age}
}

func main() {
  wasm.SetGlobalObject("tm", map[string]interface{}{ // object name is "tm"
     "newA": newA,            // make user defined function as object method named "tm.newA"
     "printf": fmt.Printf,    // make function in a standard package named "tm.printf"
  })

  done := make(chan struct{})
  <-done
}
```

### Other helper functions

```go
func JSONStringify(v interface{}) string    // wasm.JSONStringify([]int{1, 3})
func JSONParse(val string) interface{}      // wasm.JSONParse(`{"a":"b","c":1}`)
func CallObjectMethod(objName, objMethod string, args ...interface{}) js.Value
 // jsObj := wasm.CallObjectMethod("document", "getElementById", "id")
func CallFunc(funcName string, args ...interface{}) js.Value // wasm.CallFunc("add", 1, 1)
```

### Status

The package is not fully tested, so be careful.

### Contribution

Pull requests are welcome! Also, if you want to discuss something send a pull request with proposal and changes.

__Convention:__ fork the repository and make changes on your fork in a feature branch.
