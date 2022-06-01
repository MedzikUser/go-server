package plugin

import (
	"fmt"
	"log"
	"plugin"

	"github.com/MedzikUser/go-server/types"
)

// function that loads a plugin from `.so`
func LoadPlugins(plugins []string) []types.Plugin {
	o := []types.Plugin{}

	for _, p := range plugins {
		// open the `.so` plugin so that things can be extracted from it
		p, err := plugin.Open(p)
		if err != nil {
			log.Fatal(err)
		}

		// lockup `PluginName` in plugin
		pName, err := p.Lookup("PluginName")
		if err != nil {
			panic(err)
		}

		// lookup `Type` in plugin
		pType, err := p.Lookup("Type")
		if err != nil {
			panic(err)
		}

		pluginType := pType.(*string)

		if *pluginType == "command" {
			// lookup `Command` in plugin
			pCommand, err := p.Lookup("Command")
			if err != nil {
				panic(err)
			}

			// lookup `Help` in plugin
			pHelp, err := p.Lookup("Help")
			if err != nil {
				panic(err)
			}

			// lookup `F` in plugin (main function)
			f, err := p.Lookup("F")
			if err != nil {
				panic(err)
			}

			plugin := types.Plugin{
				Name:    *pName.(*string),
				Command: *pCommand.(*string),
				Help:    *pHelp.(*string),
				F:       f.(func(string, types.Client)),
			}

			o = append(o, plugin)
		} else if *pluginType == "event" {
			// lookup `Help` in plugin
			pEvent, err := p.Lookup("Event")
			if err != nil {
				panic(err)
			}

			// lookup `F` in plugin (main function)
			f, err := p.Lookup("F")
			if err != nil {
				panic(err)
			}

			plugin := types.Plugin{
				Name:  *pName.(*string),
				Event: *pEvent.(*string),
				F:     f.(func(string, types.Client)),
			}

			o = append(o, plugin)
		}
	}

	fmt.Printf("Plugins (%d):\n", len(o))
	for i, plugin := range o {
		fmt.Printf("(%d) %s\n", i+1, plugin.Name)
	}

	return o
}
