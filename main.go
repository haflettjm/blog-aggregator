package main

import (
	"blog-aggregator/internal/config"
	"errors"
	"fmt"
	"os"
)

type state struct {
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
	state.cfg.SetUser(uname)
	fmt.Printf("User has been set as:\n %v\n", uname)
	return nil
}

// func handlerLogout(state *state) {
// 	// Implement logout logic here
// }

func main() {
	// fmt.Println("Started Succesffully!")

	c := &config.Config{}

	c.Read(".gatorconfig.json")

	s := &state{
		cfg: c,
	}

	cmds := &commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Usage: programn <command>[args....]")
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	if err := cmds.run(s, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// c.SetUser("fivek77")

	// fmt.Println(c.CURR_UNAME)
}
