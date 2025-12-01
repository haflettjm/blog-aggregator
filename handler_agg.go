package main

import (
	"blog-aggregator/internal/database"
	"context"
	"fmt"
)

func handlerAgg(state *state, cmd command, user database.User) error {
	feeds := []string{}
	if len(cmd.args) == 0 {
		feeds = []string{
			"https://hnrss.org/newest",
			"https://www.wagslane.dev/index.xml",
		}
	} else {
		feeds = cmd.args
	}

	for _, feedURL := range feeds {
		feed, err := fetchFeed(context.Background(), feedURL)
		if err != nil {
			return err
		}
		// fmt.Printf("Feed: %s\n", feed.Channel.Title)
		// for _, item := range feed.Channel.Items {
		// 	fmt.Printf("- %s\n", item.Title)
		// }
		fmt.Println(feed)
	}

	return nil
}
