package main

import (
	"blog-aggregator/internal/database"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func follow(state *state, url string) error {
	ctx := context.Background()
	feeds, err := state.db.ListFeedsByURL(ctx, url)

	if err != nil {
		return err
	}
	for _, feed := range feeds {
		user, err := state.db.GetUser(ctx, state.cfg.CURR_UNAME)
		if err != nil {
			return err
		}
		if user.Name == "" {
			return fmt.Errorf("user not found")
		}
		params := database.CreateFeedFollowParams{
			ID:        uuid.New(),
			UserID:    user.ID,
			FeedID:    feed.ID,
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		}
		inserted, err := state.db.CreateFeedFollow(ctx, params)
		if err != nil {
			return err
		}
		fmt.Println(inserted)
	}

	return nil
}

func following(state *state) error {
	ctx := context.Background()
	user, err := state.db.GetUser(ctx, state.cfg.CURR_UNAME)
	if err != nil {
		return err
	}

	follows, err := state.db.ListFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	for _, feed := range follows {
		fmt.Printf("%s \n", feed.FeedName)
	}
	return nil
}

func handlerFollowing(state *state, cmd command, user database.User) error {
	if state.cfg.CURR_UNAME == "" {
		return fmt.Errorf("user not logged in")
	}
	err := following(state)
	if err != nil {
		return err
	}
	return nil
}

func handlerFollow(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: follow <url>")
	}
	return follow(state, cmd.args[0])
}
