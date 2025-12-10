package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/daitonium/go-blog-aggregator/internal/config"
	"github.com/daitonium/go-blog-aggregator/internal/database"
	"github.com/daitonium/go-blog-aggregator/internal/rss"
	"github.com/google/uuid"
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("No arguments found for command")
	}
	userName := cmd.args[0]
	if _, err := s.db.GetUser(context.Background(), userName); err != nil {
		log.Fatalf("%v, for %s", err, userName)
	}

	if err := s.cfg.SetUser(userName); err != nil {
		return err
	}

	log.Printf("The username %s has been set \n", s.cfg.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("No arguments found for command")
	}
	name := cmd.args[0]
	if len(name) == 0 {
		return errors.New("Name is empty, please add a name to register")
	}
	_, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		return errors.New("Name already exists, try with another name")
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	newUser, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}
	s.cfg.SetUser(newUser.Name)
	log.Println("New user created and set")
	log.Printf("user data: %v \n", newUser)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Println("Users deleted succesfully")
	os.Exit(0)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, usr := range users {
		if usr.Name == s.cfg.CurrentUserName {
			log.Printf("* %s (current)", usr.Name)
		} else {
			log.Printf("* %s", usr.Name)
		}
	}
	os.Exit(0)
	return nil

}

func handlerAggregator(s *state, cmd command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	log.Printf("Feed: %+v\n", feed)

	return nil
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
