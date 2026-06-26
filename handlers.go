package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kayrabulbul/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("username is required")
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
		return errors.New("user is required")
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

	_, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(cmd.arguments[0])
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Make request variable
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, &strings.Reader{})
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")

	// Get the response
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close() // Close body after function ends

	// Read body to unmarshal
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	// Unmarshal into xmlData
	var xmlData RSSFeed
	if err = xml.Unmarshal(b, &xmlData); err != nil {
		return &RSSFeed{}, err
	}

	return &xmlData, nil
}

func handleAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(html.UnescapeString(feed.Channel.Title))
	fmt.Println(feed.Channel.Link)
	fmt.Println(html.UnescapeString(feed.Channel.Description))

	for _, item := range feed.Channel.Item {
		fmt.Println(html.UnescapeString(item.Title))
		fmt.Println(item.Link)
		fmt.Println(html.UnescapeString(item.Description))
		fmt.Println(item.PubDate)
	}
	return nil
}
