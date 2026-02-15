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
	"github.com/realdnchka/gator-go/internal/rss"
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("the login handler expects a single argument, the username")
	}

	if err := s.cfg.SetUser(cmd.Args[0]); err != nil {
		return err
	}

	_, err := s.db.GetUserByName(context.Background(), cmd.Args[0])
	if err != nil {
		return errors.New("the user does not exists")
	}

	log.Printf("succesfuly loged in with username: %s\n", cmd.Args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("the login handler expects a single argument, the username")
	}

	_, err := s.db.GetUserByName(context.Background(), cmd.Args[0])
	if err == nil {
		return errors.New("the user exists")
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
	if err := s.cfg.SetUser(cmd.Args[0]); err != nil {
		return err
	}
	log.Printf("created user: %s", u.Name)
	log.Printf("username: %s, created_at: %v, updated_at: %v, id: %v", u.Name, u.CreatedAt, u.UpdatedAt, u.ID)
	return nil
}

func resetHandler(s *state, cmd command) error {
	s.db.ResetUsers(context.Background())
	log.Printf("user table succesfully was reseted")
	return nil
}

func aggHandler(s *state, cmd command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	rssFeed, err := rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	log.Printf("%v", *rssFeed)
	return nil
}

func usersHandler(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return errors.New("no users found")
	}

	log.Printf("users list:")
	for _, u := range users {
		if s.cfg.UserName == u.Name {
			fmt.Printf("%s (current)\n", u.Name)
			continue
		}
		fmt.Printf("%s\n", u.Name)
	}
	return nil
}

func addfeedHandler(s *state, cmd command) error {
	args := cmd.Args
	if len(args) < 2 {
		return errors.New("the login handler expects a two arguments, 1: RSS title; 2: RSS URL")
	}

	u, err := s.db.GetUserByName(context.Background(), s.cfg.UserName)
	if err != nil {
		return err
	}

	e := database.CreateFeedParams {
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: args[0],
		Url: args[1],
		UserID: u.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), e)
	if err != nil {
		return err
	}

	log.Printf("added feed: %v", feed)
	return nil
}

func feedsHandler(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		return errors.New("no feeds")
	}

	log.Printf("list of feeds:")
	for _, f := range feeds {
		user, err := s.db.GetUserByID(context.Background(), f.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("- %s (owner: %s; url: %s)\n", f.Name, user.Name, f.Url)
	}
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	h := c.handler[cmd.Name]
	if h == nil {
		return errors.New("no such command")
	}

	if err := h(s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(s *state, cmd command) error) {
	c.handler[name] = f
}
