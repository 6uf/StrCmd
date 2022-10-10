package StrCmd

type App struct {
	Commands              map[string]Command
	Version               float64
	AppDescription        string
	DontUseBuiltinHelpCmd bool
	Args                  []CommandArgs
}

type Command struct {
	Description string
	Subcommand  map[string]SubCmd
	Args        []string
	Action      func()
}

type SubCmd struct {
	Args   []string
	Action func()
}

type CommandArgs struct {
	Name string
	Args []GennedArgs
}

type GennedArgs struct {
	Name   string
	Value  string
	IsBool bool
}
