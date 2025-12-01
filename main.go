package main

import (
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/database"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

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
	cmds.register("reset", middlwareLoggedIn(handlerReset))
	cmds.register("users", middlwareLoggedIn(handlerListUsers))
	cmds.register("agg", middlwareLoggedIn(handlerAgg))
	cmds.register("addfeed", middlwareLoggedIn(handlerAddFeed))
	cmds.register("resetfeed", middlwareLoggedIn(handlerResetFeeds))
	cmds.register("feeds", middlwareLoggedIn(handlerListFeeds))
	cmds.register("follow", middlwareLoggedIn(handlerFollow))
	cmds.register("following", middlwareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlwareLoggedIn(handlerUnfollow))

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
