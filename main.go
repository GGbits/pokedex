package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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

type config struct {
	nextUrl string `default:"https://pokeapi.co/api/v2/location/"`
	prevUrl string
}

func main() {
	fmt.Println("Welcome to Pokedex! For help, please type \"help\" without quotes.")
	fmt.Println("To quit, simply type \"quit\"")
	fmt.Println()
	fmt.Printf("Pokedex > ")

	cfg := config{nextUrl: "https://pokeapi.co/api/v2/location/"}
	scanner := bufio.NewScanner(os.Stdin)
	commandMap := getCommands()
	for scanner.Scan() {
		cm, ok := commandMap[scanner.Text()]
		if !ok {
			fmt.Printf("unknown command: %v. See help for list of commands", scanner.Text())
		} else {
			err := cm.callback(&cfg)
			if err != nil {
				fmt.Print(err)
			}
		}
		fmt.Printf("\n\nPokedex > ")
	}
}
