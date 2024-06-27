package main

import (
	"github.com/johannessarpola/lutakkols/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		// todo
		return
	}
}
