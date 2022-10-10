package StrCmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var Current CommandArgs

func (Data *App) ParseCommand(Text string) error {
	var Args CommandArgs
	var Default Command
	var Names string
	var ParsedNames []string = strings.Split(Text, " ")
	var GennedArg = []GennedArgs{}

	for _, names := range ParsedNames {
		if d, ok := Data.Commands[names]; ok {
			Default = d
			Names = names
		}
	}
	if len(Default.Args) > 0 {
		for _, Args := range Default.Args {
			if strings.Contains(Args, "--") && strings.Contains(Text, Args) {
				GennedArg = append(GennedArg, GennedArgs{
					Name:   Args,
					Value:  "true",
					IsBool: true,
				})
			} else {
				if Name, Value := GetKey(Args, Text); Name != "" && Value != "" {
					GennedArg = append(GennedArg, GennedArgs{
						Name:   Name,
						Value:  Value,
						IsBool: strings.Contains(Value, "--"),
					})
				}
			}
		}
		Args = CommandArgs{
			Name: Names,
			Args: GennedArg,
		}
	}

	var UsingSub bool

	if Default.Subcommand != nil {
		var UptoDate SubCmd
		for _, names := range ParsedNames {
			if d, ok := Default.Subcommand[names]; ok {
				UptoDate = d
				UsingSub = true
				Names = names
				break
			}
		}
		if UsingSub {
			if len(UptoDate.Args) > 0 {
				for _, Args := range UptoDate.Args {
					if strings.Contains(Args, "--") && strings.Contains(Text, Args) {
						GennedArg = append(GennedArg, GennedArgs{
							Name:   Args,
							Value:  "true",
							IsBool: true,
						})
					} else {
						if Name, Value := GetKey(Args, Text); Name != "" && Value != "" {
							GennedArg = append(GennedArg, GennedArgs{
								Name:   Name,
								Value:  Value,
								IsBool: strings.Contains(Value, "--"),
							})
						}
					}
				}

				Args = CommandArgs{
					Name: Names,
					Args: GennedArg,
				}
			}
		}

		Current = Args
		if UptoDate.Action != nil {
			UptoDate.Action()
		}
	}

	if !UsingSub {
		if !Data.DontUseBuiltinHelpCmd {
			if strings.HasPrefix(Text, "help") {
				fmt.Println(Data.FormatHelpText())
			}
		}
		Current = Args
		if Default.Action != nil {
			Default.Action()
		}
	}

	return nil
}

func (Data *App) FormatHelpText() (Base string) {
	if Data.Version != 0 {
		Base += fmt.Sprintf("VERSION: %v\n\n", Data.Version)
	} else {
		Base += "VERSION: 1.0.0\n\n"
	}
	if Data.AppDescription != "" {
		Base += "Description: " + Data.AppDescription + "\n\n"
	}
	Base += ReturnCommandInfo(Data.Commands, " [ARGS")
	return
}

func ReturnCommandInfo(Value map[string]Command, Format string) (Base string) {
	for name, key := range Value {
		if key.Description == "" {
			key.Description = "A global command that is parsed through StrCmd (Description was empty!)"
		}

		var B string = Format
		var S string = "    [SUBCMD]\n"
		if len(key.Args) > 0 {
			for _, name := range key.Args {
				B += " " + name
			}
			B += "]"
		}

		for name, key := range key.Subcommand {
			if key.Description == "" {
				key.Description = "A global command that is parsed through StrCmd (Description was empty!)"
			}

			var B string = Format
			if len(key.Args) > 0 {
				for _, name := range key.Args {
					B += " " + name
				}
				B += "]"
			}

			if B != Format {
				S += fmt.Sprintf("     - %v | %v\n%v\n", name, key.Description, B)
			} else {
				S += fmt.Sprintf("     - %v | %v\n", name, key.Description)
			}
		}

		if B != Format && S != "    [SUBCMD]\n" {
			Base += fmt.Sprintf("  - %v | %v%v\n%v\n", name, key.Description, B, S)
		} else {
			if S != "    [SUBCMD]\n" {
				Base += fmt.Sprintf("  - %v | %v\n%v\n", name, key.Description, S)
			} else {
				Base += fmt.Sprintf("  - %v | %v\n", name, key.Description)
			}
		}
	}
	return
}

func ReturnSubCmdInfo() (Base string) {
	return
}

func GetKey(Arg, Text string) (string, string) {
	if Data := regexp.MustCompile(fmt.Sprintf(`%v ([a-zA-Z0-9.>,<?'`+"`"+`";:[{}=_*&^%$#@!~-]+)`, Arg)).FindAllStringSubmatch(Text, 1); len(Data) == 1 {
		return Arg, Data[0][1]
	}
	return "", ""
}

func (D *App) Run(inputtext string) {
	for {
		if err := D.ParseCommand(Listen(true, inputtext)); err != nil {
			panic(err)
		}
	}
}

func Listen(show bool, input string) string {
	fmt.Print(input)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func GetValue(Arg string) string {
	for _, arg := range Current.Args {
		if arg.Name == Arg {
			return arg.Value
		}
	}
	return ""
}

func GetInt(Arg string) int {
	for _, arg := range Current.Args {
		if arg.Name == Arg {
			if value, err := strconv.Atoi(arg.Value); err == nil {
				return value
			}
		}
	}
	return 0
}

func GetBool(Arg string) bool {
	for _, arg := range Current.Args {
		if arg.Name == Arg && arg.IsBool {
			return true
		}
	}
	return false
}
