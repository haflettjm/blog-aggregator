package main

import (
	"context"
	"fmt"
)

func handlerReset(state *state, cmd command) error {
	if err := state.db.DeleteAllUsers(context.Background()); err != nil {
		return err
	}
	fmt.Println("All users have been deleted")
	return nil
}

func resetFeeds(ctx context.Context, state *state) error {
	if err := state.db.DeleteFeeds(ctx); err != nil {
		return err
	}
	return nil
}
