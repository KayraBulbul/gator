package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strconv"
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
		return &RSSFeed{}, errors.New("error making request")
	}

	req.Header.Set("User-Agent", "gator")

	// Get the response
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, errors.New("error getting response")
	}
	defer res.Body.Close() // Close body after function ends

	// Read body to unmarshal
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, errors.New("error reading response body")
	}

	// Unmarshal into xmlData
	var xmlData RSSFeed
	if err = xml.Unmarshal(b, &xmlData); err != nil {
		return &RSSFeed{}, errors.New("error unmarshaling data")
	}

	return &xmlData, nil
}

func scrapeFeeds(s *state) error {
	dbFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return errors.New("error getting next feed")
	}

	err = s.db.MarkFeedFetched(context.Background(), dbFeed.ID)
	if err != nil {
		return errors.New("error marking feed fetched")
	}

	feed, err := fetchFeed(context.Background(), dbFeed.Url)
	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Item {
		PublishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return errors.New("error parsing date")
		}

		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       html.UnescapeString(item.Title),
			Url:         item.Link,
			Description: html.UnescapeString(item.Description),
			PublishedAt: sql.NullTime{
				Time:  PublishedAt,
				Valid: true,
			},
			FeedID: dbFeed.ID,
		}

		err = s.db.CreatePost(context.Background(), params)
		if err != nil {
			if strings.Contains(err.Error(), "23505") {
				continue
			}
			return errors.New("error creating post")
		}

		fmt.Printf("Post for %s was created.\n", html.UnescapeString(item.Title))
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("no argument provided")
	}
	timer := cmd.arguments[0]

	timeBetweenReqs, err := time.ParseDuration(timer)
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %s\n", timer)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	// validate arguments
	if len(cmd.arguments) != 2 {
		log.Fatal("Not enough arguments")
	}
	name := cmd.arguments[0]
	url := cmd.arguments[1]

	FeedParams := database.CreateFeedParams{
		ID:            uuid.New(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Name:          name,
		Url:           url,
		UserID:        user.ID,
		LastFetchedAt: sql.NullTime{Valid: false},
	}

	feed, err := s.db.CreateFeed(context.Background(), FeedParams)
	if err != nil {
		return err
	}

	FeedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), FeedFollowParams)
	if err != nil {
		return err
	}

	fmt.Println(feed.ID)
	fmt.Println(feed.CreatedAt)
	fmt.Println(feed.UpdatedAt)
	fmt.Println(feed.Name)
	fmt.Println(feed.Url)
	fmt.Println(feed.UserID)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Println(feed.Feedname)
		fmt.Println(feed.Url)
		fmt.Println(feed.Username)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return errors.New("need URL argument")
	}

	url := cmd.arguments[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	FeedParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), FeedParams)
	if err != nil {
		return err
	}

	fmt.Println(feed.FeedName)
	fmt.Println(feed.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return err
	}

	for _, feed := range feedFollows {
		fmt.Println(feed.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return errors.New("need URL argument")
	}
	feedURL := cmd.arguments[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return err
	}

	params := database.DeleteFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.DeleteFollow(context.Background(), params)
	if err != nil {
		return err
	}
	return nil
}

func handlerBrowse(s *state, cmd command) error {
	var limit int32 = 2
	if len(cmd.arguments) == 1 {
		parsedLimit, err := strconv.ParseInt(cmd.arguments[0], 10, 32)
		if err != nil {
			return err
		}

		limit = int32(parsedLimit)
	}

	posts, err := s.db.GetPostForUser(context.Background(), limit)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println(post)
	}
	return nil
}
