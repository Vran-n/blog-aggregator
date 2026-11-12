package main

import (
	"context"
	"database/sql"
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

func markAsFetched(queries *database.Queries, feed database.Feed) error {
	db_markfeedfetched := database.MarkFeedFetchedParams{
		ID:            feed.ID,
		UpdatedAt:     time.Now(),
		LastFetchedAt: sql.NullTime{Time: time.Now()},
	}

	if err := queries.MarkFeedFetched(context.Background(), db_markfeedfetched); err != nil {
		return fmt.Errorf("Cannot mark feed \"%s\" with url \"%s\" as fetched: %w", feed.ID, feed.Url, err)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		err := errors.New("Duration must not be empty!")
		return fmt.Errorf("%w\nUSAGE: agg [TIME](s or m or h)", err)
	}
	time_between_reqs := cmd.arguments[0]
	interval, err := time.ParseDuration(time_between_reqs)

	if err != nil {
		return fmt.Errorf("Could not parse string to time.duration: %w", err)
	}

	feeds, err := s.db_queries.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Cannot fetch next feeds: %w", err)
	}

	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		fmt.Println("Fetching all feeds:")

		for _, feed := range feeds {
			if err := markAsFetched(s.db_queries, feed); err != nil {
				return err
			}

			rss_feed, err := funcFeed(context.Background(), feed.Url)
			if err != nil {
				return err
			}

			fmt.Printf("\"%s\": %s\n", rss_feed.Channel.Title, rss_feed.Channel.Link)

			for _, item := range rss_feed.Channel.Item {
				truetime, err := time.Parse(time.RFC1123Z, item.PubDate)
				if err != nil {
					fmt.Printf("Unable to parse time of feed \"%s\": %s\n", rss_feed.Channel.Title, item.PubDate)
				}

				db_post := database.CreatePostParams{
					ID:          uuid.New(),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Title:       sql.NullString{String: item.Title},
					Url:         item.Link,
					Description: sql.NullString{String: item.Description},
					PublishedAt: sql.NullTime{Time: truetime},
				}

				s.db_queries.CreatePost(context.Background(), db_post)
			}

		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		err := errors.New("Name must not be empty!")
		return fmt.Errorf("%w\nUSAGE: addfeed [NAME] [URL]", err)
	} else if len(cmd.arguments) == 1 {
		err := errors.New("URL must not be empty!")
		return fmt.Errorf("%w\nUSAGE: addfeed [NAME] [URL]", err)
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

	fmt.Printf("\"%s\" created a new feed \"%s\" with url \"%s\"\n", user.Name, feed.Name, feed.Url)

	// this is only bad if i were to forget to change the name of the command
	// SO REMIND YOURSELF
	pseudocmd := command{
		name:      "follow",
		arguments: []string{url},
	}
	if err := handlerFollow(s, pseudocmd, user); err != nil {
		return fmt.Errorf("Failed to follow feed \"%s\": %w", url, err)
	}
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

		fmt.Printf("> \"%s\": %s | created by: \"%s\"\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}
