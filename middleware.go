package main

import (
	"context"
	"fmt"
	"vran/blogaggre/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		username := s.cfg.CurrentUsername

		user, err := s.db_queries.GetUser(context.Background(), username)
		if err != nil {
			return fmt.Errorf("Failed to get user \"%s\": %w", username, err)
		}

		return handler(s, cmd, user)
	}
}
