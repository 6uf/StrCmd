# StrCmd
A simple input/output Command handler.

Example:
```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/6uf/StrCmd"
)

func main() {
	a := StrCmd.App{
		//DontUseBuiltinHelpCmd: true,
		Version:        "v2.0.0",
		AppDescription: "A command prompt demo of the library StrCmd",
		Commands: map[string]StrCmd.Command{
			"input": {
				Description: "Test inputs!",
				Action: func() {
					fmt.Println(StrCmd.String("-name"))
				},

				Args: []string{
					"-name",
				},
			},
			"subcmdtests": {
				Subcommand: map[string]StrCmd.SubCmd{
					"1": {
						Description: "First nested sub command of the command subcmdtests",
						Action: func() {
							fmt.Println("This is a basic function that prints this message within a nested sub command!")
						},
					},
					"2": {
						Description: "Second nested command of the command subcmdtests",
						Action: func() {
							fmt.Println("Third Layer", StrCmd.Int("-test"), StrCmd.String("-name"))
						},
						Args: []string{
							"-test",
							"-name",
						},
					},
				},
				Description: "This command is a normal command but also is a gateway to more subcommands!",
				Action: func() {
					fmt.Println("Test Value:", StrCmd.String("-test"))
				},
				Args: []string{
					"-test",
				},
			},
			"loop-test": {
				Description: "A loop that will continue until ctrl+c is pressed.",
				Action: func() {
					var Terminate bool
					c := make(chan os.Signal, 1)
					signal.Notify(c, os.Interrupt)

					go func() {
						<-c
						signal.Stop(c)
						Terminate = true
					}()

					for {
						if Terminate {
							break
						} else {
							fmt.Println("Hello")
							time.Sleep(3 * time.Second)
						}
					}
				},
			},
		},
	}
	a.Run(">>: ")
}

```
