package main

import (
	"errors"
	"fmt"
	"vran/blogaggre/internal/config"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Please input a username!")
	}
	username := cmd.arguments[0]
	config.SetUser(s.cfg, username)

	fmt.Printf("Set new username: \"%s\"\n", username)
	return nil
}
