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

func aggHandler(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("expected 1 argument: time between requests")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	
	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s, user)
	}
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

func addfeedHandler(s *state, cmd command, user database.User) error {
	args := cmd.Args
	if len(args) < 2 {
		return errors.New("the addfeed command expects two arguments, 1: RSS title; 2: RSS URL")
	}

	e := database.CreateFeedParams {
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: args[0],
		Url: args[1],
		UserID: user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), e)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams {
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	log.Printf("added feed: %v", feed)
	return nil
}

func feedsHandler(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil || len(feeds) == 0 {
		return errors.New("feeds not found")
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

func followHandler(s *state, cmd command, user database.User) error {
	url := cmd.Args[0]
	db := s.db

	feed, err := db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return errors.New("no such feeed")
	}
	params := database.CreateFeedFollowParams {
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}
	_, err = db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("created feed: %s with url: %s", feed.Name, url)
	return nil
}

func followingHandler(s *state, cmd command, user database.User) error {
	db := s.db

	feeds, err := db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	if len(feeds) == 0 {
		fmt.Printf("%s doesn't follow any RSS\n", user.Name)
		return nil
	}
	fmt.Printf("%s follows next RSS:\n", user.Name)
	for _, f := range feeds {
		fmt.Printf("%s\n", f.FeedName)
	}
	return nil
}

func unfollowHandler(s *state, cmd command, user database.User) error {
	db := s.db
	url := cmd.Args[0]
	err := db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		Url: url,
		UserID: user.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("you successfully unsubcribed from %s\n", url)
	return nil
}

func browseHandler(s *state, cmd command, user database.User) error {
	db := s.db
	var limit int32
	if len(cmd.Args) == 0 {
		limit = 2
	}
	posts, err := db.GetPostsByLimit(context.Background(), database.GetPostsByLimitParams{
		UserID: user.ID,
		Limit: limit,
	})
	if err != nil {
		return err
	}
	for _, p := range posts {
		fmt.Printf("%s\n%s\n%v\n", p.Title, p.Description, p.PublishedAt)
	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error{
	f := func(s *state, cmd command) error {
		db := s.db
		user, err := db.GetUserByName(context.Background(), s.cfg.UserName)
		if err != nil {
			return err
		}
		handler(s, cmd, user)
		return nil
	}
	return f
}

func scrapeFeeds(s *state, user database.User) {
	db := s.db
	feedtf, err := db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Fatalf("cannot get next feed to fetch: %v", err)
	}
	err = db.MarkFeedFetched(context.Background(), feedtf.ID)
	if err != nil {
		log.Fatalf("cannot mark feed: %v", err)
	}

	feed, err := rss.FetchFeed(context.Background(), feedtf.Url)
	if err != nil {
		log.Fatalf("cannot fetch feed: %v", err)
	}
	
	for _, f := range feed.Channel.Items {
		pubtime, err := time.Parse(time.RFC1123, f.PubDate)
		if err != nil {
			log.Fatal(err)
		}
		
		posts, err := db.GetPosts(context.Background(), user.ID)
		if err != nil {
			log.Fatalf("cannot get posts: %v", err)
		}

		for _, p := range posts {
			p.Url = f.Link
			continue
		}

		post := database.CreatePostParams {
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: f.Title,
			Url: f.Link,
			Description: f.Description,
			PublishedAt: pubtime,
			FeedID: feedtf.ID,
		}
		
		_, err = db.CreatePost(context.Background(), post)
		if err != nil {
			log.Fatalf("cannot create a post: %v", err)
		}
	}
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
