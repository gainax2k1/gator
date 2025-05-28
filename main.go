package main

import (
	"fmt"

	"github.com/gainax2k1/gator/internal/config"
)

func main() {

	userConfig, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config: ", err)
	}

	userConfig.SetUser("Rico")
	userConfig, err = config.Read()
	if err != nil {
		fmt.Println("Error reading config: ", err)
	}

	fmt.Printf("%+v\n", userConfig)
	/*
		%v — prints the struct values, but not the field names (e.g., {postgres://example Rico})
		%+v — prints the struct values with field names (e.g., {DbURL:postgres://example CurrentUserName:Rico})
		%#v — prints the Go-syntax representation (e.g., config.Config{DbURL:"postgres://example", CurrentUserName:"Rico"})
	*/

}
