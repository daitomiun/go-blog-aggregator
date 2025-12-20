package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/daitonium/go-blog-aggregator/internal/database"
	"github.com/daitonium/go-blog-aggregator/internal/rss"
	"github.com/google/uuid"
)

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

	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}
	s.cfg.SetUser(user.Name)
	log.Println("New user created and set")
	log.Printf("user data: %v \n", user)

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
	if len(cmd.args) <= 0 {
		return errors.New("Needs at least 1 argument duration (1s, 1m, 1h)")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	log.Printf("Fetching feed at the lighting speed of %s", timeBetweenRequests.String())
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}
func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Fatalf("next feed err: %s", err)
		return
	}
	log.Println("Found a new feed to fetch")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {

	err := db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		UpdatedAt:     time.Now(),
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:            feed.ID,
	})
	fetchedFeed, err := rss.FetchFeed(context.Background(), feed.Url)

	if err != nil {
		log.Printf("Fetch Feed err: %s", err)
		return
	}
	log.Printf("Feed from database: %+v\n", feed.Url)
	log.Printf("Feed: %+v\n", fetchedFeed.Channel.Title)
	log.Println("Save posts...")

	for _, item := range fetchedFeed.Channel.Items {
		publishedDateParse, err := parsePublishedAt(item.PubDate)
		if err != nil {
			log.Printf("Error trying to parse rss time %s \n", err)
			continue
		}

		_, err = db.GetPostByUrl(context.Background(), item.Link)

		if err == nil {
			continue
		}

		err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Description: sql.NullString{String: item.Description, Valid: true},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Url:         item.Link,
			PublishedAt: publishedDateParse,
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Printf("Cannot create post %s continue with next one \n", err)
			continue
		}
		log.Println("Post Created from feed")
	}
}

func parsePublishedAt(s string) (time.Time, error) {
	var layouts = []string{
		time.RFC3339,
		time.RFC1123Z,
		time.RFC1123,
		"Mon, 2 Jan 2006 15:04:05 -0700",
	}
	var lastErr error
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}
	log.Printf("--- List of all feeds ---")
	for _, feed := range feeds {
		log.Printf("Name: %s", feed.Name)
		log.Printf("URL: %s", feed.Url)
		log.Printf("User name: %s", feed.UserName)
		fmt.Println()
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("Not sufficient arguments for command (needs title and url)")
	}

	nameArg := cmd.args[0]

	if len(nameArg) == 0 {
		return errors.New("Name cannot be null")
	}

	urlArg := cmd.args[1]

	if len(urlArg) == 0 {
		return errors.New("Url cannot be null")
	}

	feedParam := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      nameArg,
		Url:       urlArg,
		UserID:    user.ID,
	}

	newFeed, err := s.db.CreateFeed(context.Background(), feedParam)
	if err != nil {
		return err
	}
	log.Println(newFeed)

	FeedFollowArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    newFeed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), FeedFollowArgs)
	if err != nil {
		return err
	}
	log.Println("New Follow added!")
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("command takes at least 1 url argument")
	}
	urlArg := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), urlArg)
	if err != nil {
		return err
	}

	FeedFollowArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), FeedFollowArgs)

	if err != nil {
		return err
	}

	log.Println("Feed Follow created")
	log.Println(feedFollow)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	log.Printf("--- Current user `%s` follows --- \n", user.Name)
	for _, feed := range feeds {
		log.Printf("Name %s  \n", feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("command takes at least 1 url argument")
	}
	urlArg := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), urlArg)
	if err != nil {
		return err
	}

	err = s.db.UnfollowFeedByUserIdAndFeedId(
		context.Background(),
		database.UnfollowFeedByUserIdAndFeedIdParams{
			UserID: user.ID,
			FeedID: feed.ID,
		})
	if err != nil {
		return fmt.Errorf("Cannot unfollow feed: %v", err)
	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	if len(cmd.args) < 1 {
		limit = 2
		log.Printf("Defaulting to %v posts per browse \n", limit)
	} else {

		newLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			log.Printf("err: %v", err)
		}
		limit = int32(newLimit)

		log.Printf("%v posts per browse \n", limit)
	}

	posts, err := s.db.GetPostsByUserId(context.Background(), database.GetPostsByUserIdParams{UserID: user.ID, Limit: limit})
	if err != nil {
		log.Println("Cannot look at posts")
		return err
	}
	log.Println("------ Browse your followed feeds! ------")
	for _, post := range posts {
		fmt.Printf("--- %v ---\n", post.Title)
		fmt.Printf("Desc: %v \n", post.Description.String)
		fmt.Printf("Published on: %v \n", post.PublishedAt)
		fmt.Printf("Url: %v \n", post.Url)
		fmt.Println("---------")

	}
	return nil
}
