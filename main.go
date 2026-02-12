package main

import (
	"log"
	"os"

	"github.com/realdnchka/gator-go/internal/config"
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
	state.Config = &cfg

	if len(os.Args) <= 1 {
		log.Fatalf("missing command")
	}

	if len(os.Args) == 2 {
		log.Fatalf("missing parameter")
	}

	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	cmds.run(&state, cmd)
}
