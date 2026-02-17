package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/realdnchka/gator-go/internal/config"
	"github.com/realdnchka/gator-go/internal/database"

	_ "github.com/lib/pq"
)



func StateInit() state {
	return state{}
}

func CommandsInit() commands{
	cmds := commands{
		handler: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", resetHandler)
	cmds.register("users", usersHandler)
	cmds.register("agg", middlewareLoggedIn(aggHandler))
	cmds.register("addfeed", middlewareLoggedIn(addfeedHandler))
	cmds.register("feeds", feedsHandler)
	cmds.register("follow", middlewareLoggedIn(followHandler))
	cmds.register("following", middlewareLoggedIn(followingHandler))
	cmds.register("unfollow", middlewareLoggedIn(unfollowHandler))
	cmds.register("browse", middlewareLoggedIn(browseHandler))
	return cmds
}

func main() {
	state := StateInit()
	cmds := CommandsInit()

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	state.cfg = &cfg

	if len(os.Args) <= 1 {
		log.Fatalf("missing command")
	}

	db, err := sql.Open("postgres", state.cfg.DBUrl)
	state.db = database.New(db)

	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	cmd := command {
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err = cmds.run(&state, cmd); err != nil {
		log.Fatalf("error: %v", err)
	}
}
