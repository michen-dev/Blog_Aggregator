package commands

import (
	"context"
	"fmt"

	"github.com/michen-dev/Blog_Aggregator/rss"
)

func HandlerRSS(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.Name)
	}

	_, err := rss.FetchFeed(context.Background(), feed_url)
	if err != nil {
		return err
	}
	return nil
}
