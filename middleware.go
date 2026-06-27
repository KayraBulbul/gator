package main

import (
	"context"

	"github.com/kayrabulbul/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		innerUser, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
		if err != nil {
			return err
		}
		return handler(s, cmd, innerUser)
	}
}
