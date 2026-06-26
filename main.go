package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/kayrabulbul/gator/internal/config"
	"github.com/kayrabulbul/gator/internal/database"
	_ "github.com/lib/pq"
)

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

	// register commands
	err = appCommands.register("login", handlerLogin)
	if err != nil {
		log.Fatalf("Error registering login: %v", err)
	}
	err = appCommands.register("register", handlerRegister)
	if err != nil {
		log.Fatalf("Error registering register: %v", err)
	}
	err = appCommands.register("reset", handlerReset)
	if err != nil {
		log.Fatalf("Error registering reset: %v", err)
	}
	err = appCommands.register("users", handlerUsers)
	if err != nil {
		log.Fatalf("Error registering users: %v", err)
	}
	err = appCommands.register("agg", handlerAgg)
	if err != nil {
		log.Fatalf("Error registering agg: %v", err)
	}
	err = appCommands.register("addfeed", addFeed)
	if err != nil {
		log.Fatalf("Error registering addfeed: %v", err)
	}
	err = appCommands.register("feeds", handlerFeeds)
	if err != nil {
		log.Fatalf("Error registering feeds: %v", err)
	}
	err = appCommands.register("follow", handlerFollow)
	if err != nil {
		log.Fatalf("Error registering follow: %v", err)
	}
	err = appCommands.register("following", handlerFollowing)
	if err != nil {
		log.Fatalf("Error registering following: %v", err)
	}

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
