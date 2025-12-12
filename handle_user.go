package main

import (
	"fmt"
	"time"
	"context"
	"github.com/google/uuid"
	"github.com/lnix1/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	_, dbCheck := s.db.GetUser(context.Background(), name)
	if dbCheck != nil {
		return fmt.Errorf("user does not exist: %w", dbCheck)
	}

	err := s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	createArgs := database.CreateUserParams{
		ID: uuid.New(), 
		CreatedAt: time.Now(), 
		UpdatedAt: time.Now(),
		Name: name,
	}

	i, err := s.db.CreateUser(context.Background(), createArgs)
	if err != nil {
		return fmt.Errorf("failed to create user in db: %w", err)
	}

	err = handlerLogin(s, cmd)
	if err != nil {
		return err
	}

	fmt.Printf("User was created with ID: %s, CreatedAt: %s, UpdatedAt: %s, Name: %s \n", 
		i.ID, i.CreatedAt, i.UpdatedAt, i.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to reset Users table: %w", err)
	}

	fmt.Println("Successfully reset Users table.")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	currUsers, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to retrieve users list: %w", err)
	}

	for _, user := range currUsers {
		if user == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current) \n", user)
			continue
		}
		fmt.Printf("* %s \n", user)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	//targetUrl := cmd.Args[0]
	targetUrl := "https://www.wagslane.dev/index.xml"
	Feed, err := fetchFeed(context.Background(), targetUrl)
	if err != nil {
		return fmt.Errorf("Error fetching RSS Feed: %w", err)
	}

	fmt.Println(*Feed)

	return nil
}

func handlerAddFeed(s *state, cmd command, currUser database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Not enough args")
	}
	
	currentUserId := currUser.ID
	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	createArgs := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: feedName,
		Url: feedUrl,
		UserID: currentUserId,
	}

	i, err := s.db.CreateFeed(context.Background(), createArgs)
	if err != nil {
		return fmt.Errorf("failed to create feed in db: %w", err)
	}
	fmt.Printf("%+v \n", i)

	err = handlerFollow(s, command{Name: "follow", Args: []string{feedUrl}}, currUser)
	if err != nil {
		return fmt.Errorf("error registering current user to follow newly added feed: %w", err)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	currFeeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to retrieve feeds list: %w", err)
	}

	for _, feed := range currFeeds {
		fmt.Printf("%+v \n", feed)
	}
	return nil
}

func handlerFollow(s *state, cmd command, currUser database.User) error {
	currentUserId := currUser.ID
	
	targetFeedId, err := s.db.GetFeedId(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("feed url does not exist in db: %w", err)
	}

	createArgs := database.CreateFeedFollowsParams{
		CreatedAt:	time.Now(),
		UpdatedAt: 	time.Now(),
		UserID: 	currentUserId,
		FeedID:    	targetFeedId,
	}
	feedFollowsRow, err := s.db.CreateFeedFollows(context.Background(), createArgs)
	if err != nil {
		return fmt.Errorf("error adding feed_follows instance for user: %w", err)
	}

	fmt.Printf("User %s now following Feed %s \n", feedFollowsRow.UserName, feedFollowsRow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, currUser database.User) error {
	currentUserId := currUser.ID

	currentUserFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), currentUserId)
	if err != nil {
		return fmt.Errorf("error retrieving current user's feed follows: %w", err)
	}

	fmt.Printf("%s is following these feeds: \n", s.cfg.CurrentUserName)
	for _, feed := range currentUserFeeds {
		fmt.Printf(" - %s \n", feed.FeedName)
	}
	
	return nil
}

func handlerUnfollow(s *state, cmd command, currUser database.User) error {
	currentUserId := currUser.ID
	
	targetFeedId, err := s.db.GetFeedId(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("feed url does not exist in db: %w", err)
	}
	
	deleteArgs:= database.RemoveFeedFollowParams{
		UserID: 	currentUserId,
		FeedID:		targetFeedId,
	}

	err = s.db.RemoveFeedFollow(context.Background(), deleteArgs)
	if err != nil {
		return fmt.Errorf("error deleting feed_follows instance for user: %w", err)
	}

	fmt.Printf("User %s has unfollowed Feed %s \n", currUser.Name, cmd.Args[0])

	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, currUser database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currUser, dbCheck := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if dbCheck != nil {
			return fmt.Errorf("user does not exist: %w", dbCheck)
		}

		return handler(s, cmd, currUser)
	}
}
