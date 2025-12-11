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
