package main

import (
	"errors"
	"fmt"

	"github.com/realdnchka/gator-go/internal/config"
)

type command struct {
	Name string
	Args []string
}

type state struct {
	Config *config.Config
}

type commands struct {
	handler map[string]func(*state, command) error
}

func StateInit() state {
	return state{}
}

func CommandsInit() {
	cmds.register("login", handlerLogin)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("the login handler expects a single argument, the username\n")
	}

	s.Config.SetUser(cmd.Args[0])
	fmt.Printf("Succesfuly loged in with username: %s\n", cmd.Args[0])
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
