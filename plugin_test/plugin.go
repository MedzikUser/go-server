package main

import (
	"github.com/MedzikUser/go-server/types"
)

var PluginName string = "test"
var Command string = "/test"
var Help string = "test command from lib"

func F(input string, client types.Client) {
	client.Send("This executed from test plugin!")
}
