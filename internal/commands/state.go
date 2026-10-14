package commands

import (
	"github.com/michen-dev/Blog_Aggregator/internal/config"
	"github.com/michen-dev/Blog_Aggregator/internal/database"
)


type State struct {
	Cfg *config.Config
	Db *database.Queries
}