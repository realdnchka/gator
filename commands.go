package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/realdnchka/gator-go/internal/config"
	"github.com/realdnchka/gator-go/internal/database"
)

type command struct {
	Name string
	Args []string
}

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type commands struct {
	handler map[string]func(*state, command) error
}

func StateInit() state {
	return state{}
}

func CommandsInit() {
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("the login handler expects a single argument, the username")
	}

	if err := s.cfg.SetUser(cmd.Args[0]); err != nil {
		return err
	}
	fmt.Printf("Succesfuly loged in with username: %s\n", cmd.Args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("the login handler expects a single argument, the username")
	}

	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}
	u := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}
	_, err = s.db.CreateUser(context.Background(), u)
	if err != nil {
		return err
	}
	s.cfg.UserName = u.Name
	fmt.Printf("created user: %s", u.Name)
	log.Printf("username: %s, created_at: %v, updated_at: %v, id: %v", u.Name, u.CreatedAt, u.UpdatedAt, u.ID)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	if err := c.handler[cmd.Name](s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(s *state, cmd command) error) {
	c.handler[name] = f
}
