package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gainax2k1/gator/internal/config"
	"github.com/gainax2k1/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type state struct {
	db       *database.Queries
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

func users(s *state, cmd command) error { /// THIS IS PLACEHOLDER CODE< NEEDS WORK!!!
	usernames, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(usernames) == 0 {
		return fmt.Errorf("error: no usersnames found")
	}

	for _, username := range usernames {
		fmt.Printf(" * %s", username)
		if username == s.appState.CurrentUserName {
			fmt.Print(" (current)")
		}
		fmt.Printf("\n")
	}
	return nil

	/* i don'tknow if any of this is right for this
	if gatorCommand, exists :- c.cliCommands[cmd.name]; exists{
		return gatorCommand(s, cmd)
	} else {
		return fmt.Errorf(("command not found"))
	}
	*/
}

// Create a register handler and register it with the commands. Usage:
func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("register handler expects a single argument (the username); none received")
	}

	var newUser database.CreateUserParams

	newUser.Name = cmd.arguments[0]
	newUser.ID = uuid.New()
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	_, err := s.db.CreateUser(context.Background(), newUser)
	if err != nil {
		value, ok := err.(*pq.Error)
		if ok {
			if value.Code == pq.ErrorCode("23505") {
				fmt.Println("user already exists, exiting...")
				os.Exit(1)

			}
		}
		return fmt.Errorf("error creating user: %w", err)
	}
	s.appState.CurrentUserName = newUser.Name

	fmt.Printf("\nUser '%s' has been registered\n", newUser.Name)

	err = s.appState.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("\nUser has been set to %s\n", cmd.arguments[0])

	//error checking logging
	fmt.Println(newUser)

	return nil

}

func handlerLogin(s *state, cmd command) error {

	/*
		- If the command's arg's slice is empty, return an error; the login handler expects a single argument, the username.
		- Use the state's access to the config struct to set the user to the given username. Remember to return any errors.
		- Print a message to the terminal that the user has been set.*/

	if len(cmd.arguments) == 0 {
		return fmt.Errorf("login handler expects a single argument (the username); none received")
	}

	// Update the login command handler to error (and exit with code 1) if the given username doesn't exist in the database.
	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		fmt.Printf("username doesn't exist in database: %s", cmd.arguments[0])
		os.Exit(1)
	}

	err = s.appState.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("\nUser has been set to %s\n", cmd.arguments[0])
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		fmt.Println("error reseting databse: ", err)
		os.Exit(1) //maybe overkill? maybe just return err?
	}
	fmt.Println("successfully reset database.")
	os.Exit(0)
	return nil //not sure why neccessary
}
