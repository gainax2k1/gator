package main

import (
	"fmt"

	"github.com/gainax2k1/gator/internal/config"
)

type state struct {
	appState *config.Config
}

type command struct {
	name      string
	arguments []string
}

/*
Create a commands struct. This will hold all the commands the CLI can handle.
 Add a map[string]func(*state, command) error field to it. This will be a map of command names to their handler functions.
*/

type commands struct {
	cliCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error { // This method runs a given command with the provided state if it exists.

	if gatorCommand, exists := c.cliCommands[cmd.name]; exists {
		return gatorCommand(s, cmd)
	} else {
		return fmt.Errorf("command not found")
	}

}

func (c *commands) register(name string, f func(*state, command) error) { // This method registers a new handler function for a command name.
	c.cliCommands[name] = f
}

func handlerLogin(s *state, cmd command) error {

	/*
		- If the command's arg's slice is empty, return an error; the login handler expects a single argument, the username.
		- Use the state's access to the config struct to set the user to the given username. Remember to return any errors.
		- Print a message to the terminal that the user has been set.*/

	if len(cmd.arguments) == 0 {
		return fmt.Errorf("login handler expects a single argument (the username); none received")
	}

	err := s.appState.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("\nUser has been set to %s\n", cmd.arguments[0])
	return nil
}
