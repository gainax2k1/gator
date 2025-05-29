package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gainax2k1/gator/internal/config"
)

func main() {

	userConfig, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config: ", err)
	}

	/* removed:
	userConfig.SetUser("Rico")
	userConfig, err = config.Read()
	if err != nil {
		fmt.Println("Error reading config: ", err)
	}
	*/
	gatorState := state{appState: &userConfig}

	gatorCommands := commands{cliCommands: make(map[string]func(*state, command) error)}

	gatorCommands.register("login", handlerLogin)
	gatorArgs := os.Args

	if len(gatorArgs) < 2 {
		log.Fatal("Required command name missing; exiting program")
	}

	commandName := gatorArgs[1]
	commandArgs := gatorArgs[2:]

	gatorCommand := command{name: commandName, arguments: commandArgs}
	err = gatorCommands.run(&gatorState, gatorCommand)
	if err != nil {
		log.Fatal("Fatal error: ", err)
	}

	/*
		You'll need to split the command-line arguments into the command name and the arguments slice to create a command instance.
		Use the commands.run method to run the given command and print any errors returned.

	*/

	// for debugging: fmt.Printf("%+v\n", userConfig)
	/*
		%v — prints the struct values, but not the field names (e.g., {postgres://example Rico})
		%+v — prints the struct values with field names (e.g., {DbURL:postgres://example CurrentUserName:Rico})
		%#v — prints the Go-syntax representation (e.g., config.Config{DbURL:"postgres://example", CurrentUserName:"Rico"})
	*/

}
