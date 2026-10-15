package commands

import (
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"

	"github.com/michen-dev/Blog_Aggregator/internal/database"
	"github.com/michen-dev/Blog_Aggregator/rss"
)


func HandlerRSS(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.Name)
	}

	_, err := rss.FetchFeed(context.Background(), feed_url)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s <feed_name> <feed_url>", cmd.Name)
	}

	ctx := context.Background()
	user, err := s.Db.GetUser(ctx, s.Cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	feed_params := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.Args[0],
		Url: cmd.Args[1],
		UserID: user.ID,
	}

	feed, err := s.Db.CreateFeed(ctx, feed_params)
	if err != nil {
		return fmt.Errorf("Creating feed failed: %v", err)
	}

	fmt.Printf("Feed Data: %+v\n", feed)
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.Name)
	}

	ctx := context.Background()
	feeds, err := s.Db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("Failed fetching feeds: %v", err) 
	}

	for i, feed := range feeds {
		fmt.Printf("name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		user, err := s.Db.GetFeedUser(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("Failed getting user: %v", err)
		}
		fmt.Printf("username: %s\n", user.Name)
		if i != len(feeds) - 1 {
			fmt.Println("=====================================")
		}
	}

	return nil
}