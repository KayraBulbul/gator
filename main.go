package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/kayrabulbul/gator/internal/config"
)

type state struct {
	pntr *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.commandMap[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.commandMap[name] = f
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Username is required")
	}

	err := s.pntr.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Printf("User set to %s", cmd.arguments[0])
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	appState := &state{&cfg}
	appCommands := commands{make(map[string]func(*state, command) error)}
	arguments := os.Args

	if len(arguments) < 2 {
		log.Fatal("Not enough arguments")
	}

	appCommand := command{
		name:      arguments[1],
		arguments: arguments[2:],
	}

	appCommands.register("login", handlerLogin)
	err = appCommands.run(appState, appCommand)
	if err != nil {
		log.Fatalf("Error running command: %v", err)
	}
}
