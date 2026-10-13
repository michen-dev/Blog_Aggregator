package commands

import (
	"fmt"
)


func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	err := s.Cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("Failed setting username: %v", err)
	}

	fmt.Printf("Successfully set username: %s!\n", name)
	return nil
}