package main

import (
	"fmt"
	"log"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	username := user.Username

	fmt.Printf("Username: %s\n", username)

	App()
}
