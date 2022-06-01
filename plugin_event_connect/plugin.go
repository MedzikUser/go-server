package main

import (
	"fmt"

	"github.com/MedzikUser/go-server/types"
)

// plugin name
var PluginName string = "welcome"

// plugin type
var Type string = "event"

// event name (onConnect)
var Event string = "onConnect"

// plugin main function must be named `F`
func F(_empty string, client types.Client) {
	// send message to client
	client.Send(fmt.Sprintf("You connected from ip `%s` to server `%s`", client.Conn.RemoteAddr(), client.Conn.LocalAddr()))
}
