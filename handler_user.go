package main

import (
	"blog-aggregator/internal/database"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

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

func handlerListUsers(state *state, cmd command, user database.User) error {
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
