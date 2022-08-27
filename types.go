package StrCmd

type App struct {
	Commands map[string]Command
	Args     []CommandArgs
}

type Command struct {
	Subcommands map[string]Command
	HasSub      bool
	Args        []string
	NeedArgs    bool
	Action      func()
	args        []GennedArgs
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
