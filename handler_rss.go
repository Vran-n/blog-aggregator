package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
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
