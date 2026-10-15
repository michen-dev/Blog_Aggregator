package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"time"
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


func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Declare client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create & Set up Header
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "gator")
	req.Header.Set("Content-Type", "application/xml")

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed requesting: %v", err)
	}
	defer resp.Body.Close()

	// Decode into struct
	decoder := xml.NewDecoder(resp.Body)
	var rss_feed RSSFeed
	err = decoder.Decode(&rss_feed)
	if err != nil {
		return nil, fmt.Errorf("Error decoding XML response: %v", err)
	}

	rss_feed.Channel.Title = html.UnescapeString(rss_feed.Channel.Title)
	rss_feed.Channel.Description = html.UnescapeString(rss_feed.Channel.Description)
	for i, item := range rss_feed.Channel.Item {
		item.Description = html.UnescapeString(item.Description)
		item.Title = html.UnescapeString(item.Title)
		rss_feed.Channel.Item[i] = item
	}

	fmt.Printf("XML Data: %+v\n", rss_feed)
	return &rss_feed, nil
}