package commands

import (
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
	"log"

	"github.com/michen-dev/Blog_Aggregator/internal/database"
)


func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]
	cnt := context.Background()
	_, err := s.Db.GetUser(cnt, name)
	if err != nil {
		return fmt.Errorf("User does not exist")
	}

	err = s.Cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("Failed setting current user: %v", err)
	}

	fmt.Printf("Successfully Login: %s!\n", name)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]
	cnt := context.Background()
	_, err := s.Db.GetUser(cnt, name)
	if err == nil {
		return fmt.Errorf("Username already existed")
	}

	user_params := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
	}

	user, err := s.Db.CreateUser(cnt, user_params)
	if err != nil {
		return fmt.Errorf("Error creating new user: %v", err)
	}

	fmt.Println("Creating new user successfully!")
	log.Print(user)

	err = s.Cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("Failed setting current user: %v", err)
	}

	fmt.Printf("Successfully Login: %s!\n", name)

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.Name)
	}

	err := s.Db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("Failed resetting")
	}
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.Name)
	}

	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed fetching users data")
	}
	
	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if user.Name == s.Cfg.Current_user_name {
			fmt.Print(" (current)")
		}
		fmt.Print("\n")
	}

	return nil
}