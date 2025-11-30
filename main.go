package main

import (
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/database"
	"context"
	"database/sql"
	"errors"
	"fmt"
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
