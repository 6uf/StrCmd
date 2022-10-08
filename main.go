package StrCmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var AppConfig *App

func GetAllNames(Text string) []string {
	return strings.Split(Text, " ")
}

func (Data *App) ParseCommand(Text string) error {
	AppConfig = Data
	var Args []CommandArgs
	var Default Command
	var Names string
	for _, names := range GetAllNames(Text) {
		if d, ok := Data.Commands[names]; ok {
			Default = d
			Names = names
		}
	}
	if len(Default.Args) > 0 {
		for _, Args := range Default.Args {
			if strings.Contains(Args, "--") && strings.Contains(Text, Args) {
				Default.args = append(Default.args, GennedArgs{
					Name:   Args,
					Value:  "true",
					IsBool: true,
				})
			} else {
				if Name, Value := GetKey(Args, Text); Name != "" && Value != "" {
					Default.args = append(Default.args, GennedArgs{
						Name:   Name,
						Value:  Value,
						IsBool: strings.Contains(Value, "--"),
					})
				}
			}
		}
		Args = append(Args, CommandArgs{
			Name: Names,
			Args: Default.args,
		})
	}

	var UsingSubCmd bool
	if Default.Subcommand != nil {
		var UptoDate SubCmd
		for _, names := range GetAllNames(Text) {
			if d, ok := Default.Subcommand[names]; ok {
				UptoDate = d
				UsingSubCmd = true
				break
			}
		}

		if UsingSubCmd {
			if len(UptoDate.Args) > 0 {
				for _, Args := range UptoDate.Args {
					if strings.Contains(Args, "--") && strings.Contains(Text, Args) {
						Default.args = append(Default.args, GennedArgs{
							Name:   Args,
							Value:  "true",
							IsBool: true,
						})
					} else {
						if Name, Value := GetKey(Args, Text); Name != "" && Value != "" {
							Default.args = append(Default.args, GennedArgs{
								Name:   Name,
								Value:  Value,
								IsBool: strings.Contains(Value, "--"),
							})
						}
					}
				}
				Args = append(Args, CommandArgs{
					Name: Names,
					Args: UptoDate.args,
				})
			}
			Data.Args = Args
			if UptoDate.Action != nil {
				UptoDate.Action()
			}
		}
	}

	if !UsingSubCmd {
		Data.Args = Args
		if Default.Action != nil {
			Default.Action()
		}
	}

	Data.Args = []CommandArgs{}
	return nil
}

func GetKey(Arg, Text string) (string, string) {
	if Data := regexp.MustCompile(fmt.Sprintf(`%v ([a-zA-Z0-9.>,<?'`+"`"+`";:[{}=_*&^%$#@!~-]+)`, Arg)).FindAllStringSubmatch(Text, 1); len(Data) == 1 {
		return Arg, Data[0][1]
	}
	return "", ""
}

func (D *App) Run(inputtext string) {
	for {
		D.ParseCommand(Listen(true, inputtext))
	}
}

func Listen(show bool, input string) string {
	fmt.Print(input)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func GetValue(Arg string) string {
	for _, arg := range AppConfig.Args {
		for _, arg := range arg.Args {
			if arg.Name == Arg {
				return arg.Value
			}
		}
	}
	return ""
}

func GetInt(Arg string) int {
	for _, arg := range AppConfig.Args {
		for _, arg := range arg.Args {
			if arg.Name == Arg {
				if value, err := strconv.Atoi(arg.Value); err == nil {
					return value
				}
			}
		}
	}
	return 0
}

func GetBool(Arg string) bool {
	for _, arg := range AppConfig.Args {
		for _, arg := range arg.Args {
			if arg.Name == Arg && arg.IsBool {
				return true
			}
		}
	}
	return false
}
