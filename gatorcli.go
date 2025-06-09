package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gainax2k1/gator/internal/config"
	"github.com/gainax2k1/gator/internal/database"
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

func users(s *state, cmd command) error {
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

}

// Create a register handler and register it with the commands. Usage:
func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("register handler expects a single argument (the username); none received")
	}

	var newUser database.CreateUserParams

	newUser.Name = cmd.arguments[0]
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
		return err
		os.Exit(1) //maybe overkill? maybe just return err?
	}
	fmt.Println("successfully reset database.")
	os.Exit(0)
	return nil //not sure why neccessary for compiler
}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"

	// Create a context with 10-second timeout, instead of just background
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rssfeed, err := fetchFeed(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	fmt.Printf("%+v\n", rssfeed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {

	if len(cmd.arguments) < 2 {
		return fmt.Errorf("add feed handler expects two arguments (<name of feed> <url>); missing one or both ")
	}

	current_user := s.appState.CurrentUserName
	user_uuid, err := s.db.GetUserIDBName(context.Background(), current_user)
	if err != nil {
		return fmt.Errorf("error retrieving current users id: %w", err)
	}
	var newFeed database.CreateFeedParams

	newFeed.Name = cmd.arguments[0]
	newFeed.Url = cmd.arguments[1]
	newFeed.CreatedAt = time.Now()
	newFeed.UpdatedAt = time.Now()
	newFeed.UserID = user_uuid

	feed, err := s.db.CreateFeed(context.Background(), newFeed)
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	fmt.Printf("New feed created.\n Feed name: %s\nurl: %s\n", newFeed.Name, newFeed.Url) // nicely formated
	//fmt.Printf("%+v\n", newFeed)                                                          // ugly, brute force the whole dang thing

	// CH4:L1 automatically create a feed follow record for the current user when they add a feed.

	var feed_follow_params database.CreateFeedFollowParams
	feed_follow_params.FeedID = feed.ID
	feed_follow_params.UserID = user_uuid

	_, err = s.db.CreateFeedFollow(context.Background(), feed_follow_params)
	if err != nil {
		return fmt.Errorf("error creating feedfollow record: %w", err)
	}

	return nil
	/*
	   Add a new command called addfeed. It takes two args:
	   name: The name of the feed
	   url: The URL of the feed
	   At the top of the handler, get the current user from the database and connect the feed to that user.

	   If everything goes well, print out the fields of the new feed record.
	*/
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		return fmt.Errorf("error: no feeds found")
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error looking up user id: %w", err)
		}

		fmt.Printf("Feed name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("By User: %s\n", user.Name)
	}

	return nil
	/*
		Add a new feeds handler. It takes no arguments and prints all the feeds in the database to the console. Be sure to include:

		The name of the feed
		The URL of the feed
		The name of the user that created the feed (you might need a new SQL query)
	*/
}

func handerFollow(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("follow handler expects one argument (<url>); missing url ")
	}
	url := cmd.arguments[0]
	current_user := s.appState.CurrentUserName

	feed_uuid, err := s.db.GetFeedIDbyURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error looking up feed id: %w", err)
	}

	user_uuid, err := s.db.GetUserIDBName(context.Background(), current_user)
	if err != nil {
		return fmt.Errorf("error looking up current user id: %w", err)
	}

	var feed_follow_params database.CreateFeedFollowParams
	feed_follow_params.FeedID = feed_uuid
	feed_follow_params.UserID = user_uuid

	feedFollowRecord, err := s.db.CreateFeedFollow(context.Background(), feed_follow_params)
	if err != nil {
		return err
	}
	fmt.Printf("Feed name: %s\n", feedFollowRecord.FeedName)
	fmt.Printf("Current user: %s\n", current_user)

	/*
		Add a follow command. It takes a single url argument and creates a new feed follow record for the current user.
		It should print the name of the feed and the current user once the record is created (which the query we just made should support).
		You'll need a query to look up feeds by URL.*/

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.arguments) > 0 {
		return fmt.Errorf("following handler expects no arguments")
	}
	username := s.appState.CurrentUserName
	user_uuid, err := s.db.GetUserIDBName(context.Background(), username)
	if err != nil {
		return fmt.Errorf("error getting user id: %w", err)
	}
	feed_follows_list, err := s.db.GetFeedFollowsForUser(context.Background(), user_uuid)
	if err != nil {
		return fmt.Errorf("error retrieving feed follows: %w", err)
	}
	for _, feed := range feed_follows_list {
		feed_name, err := s.db.GetFeedNameByUUID(context.Background(), feed.FeedID)
		if err != nil {
			return fmt.Errorf("error looking up feed name by id: %W", err)
		}
		fmt.Printf("Feed name: %s\n", feed_name)
	}
	return nil
}
