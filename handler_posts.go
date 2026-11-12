package main

import (
	"context"
	"fmt"
	"strconv"
)

func handlerBrowse(s *state, cmd command) error {
	var limit string
	if len(cmd.arguments) == 0 {
		limit = "2"
	} else {
		limit = cmd.arguments[0]
	}

	intlimit, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		return fmt.Errorf("Unable to convert limit argument to int32: %w", err)
	}

	posts, err := s.db_queries.GetPosts(context.Background(), int32(intlimit))
	if err != nil {
		return fmt.Errorf("Cannot get posts: %w", err)
	}

	i := 1
	for _, post := range posts {
		fmt.Printf("POST #%d | Published at: %v\n", i, post.PublishedAt.Time)
		fmt.Printf("Title: %s\n", post.Title.String)
		fmt.Println("--------------------------------------------------")
		fmt.Printf("%s\n", post.Description.String)
		fmt.Println("==================================================")
		i++
	}
	return nil
}
