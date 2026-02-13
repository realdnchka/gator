package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/realdnchka/gator-go/internal/config"
	"github.com/realdnchka/gator-go/internal/database"

	_ "github.com/lib/pq"
)

var cmds commands = commands{
	handler: make(map[string]func(*state, command) error),
}

func main() {
	state := StateInit()
	CommandsInit()

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	state.cfg = &cfg

	if len(os.Args) <= 1 {
		log.Fatalf("missing command")
	}
	if len(os.Args) == 2 {
		log.Fatalf("missing parameter")
	}
	db, err := sql.Open("postgres", state.cfg.DBUrl)
	state.db = database.New(db)

	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err = cmds.run(&state, cmd); err != nil {
		log.Fatalf("could not run a command: %v", err)
	}
}
