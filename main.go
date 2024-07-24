package main

import (
	"bufio"
	"fmt"
	"os"
)

var nextMap string = "https://pokeapi.co/api/v2/location/"
var prevMap string

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type pokeLocationResult struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
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
				fmt.Print(err)
			}
		}
		fmt.Printf("\n\nPokedex > ")
	}
}
