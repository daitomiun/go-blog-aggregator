package main

import (
	"errors"
	"github.com/daitonium/go-blog-aggregator/internal/config"
	"github.com/daitonium/go-blog-aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.registeredCommands[cmd.name]
	if !exists {
		return errors.New("Command not found")
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
