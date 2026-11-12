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

	user, err := s.db_queries.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("Failed to get user \"%s\": %w", username, err)
	}

	if err := config.SetUser(s.cfg, user.Name); err != nil {
		return fmt.Errorf("Failed to insert new user in config file: %w", err)
	}

	fmt.Printf("Welcome \"%s\"!\n", user.Name)
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

	if err := config.SetUser(s.cfg, user.Name); err != nil {
		return fmt.Errorf("Failed to insert new user in config file: %w", err)
	}

	fmt.Printf("\"%s\" you are now registered!\n", user.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db_queries.RemoveUsers(context.Background()); err != nil {
		return fmt.Errorf("Failed to remove all users: %w", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db_queries.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get all users: %w", err)
	}

	for _, user := range users {
		if s.cfg.CurrentUsername == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}
