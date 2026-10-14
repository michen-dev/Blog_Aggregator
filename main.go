package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/michen-dev/Blog_Aggregator/internal/commands"
	"github.com/michen-dev/Blog_Aggregator/internal/config"
	"github.com/michen-dev/Blog_Aggregator/internal/database"
)

func main() {
	// Read current config JSON file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed reading: %v", err)
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", cfg.Db_url)
	// Create database queries
	dbQueries := *database.New(db)
	
	// Initialize state
	var state commands.State
	state.Cfg = &cfg
	state.Db = &dbQueries

	// Initialize cmds and register login handler
	var cmds commands.Commands
	cmds.Commands = make(map[string]func(*commands.State, commands.Command) error)
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users", commands.HandlerUsers)

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