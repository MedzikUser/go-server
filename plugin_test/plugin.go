package main

import (
	"github.com/MedzikUser/go-server/types"
)

// plugin name
var PluginName string = "test"

// command after which the `F` function will be executed
var Command string = "/test"

// the help message shown after typing /help
var Help string = "test command from lib"

// plugin main function must be named `F`
func F(input string, client types.Client) {
	// send message to client
	client.Send("This executed from test plugin!")
}
