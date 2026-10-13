package main

import (
	"log"
	"os"

	"github.com/michen-dev/Blog_Aggregator/internal/commands"
	"github.com/michen-dev/Blog_Aggregator/internal/config"
)

func main() {
	// Read current config JSON file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed reading: %v", err)
	}
	
	// Initialize state
	var state commands.State
	state.Cfg = &cfg

	// Initialize cmds and register login handler
	var cmds commands.Commands
	cmds.Commands = make(map[string]func(*commands.State, commands.Command) error)
	cmds.Register("login", commands.HandlerLogin)

	// Receive CLI args & Build input cmd
	if len(os.Args) < 2 {
		log.Fatal("Please provide command name")
	}
	cmd_name := os.Args[1]
	cmd_args := os.Args[2:]
	cmd := commands.Command{
		Name: cmd_name,
		Args: cmd_args,
	}

	// Run command
	err = cmds.Run(&state, cmd)
	if err != nil {
		log.Fatal(err)
	}
}