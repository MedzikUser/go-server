package main

import (
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
	directory := "plugins"

	outputDirRead, err := os.Open(directory)
	if err != nil {
		panic(err)
	}

	outputDirFiles, err := outputDirRead.Readdir(0)
	if err != nil {
		panic(err)
	}

	var plugins_files []string

	for outputIndex := range outputDirFiles {
		outputFileHere := outputDirFiles[outputIndex]
		plugins_files = append(plugins_files, fmt.Sprintf("%s/%s", directory, outputFileHere.Name()))
	}

	plugins := plugin.LoadPlugins(plugins_files)

	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		panic("unexcepted error when tcp server starting: " + err.Error())
	}

	println("TCP server started!")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept Error", err)
		}

		go func() {
			log.Println("Accepted", conn.RemoteAddr())

			client := types.Client{Conn: conn}

			ClientManager = append(ClientManager, client)

			handleConnection(client, plugins)
		}()
	}
}

func handleConnection(client types.Client, plugins []types.Plugin) {
	defer client.Conn.Close()

	for {
		buf := make([]byte, 1024)

		reqLen, err := client.Conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading buffer:", err.Error())
			break
		}

		handleCommand(string(buf[0:reqLen]), client, plugins)
	}
}

func handleCommand(input string, client types.Client, plugins []types.Plugin) {
	str := strings.Split(input, " ")

	if len(str) <= 0 {
		client.Send("empty buffer")
		return
	}

	command := str[0]

	executed := false

	for _, plugin := range plugins {
		if command == *plugin.Command {
			plugin.F(input, client)

			executed = true

			break
		}
	}

	if !executed {
		switch command {
		case "/help", "/h":
			var str string = "/help - show help\n"
			str = str + "/plugins - list all plugins\n"
			str = str + "/broadcast - list all plugins"

			for _, plugin := range plugins {
				str = fmt.Sprintf("%s\n%s - %s", str, *plugin.Command, *plugin.Help)
			}

			client.Send(str)

		case "/plugins", "/pl":
			var str string = "Plugins:"

			for i, plugin := range plugins {
				str = fmt.Sprintf("%s\n(%d) %s", str, i+1, *plugin.Name)
			}

			client.Send(str)

		case "/broadcast":
			for _, client := range ClientManager {
				client.Send(strings.Join(str[1:], ""))
			}

		default:
			client.Send("unknown command")
		}
	}
}
