package main

import (
	"context"
	"errors"
	"fmt"
	"time"
	"vran/blogaggre/internal/database"

	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Site name must not be empty!")
	}
	url := cmd.arguments[0]

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
		return fmt.Errorf("Failed to create new feed follow for \"%s\" with url \"%s\": %w", user.Name, url, err)
	}

	fmt.Printf("\"%s\" is now following a \"%s\" feed with url \"%s\"\n", feedfollow.Username, feedfollow.Feedname, url)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedfollows, err := s.db_queries.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Failed to get all feed follows: %w", err)
	}

	for _, follows := range feedfollows {
		feed, err := s.db_queries.GetFeedFromID(context.Background(), follows.FeedID)
		if err != nil {
			return fmt.Errorf("Failed to get feed with ID \"%s\": %w", follows.FeedID, err)
		}

		fmt.Printf("* %s ==> %s | %s\n", follows.Username, follows.Feedname, feed.Url)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Site name must not be empty!")
	}
	url := cmd.arguments[0]

	feed, err := s.db_queries.GetFeedFromURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Failed to get feed from url \"%s\": %w", url, err)
	}

	db_removefollow := database.RemoveFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	if err := s.db_queries.RemoveFollow(context.Background(), db_removefollow); err != nil {
		return fmt.Errorf("Failed to remove \"%s\" from feed follows: %w", user.Name, err)
	}

	fmt.Printf("\"%s\" unfollowed feed \"%s\" with url \"%s\"\n", user.Name, feed.Name, feed.Url)
	return nil
}
