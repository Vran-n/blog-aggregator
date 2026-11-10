package main

import (
	"log"
	"os"
	"vran/blogaggre/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	mycommands := commands{container: map[string]func(*state, command) error{}}
	mycommands.register("login", handlerLogin)

	c, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	current_state := state{
		cfg: &c,
	}

	if len(os.Args) < 2 {
		log.Fatal("Missing command!")
	}

	cmd_name := os.Args[1]
	cmd_args := os.Args[2:]
	mycmd := command{
		name:      cmd_name,
		arguments: cmd_args,
	}

	if err := mycommands.run(&current_state, mycmd); err != nil {
		log.Fatal(err)
	}
}
