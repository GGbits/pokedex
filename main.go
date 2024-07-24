package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandHelp() error {
	fmt.Println("\nWelcome to the Pokedex!\nUsage:")
	fmt.Println()
	for _, c := range getCommands() {
		fmt.Printf("%v: %v\n", c.name, c.description)
	}
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func main() {
	fmt.Println("Welcome to Pokedex! For help, please type \"help\" without quotes.")
	fmt.Println("To quit, simply type \"quit\"")
	fmt.Println()
	fmt.Printf("Pokedex > ")
	scanner := bufio.NewScanner(os.Stdin)
	commandMap := getCommands()
	for scanner.Scan() {
		s, ok := commandMap[scanner.Text()]
		if !ok {
			fmt.Printf("unknown command: %v. See help for list of commands", scanner.Text())
		} else {
			err := s.callback()
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Printf("\n\nPokedex > ")
	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}
