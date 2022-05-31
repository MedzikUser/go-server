package plugin

import (
	"fmt"
	"log"
	"plugin"

	"github.com/MedzikUser/go-server/types"
)

func LoadPlugins(plugins []string) []types.Plugin {
	o := []types.Plugin{}

	for _, p := range plugins {
		p, err := plugin.Open(p)
		if err != nil {
			log.Fatal(err)
		}

		pName, err := p.Lookup("PluginName")
		if err != nil {
			panic(err)
		}

		pCommand, err := p.Lookup("Command")
		if err != nil {
			panic(err)
		}

		pHelp, err := p.Lookup("Help")
		if err != nil {
			panic(err)
		}

		f, err := p.Lookup("F")
		if err != nil {
			panic(err)
		}

		F := f.(func(string, types.Client))

		pl := types.Plugin{
			Name:    pName.(*string),
			Command: pCommand.(*string),
			Help:    pHelp.(*string),
			F:       F,
		}

		o = append(o, pl)
	}

	fmt.Printf("Plugins (%d):\n", len(o))
	for i, plugin := range o {
		fmt.Printf("(%d) %s\n", i+1, *plugin.Name)
	}

	return o
}
