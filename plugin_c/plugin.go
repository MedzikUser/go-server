package main

// #include "foo.h"
import "C"

import (
	"github.com/MedzikUser/go-server/types"
)

// plugin name
var PluginName string = "c-lib"

// plugin type
var Type string = "command"

// command after which the `F` function will be executed
var Command string = "/c"

// the help message shown after typing /help
var Help string = "command from C lib"

// plugin main function must be named `F`
func F(input string, client types.Client) {
	out := C.foo(C.CString(input))
	client.Send(C.GoString(out))
}
