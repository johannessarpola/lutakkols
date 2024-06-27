package main

import (
	"github.com/johannessarpola/go-lutakko-gigs/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		// todo
		return
	}
}
