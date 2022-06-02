package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/MedzikUser/go-server/plugin"
	"github.com/MedzikUser/go-server/types"
)

var ClientManager []types.Client

func main() {
	var port string
	flag.StringVar(&port, "p", "9999", "TCP server port")
	flag.Parse()

	// folder where the `.so` plugin files are located
	pluginDirectory := "plugins"

	outputDirRead, err := os.Open(pluginDirectory)
	if err != nil {
		panic(err)
	}

	outputDirFiles, err := outputDirRead.Readdir(0)
	if err != nil {
		panic(err)
	}

	var plugins_files []string

	// one by one, add the plugin path to the `plugins_files` variable
	for outputIndex := range outputDirFiles {
		outputFileHere := outputDirFiles[outputIndex]
		plugins_files = append(plugins_files, fmt.Sprintf("%s/%s", pluginDirectory, outputFileHere.Name()))
	}

	// load plugins
	plugins := plugin.LoadPlugins(plugins_files)

	// listen TCP server
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		panic("unexcepted error when tcp server starting: " + err.Error())
	}

	println("TCP server started!")

	for {
		// when the user connects
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept Error", err)
		}

		// handle a client in gorouitine
		go func() {
			log.Println("New Client", conn.RemoteAddr())

			client := types.Client{Conn: conn}

			ClientManager = append(ClientManager, client)

			handleConnection(client, plugins)
		}()
	}
}

// handle client connection
func handleConnection(client types.Client, plugins []types.Plugin) {
	// finally terminate the connection
	defer client.Conn.Close()

	// run events from plugins (onConnect)
	for _, plugin := range plugins {
		if plugin.Event == "onConnect" {
			plugin.F("", client)
		}
	}

	// loop to avoid terminating the connection after one command
	for {
		buf := make([]byte, 1024)

		reqLen, err := client.Conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading buffer:", err.Error())
			break
		}

		// run events from plugins (onSend)
		for _, plugin := range plugins {
			if plugin.Event == "onSend" {
				plugin.F("", client)
			}
		}

		handleCommand(string(buf[0:reqLen]), client, plugins)
	}
}

// handle client command
func handleCommand(input string, client types.Client, plugins []types.Plugin) {
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\r", "")
	str := strings.Split(input, " ")

	if len(str) <= 0 {
		client.Send("empty buffer")
		return
	}

	command := str[0]

	executed := false

	// search for the command in plugins
	for _, plugin := range plugins {
		if command == plugin.Command {
			plugin.F(input, client)

			executed = true

			break
		}
	}

	// commands that are not from plugins
	if !executed {
		switch command {
		case "/help", "/h":
			var str string
			str = str + "/help - show help\n"
			str = str + "/plugins - list all plugins\n"
			str = str + "/broadcast <message> - send message to all clients"
			str = str + "/disconnect - close connection"

			for _, plugin := range plugins {
				if len(plugin.Command) <= 0 {
					continue
				}

				str = fmt.Sprintf("%s\n%s - %s", str, plugin.Command, plugin.Help)
			}

			client.Send(str)

		case "/plugins", "/pl":
			var str string = "Plugins:"

			for i, plugin := range plugins {
				str = fmt.Sprintf("%s\n(%d) %s", str, i+1, plugin.Name)
			}

			client.Send(str)

		case "/clients":
			var str string

			for i, client := range ClientManager {
				str = fmt.Sprintf("%sip=%s id=%d\n", str, client.Conn.RemoteAddr(), i+1)
			}

			client.Send(str)

		case "/broadcast":
			if len(str) <= 1 {
				client.Send("Usage: /broadcast <message>")
			}

			for _, client := range ClientManager {
				client.Send(strings.Join(str[1:], " "))
			}

		case "/disconnect", "/exit", "/close":
			client.Conn.Close()

		default:
			client.Send("unknown command")
		}
	}
}
