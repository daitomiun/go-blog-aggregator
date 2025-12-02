package main

import (
	"log"
	"os"

	"github.com/daitonium/go-blog-aggregator/internal/config"
)

func main() {

	newConf, err := config.Read()
	if err != nil {
		log.Fatalf("Cannot read config: %v", err)
	}
	s := &state{}
	s.config = &newConf
	cmds := &commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Less than 2 arguments")
	}
	name, args := os.Args[1], os.Args[2:]
	if err := cmds.run(s, command{name, args}); err != nil {
		log.Fatal(err)
	}
}
