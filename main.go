package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/daitonium/go-blog-aggregator/internal/config"
	"github.com/daitonium/go-blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Cannot read config: %v", err)
	}
	s := &state{}
	s.cfg = &cfg
	cmds := &commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	if len(os.Args) < 2 {
		log.Fatal("Less than 2 arguments")
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	s.db = dbQueries

	name, args := os.Args[1], os.Args[2:]
	if err := cmds.run(s, command{name, args}); err != nil {
		log.Fatal(err)
	}

}
