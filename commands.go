package main

import "fmt"

type command struct {
	name      string
	arguments []string
}

type commands struct {
	container map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	command, ok := c.container[cmd.name]
	if !ok {
		return fmt.Errorf("Unknown command \"%s\"", cmd.name)
	}

	if err := command(s, cmd); err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.container[name] = f
}
