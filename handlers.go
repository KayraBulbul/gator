package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kayrabulbul/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Username is required")
	}

	name := strings.ToLower(cmd.arguments[0])

	if _, err := s.db.GetUser(context.Background(), name); err != nil {
		log.Fatal("User doesn't exist")
	}

	err := s.cfg.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Printf("User set to %s", cmd.arguments[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("User is required")
	}

	name := strings.ToLower(cmd.arguments[0])
	if _, err := s.db.GetUser(context.Background(), name); err == nil {
		log.Fatal("User already exists")
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	s.db.CreateUser(context.Background(), params)

	err := s.cfg.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been created", name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		log.Fatalf("Couldn't reset database %v", err)
	}
	fmt.Println("Database reset successful")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("Couldn't get all users %v", err)
	}

	for _, user := range users {
		if s.cfg.Current_user_name == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}
