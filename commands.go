package main

import (
	"errors"
	"github.com/daitonium/go-blog-aggregator/internal/config"
	"log"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("No arguments found for command")
	}
	userName := cmd.args[0]

	if err := s.config.SetUser(userName); err != nil {
		return err
	}

	log.Printf("The username %s has been set \n", s.config.CurrentUserName)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.registeredCommands[cmd.name]
	if !exists {
		log.Println("test exists")
		return errors.New("Command not found")
	}
	return handler(s, cmd)
}
func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
