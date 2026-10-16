package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/michen-dev/Blog_Aggregator/internal/database"
)

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s <feed_name> <feed_url>", cmd.Name)
	}

	ctx := context.Background()

	_, err := s.Db.GetFeed(ctx, cmd.Args[1])
	if err == nil {
		return fmt.Errorf("This feed already existed")
	}

	feed_params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}
	feed, err := s.Db.CreateFeed(ctx, feed_params)
	if err != nil {
		return fmt.Errorf("Creating feed failed: %v", err)
	}

	feed_follow_params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feed_follow, err := s.Db.CreateFeedFollow(ctx, feed_follow_params)
	if err != nil {
		return fmt.Errorf("Creating feed_follow failed: %v", err)
	}

	fmt.Printf("Feed Created Successfully: %+v\n", feed)
	fmt.Printf("Feed_Following Created Successfully: %s - %s\n", feed_follow.UserName, feed_follow.FeedName)
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
		if i != len(feeds)-1 {
			fmt.Println("=====================================")
		}
	}

	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s \"<url>\"", cmd.Name)
	}

	ctx := context.Background()

	url := cmd.Args[0]
	feed, err := s.Db.GetFeed(ctx, url)
	if err != nil {
		return err
	}

	feed_follow_params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feed_follow, err := s.Db.CreateFeedFollow(ctx, feed_follow_params)
	if err != nil {
		return fmt.Errorf("Error creating feed_follow: %v", err)
	}
	fmt.Println("Successfully created feed_follow!")
	fmt.Printf("%s - %s\n", feed_follow.UserName, feed_follow.FeedName)
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.Name)
	}

	ctx := context.Background()

	following, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("Failed querying feed_follows for user: %v", err)
	}

	for _, feed_follow := range following {
		fmt.Println(feed_follow.FeedName)
	}
	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s \"<url>\"", cmd.Name)
	}

	ctx := context.Background()

	feed, err := s.Db.GetFeed(ctx, cmd.Args[0])
	if err != nil {
		return err
	}

	delete_params := database.DeleteFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}
	err = s.Db.DeleteFeedFollow(ctx, delete_params)
	if err != nil {
		return err
	}

	fmt.Println("Successfully Unfollowed!")
	return nil
}

// ====== Wrapper Function
func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		ctx := context.Background()

		user, err := s.Db.GetUser(ctx, s.Cfg.Current_user_name)
		if err != nil {
			return fmt.Errorf("Error querying user: %v", err)
		}

		err = handler(s, cmd, user)
		if err != nil {
			return err
		}
		return nil
	}
}