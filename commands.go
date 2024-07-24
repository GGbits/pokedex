package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func commandExit() error {
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("\nWelcome to the Pokedex!\nUsage:")
	fmt.Println()
	for _, c := range getCommands() {
		fmt.Printf("%v: %v\n", c.name, c.description)
	}
	return nil
}

func commandMap() error {
	//TODO: commandMap functionality
	res, err := http.Get(nextMap)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		return fmt.Errorf("Response failed with status code: %d and\nbody:%s\n", res.StatusCode, body)
	}

	plResult := pokeLocationResult{}
	err = json.Unmarshal(body, &plResult)
	if err != nil {
		return fmt.Errorf("failed to convert json to location data: %s", err)
	}

	for _, res := range plResult.Results {
		println(res.Name)
	}
	prevMap = nextMap
	nextMap = plResult.Next
	return nil
}

func commandMapb() error {
	//TODO: commandMapb functionality
	return nil
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
