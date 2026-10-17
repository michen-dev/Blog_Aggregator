package commands

import (
	"context"
	"encoding/xml"
	"database/sql"
	"fmt"
	"html"
	"net/http"
	"time"
	"strconv"

	"github.com/google/uuid"
	"github.com/michen-dev/Blog_Aggregator/internal/database"
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

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <time_between_reqs>", cmd.Name)
	}

	// Getting time between reqs arg
	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error parsing time value : %v", err)
	}

	// Set up ticker and run
	ticker := time.NewTicker(time_between_reqs)
	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			fmt.Printf("Error scraping feed: %v\n", err)
		}
		fmt.Println("=====================================")
	}
}

func HandlerBrowse(s *State, cmd Command, user database.User) error {
	limit := 2
	if len(cmd.Args) > 1 {
		return fmt.Errorf("Usage: %s <limit>", cmd.Name)
	} else if (len(cmd.Args) == 1) {
		var err error
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("Error converting limit to int: %v", err)
		}
	}

	posts, err := s.Db.GetPostsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error querying posts: %v", err)
	}

	for i := range limit {
		fmt.Println(posts[i].Title)
	}
	return nil
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
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

	// Formatting string values
	rss_feed.Channel.Title = html.UnescapeString(rss_feed.Channel.Title)
	rss_feed.Channel.Description = html.UnescapeString(rss_feed.Channel.Description)
	for i, item := range rss_feed.Channel.Item {
		item.Description = html.UnescapeString(item.Description)
		item.Title = html.UnescapeString(item.Title)
		rss_feed.Channel.Item[i] = item
	}

	// Print and return
	// fmt.Printf("XML Data: %+v\n", rss_feed)
	return &rss_feed, nil
}


func scrapeFeeds(s *State) error {
	ctx := context.Background()

	// Get next feed to fetch
	next_feed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("Error querying next feed to fetch: %v", err)
	}

	// Mark feed as fetched
	_, err = s.Db.MarkFeedFetched(ctx, next_feed.ID)
	if err != nil {
		return fmt.Errorf("Mark feed fetched failed: %v", err)
	}

	// Fetch feed
	fmt.Printf("Fetching feed: %s\n", next_feed.Name)
	rss_feed, err := fetchFeed(ctx, next_feed.Url)
	if err != nil {
		return fmt.Errorf("Error fetching feed: %v", err)
	}

	// Print fetch feed
	fmt.Printf("Fetched feed: %s\n", rss_feed.Channel.Title)
	// Save posts
	for _, item := range rss_feed.Channel.Item {
		// Parse pubDate and description
		parsedPubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Printf("Error parsing published date: %v\n", err)
		}
		parsedDescription := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		// Declare params & Create post
		params := database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: item.Title,
			Url: item.Link,
			Description: parsedDescription,
			PublishedAt: parsedPubDate,
			FeedID: next_feed.ID,
		}
		post, err := s.Db.CreatePost(ctx, params)
		if err != nil {
			return fmt.Errorf("Error saving post: %v", err)
		}
		fmt.Println(post.Title)
	}
	return nil
}