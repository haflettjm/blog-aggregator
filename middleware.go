package main

import (
	"blog-aggregator/internal/database"
	"context"
	"fmt"
)

func middlwareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		if s.cfg.CURR_UNAME == "" {
			return fmt.Errorf("User Not logged In.")
		}
		user, err := s.db.GetUser(context.Background(), s.cfg.CURR_UNAME)
		if err != nil {
			return fmt.Errorf("couldn't get user: %w", err)
		}
		return handler(s, cmd, user)
	}
}
