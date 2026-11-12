package main

import (
	"context"
	"errors"
	"fmt"
	"time"
	"vran/blogaggre/internal/database"

	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Site name must not be empty!")
	}
	username := s.cfg.CurrentUsername
	url := cmd.arguments[0]

	user, err := s.db_queries.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("Failed to get user \"%s\": %w", username, err)
	}

	feed, err := s.db_queries.GetFeedFromURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Failed to get feed from url \"%s\": %w", url, err)
	}

	db_feedfollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedfollow, err := s.db_queries.CreateFeedFollow(context.Background(), db_feedfollow)
	if err != nil {
		return fmt.Errorf("Failed to create new feed follow for \"%s\" with url \"%s\": %w", username, url, err)
	}

	fmt.Printf("Added feed follow for \"%s\" with feed name \"%s\"\n", feedfollow.Username, feedfollow.Feedname)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	username := s.cfg.CurrentUsername

	user, err := s.db_queries.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("Failed to get user \"%s\": %w", username, err)
	}

	feedfollows, err := s.db_queries.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Failed to get all feed follows: %w", err)
	}

	for _, follows := range feedfollows {
		fmt.Printf("* %s ==> %s | %s\n", follows.Username, follows.Feedname, follows.Feedurl)
	}
	return nil
}
