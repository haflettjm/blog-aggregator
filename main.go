package main

import (
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/database"
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func follow(state *state, url string) error {
	feeds, err := state.db.ListFeedsByURL(context.Background(), url)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	c := &config.Config{}
	c.Read(".gatorconfig.json")

	db, err := sql.Open("postgres", c.DB_URL)
	if err != nil {
		fmt.Println("Error Connecting to the DB:", err)
		os.Exit(1)
	}
	defer db.Close()

	s := &state{cfg: c, db: database.New(db)}

	cmds := &commands{handlers: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("resetfeed", handlerResetFeeds)
	cmds.register("feeds", handlerListFeeds)

	if len(os.Args) < 2 {
		fmt.Println("Usage: program <command> [args...]")
		os.Exit(1)
	}

	cmd := command{name: os.Args[1], args: os.Args[2:]}
	if err := cmds.run(s, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
