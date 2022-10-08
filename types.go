package StrCmd

type App struct {
	Commands map[string]Command
	Args     []CommandArgs
}

type Command struct {
	Subcommand SubCmd
	Args       []string
	Action     func()
	args       []GennedArgs
}

type SubCmd struct {
	Name   string
	Args   []string
	Action func()
	args   []GennedArgs
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
