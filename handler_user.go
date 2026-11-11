package main

import (
	"context"
	"errors"
	"fmt"
	"time"
	"vran/blogaggre/internal/config"
	"vran/blogaggre/internal/database"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Username must not be empty!")
	}
	username := cmd.arguments[0]

	_, err := s.db_queries.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("Failed to get user \"%s\": %w", username, err)
	}

	if err := config.SetUser(s.cfg, username); err != nil {
		return fmt.Errorf("Failed to insert new user in config file: %w", err)
	}

	fmt.Printf("Logged in user: \"%s\"\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Username must not be empty!")
	}
	username := cmd.arguments[0]

	db_userparams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	user, err := s.db_queries.CreateUser(context.Background(), db_userparams)
	if err != nil {
		return fmt.Errorf("Failed to create new user: %w", err)
	}

	if err := config.SetUser(s.cfg, username); err != nil {
		return fmt.Errorf("Failed to insert new user in config file: %w", err)
	}

	fmt.Printf("Registered \"%s\"\n", username)
	fmt.Printf("user: %v\n", user)
	return nil
}
