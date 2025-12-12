package main

import (
	"net/http"
	"encoding/xml"
	"fmt"
	"io"
	"context"
	"html"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error connecting to URL: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil || res.StatusCode > 299 {
		return &RSSFeed{}, fmt.Errorf("Error executing request: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error reading response body: %w", err)
	}

	defer res.Body.Close()

	var newFeed RSSFeed
	err = xml.Unmarshal(body, &newFeed)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error extracting information from response body: %w", err)
	}

	unescape(&newFeed)

	return &newFeed, nil
}

func unescape(r *RSSFeed) {
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)

	for _, item := range r.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
}
