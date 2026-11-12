package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
	"vran/blogaggre/internal/database"

	"github.com/google/uuid"
)

type RSS_Feed struct {
	Channel struct {
		Title       string     `xml:"title"`
		Link        string     `xml:"link"`
		Description string     `xml:"description"`
		Item        []RSS_Item `xml:"item"`
	} `xml:"channel"`
}

func (rfeed *RSS_Feed) decodeEscapedHTML() {
	rfeed.Channel.Title = html.UnescapeString(rfeed.Channel.Title)
	rfeed.Channel.Description = html.UnescapeString(rfeed.Channel.Description)

	for i := range rfeed.Channel.Item {
		rfeed.Channel.Item[i].decodeEscapedHTML()
	}
}

type RSS_Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (ritem *RSS_Item) decodeEscapedHTML() {
	ritem.Title = html.UnescapeString(ritem.Title)
	ritem.Description = html.UnescapeString(ritem.Description)
}

func funcFeed(ctx context.Context, feedURL string) (*RSS_Feed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSS_Feed{}, fmt.Errorf("HTTP Unable to get new request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return &RSS_Feed{}, fmt.Errorf("HTTP Unable to get new response: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSS_Feed{}, err
	}

	var rss_feed RSS_Feed
	if err := xml.Unmarshal(data, &rss_feed); err != nil {
		return &RSS_Feed{}, fmt.Errorf("Unable to unmarshal bytes: %w", err)
	}

	rss_feed.decodeEscapedHTML()

	return &rss_feed, nil
}

func handlerAgg(s *state, cmd command) error {
	_, err := funcFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Site name must not be empty!")
	} else if len(cmd.arguments) == 1 {
		return errors.New("URL must not be empty!")
	}
	name := cmd.arguments[0]
	url := cmd.arguments[1]

	db_feedparams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	// cant ommit the first variable else it bitches
	feed, err := s.db_queries.CreateFeed(context.Background(), db_feedparams)
	if err != nil {
		return fmt.Errorf("Failed to get feed \"%s\": %w", url, err)
	}

	// this is only bad if i were to forget to change the name of the command
	// SO REMIND YOURSELF
	pseudocmd := command{
		name:      "follow",
		arguments: []string{url},
	}

	if err := handlerFollow(s, pseudocmd, user); err != nil {
		return fmt.Errorf("Failed to follow feed \"%s\": %w", url, err)
	}

	fmt.Printf("\"%s\" added a new feed \"%s\" with url \"%s\"\n", user.Name, feed.Name, feed.Url)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db_queries.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get all feeds: %w", err)
	}

	for _, feed := range feeds {
		user, err := s.db_queries.GetUserFromID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Failed to get user from id \"%s\": %w", feed.UserID, err)
		}

		fmt.Printf("> \"%s\" | %s | created by: \"%s\"\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}
