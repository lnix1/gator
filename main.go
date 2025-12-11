package main

import (
	"fmt"
	"os"
	"log"
	"github.com/lnix1/gator/internal/config"
)

type state struct {
	config 	*config.Config
}

type command struct {
	name	string
	args 	[]string
}

type commands struct {
	commandMap 	map[string]func(*state, command) error
}

func (c commands) run(s *state, cmd command) error {
	commandFunc, ok := c.commandMap[cmd.name]
	if !ok {
		return fmt.Errorf("Command does not exist")
	}

	err := commandFunc(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c commands) register(name string, f func(*state, command)	error) {
	c.commandMap[name] = f
}

func handleLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Fatal("Login command expects a single arg, username \n")
	}

	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Println("User has been set")
	return nil
}

func main() {
	currentConfig, _ := config.Read()
	currentState := state{config: &currentConfig}
	currCommands := commands{commandMap: map[string]func(*state, command) error {}}
	currCommands.register("login", handleLogin)
	userArgs := os.Args[1:]
	if len(userArgs) < 1 {
		log.Fatal("not enough arguments passed")
	}
	currCommand := command{name: userArgs[0], args: userArgs[1:]}
	err := currCommands.run(&currentState, currCommand)
	if err != nil {
		fmt.Printf("Got an error running command: %v", err)
	}
}
