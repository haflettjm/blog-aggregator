package main

import (
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/database"
	"blog-aggregator/internal/rss"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(s *state, cmd command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return errors.New("command not found")
	}
	return handler(s, cmd)
}

func handlerLogin(state *state, cmd command) error {
	if len(cmd.args) < 1 || cmd.args[0] == "" {
		return errors.New("invalid arguments")
	}
	uname := cmd.args[0]

	_, err := state.db.GetUser(context.Background(), uname)

	if err != nil {
		return fmt.Errorf("user %s not found", uname)
	}

	state.cfg.SetUser(uname)
	fmt.Println(uname)
	return nil
}

func handlerRegister(state *state, cmd command) error {
	if len(cmd.args) < 1 || cmd.args[0] == "" {
		return errors.New("invalid arguments")
	}
	uname := cmd.args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      uname,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := state.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}
	state.cfg.SetUser(uname)
	fmt.Println("User has been registered as:", uname)
	return nil
}

func handlerReset(state *state, cmd command) error {
	if err := state.db.DeleteAllUsers(context.Background()); err != nil {
		return err
	}
	fmt.Println("All users have been deleted")
	return nil
}

func handlerListUsers(state *state, cmd command) error {
	users, err := state.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == state.cfg.CURR_UNAME {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Println(user.Name)
		}
	}
	return nil
}

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

func handlerAgg(state *state, cmd command) error {
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

func resetFeeds(ctx context.Context, state *state) error {
	if err := state.db.DeleteFeeds(ctx); err != nil {
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

// func handlerLogout(state *state) {
// 	// Implement logout logic here
// }

func main() {
	// fmt.Println("Started Succesffully!")

	c := &config.Config{}

	c.Read(".gatorconfig.json")
	db, err := sql.Open("postgres", c.DB_URL)
	if err != nil {
		fmt.Println("Error Connecting to the DB: \n", err)
		os.Exit(1)
	}
	defer db.Close()
	dbQueries := database.New(db)
	s := &state{
		cfg: c,
		db:  dbQueries,
	}

	cmds := &commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	//cmds.register("feeds", handlerListFeeds)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("resetfeed", handlerResetFeeds)
	cmds.register("feeds", handlerListFeeds)

	if len(os.Args) < 2 {
		fmt.Println("Usage: programn <command>[args....]")
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := cmds.run(s, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// c.SetUser("fivek77")

	// fmt.Println(c.CURR_UNAME)
}
