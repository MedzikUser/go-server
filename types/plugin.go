package types

type Plugin struct {
	Name    string
	Event   string // plugin events (`onConnect`)
	Command string
	Help    string
	F       func(string, Client)
}
