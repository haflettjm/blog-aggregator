package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"blog-aggregator/internal/database"
	"blog-aggregator/internal/rss"
)

func fetchFeed(ctx context.Context, feedURL string) (*rss.RSSFeed, error) {
	feed := &rss.RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = feed.Parse(data)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

func addFeed(ctx context.Context, state *state, url string, name string) error {
	user := state.cfg.CURR_UNAME

	dbUser, err := state.db.GetUser(ctx, user)
	if err != nil {
		return err
	}

	feed, err := fetchFeed(ctx, url)
	if err != nil {
		return err
	}
	feedParams := database.CreateFeedParams{
		UserID:      dbUser.ID,
		Url:         url,
		Name:        name,
		Description: feed.Channel.Description,
	}
	_, err = state.db.CreateFeed(ctx, feedParams)

	if err != nil {
		return err
	}

	// fmt.Printf("Feed added: %s (%s)\n", url, user)

	fmt.Println(state.db.GetFeedsByUserId(ctx, dbUser.ID))
	return nil
}

func handlerAddFeed(state *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("Usage: addfeed <url> <name>")
	}

	if state.cfg.CURR_UNAME == "" {
		return errors.New("User not logged in")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	if err := addFeed(context.Background(), state, url, name); err != nil {
		return err
	}
	return nil
}

func handlerResetFeeds(state *state, cmd command) error {
	if err := resetFeeds(context.Background(), state); err != nil {
		return err
	}
	return nil
}

func handlerListFeeds(state *state, cmd command) error {
	if err := listFeeds(context.Background(), state); err != nil {
		return err
	}
	return nil
}

func listFeeds(ctx context.Context, state *state) error {
	users, err := state.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	for _, user := range users {
		feeds, err := state.db.GetFeedsByUserId(ctx, user.ID)
		if err != nil {
			fmt.Println(err)
		}
		for _, feed := range feeds {
			fmt.Println("Name: ", feed.Name)
			fmt.Println("URL: ", feed.Url)
			fmt.Println("Added By: ", user.Name)
		}
	}
	return nil
}
