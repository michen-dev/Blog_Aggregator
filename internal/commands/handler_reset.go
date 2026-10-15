package commands

import (
	"fmt"
	"context"
)

func HandlerReset(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.Name)
	}

	err := s.Db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("Failed resetting: %v", err)
	}
	return nil
}