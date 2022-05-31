package types

type Plugin struct {
	Name    *string
	Command *string
	Help    *string
	F       func(string, Client)
}
