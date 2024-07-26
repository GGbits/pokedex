package main

import (
	"fmt"
	"os"
)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Gets next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Gets last 20 locations",
			callback:    commandMapb,
		},
	}
}

func commandExit(cfg *config) error {
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("\nWelcome to the Pokedex!\nUsage:")
	fmt.Println()
	for _, c := range getCommands() {
		fmt.Printf("%v: %v\n", c.name, c.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	//TODO: Make this a function
	webResult, ok := cfg.cache.Get(cfg.nextUrl)
	if !ok {
		wr, err := getApiResponse(cfg.nextUrl)
		if err != nil {
			return fmt.Errorf("an error occured while getting API response: %s", err)
		}
		cfg.cache.Add(cfg.nextUrl, wr)
		webResult = wr
	}

	plResult, err := unmarshallPokeLocationResult(webResult)
	if err != nil {
		return fmt.Errorf("an error occured while serializing the json: %s", err)
	}

	for _, res := range plResult.Results {
		println(res.Name)
	}
	cfg.prevUrl = plResult.Previous
	cfg.nextUrl = plResult.Next
	return nil
}

func commandMapb(cfg *config) error {
	if cfg.prevUrl == "" {
		return fmt.Errorf("no previous location set exists. Please try using \"map\" at least twice before \"mapb\"")
	}

	//TODO: Make this a function
	webResult, ok := cfg.cache.Get(cfg.prevUrl)
	if !ok {
		wr, err := getApiResponse(cfg.prevUrl)
		if err != nil {
			return fmt.Errorf("an error occured while getting API response: %s", err)
		}
		cfg.cache.Add(cfg.prevUrl, wr)
		webResult = wr
	}

	webResult, err := getApiResponse(cfg.prevUrl)
	if err != nil {
		return fmt.Errorf("an error occured while getting API response. Could be internet connection?: %s", err)
	}
	plResult, err := unmarshallPokeLocationResult(webResult)
	if err != nil {
		return fmt.Errorf("an error occured while serializing the json: %s", err)
	}

	for _, res := range plResult.Results {
		println(res.Name)
	}
	cfg.prevUrl = plResult.Previous
	cfg.nextUrl = plResult.Next
	return nil
}
