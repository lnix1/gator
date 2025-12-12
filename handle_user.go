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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Not enough args")
	}
	
	currUser, dbCheck := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if dbCheck != nil {
		return fmt.Errorf("user does not exist: %w", dbCheck)
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
	fmt.Printf("%+v", i)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	currFeeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to retrieve feeds list: %w", err)
	}

	for _, feed := range currFeeds {
		fmt.Printf("%+v", feed)
	}
	return nil
}
