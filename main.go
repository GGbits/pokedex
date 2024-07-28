package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedex/internal"
	"strings"
	"time"
)

func main() {
	fmt.Println("Welcome to Pokedex! For help, please type \"help\" without quotes.")
	fmt.Println("To quit, simply type \"quit\"")
	fmt.Println()
	fmt.Printf("Pokedex > ")

	cfg := config{
		nextUrl:       "https://pokeapi.co/api/v2/location-area/",
		locBaseUrl:    "https://pokeapi.co/api/v2/location-area/",
		pokeUrl:       "https://pokeapi.co/api/v2/pokemon/",
		caughtPokemon: make(map[string]pokemonResult),
		cache:         internal.NewCache(15 * time.Second),
	}
	scanner := bufio.NewScanner(os.Stdin)
	commandMap := getCommands()
	for scanner.Scan() {
		words := cleanInput(scanner.Text())
		cm, ok := commandMap[words[0]]
		if !ok {
			fmt.Printf("unknown command: %v. See help for list of commands", scanner.Text())
		} else {
			if len(words) == 1 {
				err := cm.callback(&cfg)
				if err != nil {
					fmt.Print(err)
				}
			} else {
				err := cm.callback(&cfg, words[1:]...)
				if err != nil {
					fmt.Print(err)
				}
			}

		}
		fmt.Printf("\n\nPokedex > ")
	}
}

func cleanInput(readLine string) []string {
	lower := strings.ToLower(readLine)
	split := strings.Fields(lower)
	return split
}
