package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/kayrabulbul/gator/internal/config"
	"github.com/kayrabulbul/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	dbURL := cfg.Connection_string
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	dbQueries := database.New(db)

	appState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}
	appCommands := commands{make(map[string]func(*state, command) error)}
	appCommands.register("login", handlerLogin)
	appCommands.register("register", handlerRegister)
	appCommands.register("reset", handlerReset)
	appCommands.register("users", handlerUsers)

	arguments := os.Args

	if len(arguments) < 2 {
		log.Fatal("Not enough arguments")
	}

	appCommand := command{
		name:      arguments[1],
		arguments: arguments[2:],
	}

	err = appCommands.run(appState, appCommand)
	if err != nil {
		log.Fatalf("Error running command: %v", err)
	}
}
