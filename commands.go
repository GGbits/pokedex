package main

import (
	"fmt"
	"math/rand"
	"os"
)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"catch": {
			name:        "catch [pokemon]",
			description: "Attempt to catch [pokemon]",
			callback:    commandCatch,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore [area]",
			description: "Lists pokemon in an area",
			callback:    commandExplore,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"inspect": {
			name:        "inspect [pokemon]",
			description: "Displays stats for [pokemon]",
			callback:    commandInspect,
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

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("could not catch a pokemon, no pokemon specified. Please try 'catch [pokemon]")
	}
	webResult, err := queryCache(cfg, cfg.pokeUrl+args[0])
	if err != nil {
		return err
	}

	pokeResult, err := unmarshallPokemonResult(webResult)
	if err != nil {
		return fmt.Errorf("an error occured while serializing the json: %s", err)
	}

	fmt.Printf("attempting to catch a %v\n", pokeResult.Name)
	pExp := pokeResult.BaseExperience
	rndInt := rand.Intn(500)
	if rndInt > pExp {
		fmt.Printf("You caught a %v\n", pokeResult.Name)
		cfg.caughtPokemon[pokeResult.Name] = pokeResult
	} else {
		fmt.Printf("%v escaped...\n", pokeResult.Name)
	}
	return nil
}

func commandExit(cfg *config, args ...string) error {
	os.Exit(0)
	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("can't explore, no area provided. Please try 'explore [area name]'\n")
	}
	fmt.Printf("Exploring area %v...\n", args)
	exploreUrl := cfg.locBaseUrl + args[0]
	webResult, err := queryCache(cfg, exploreUrl)
	if err != nil {
		return err
	}
	exResult, err := unmarshallPokeExploreResult(webResult)
	if err != nil {
		return err
	}
	fmt.Println("Found the following pokemon...")
	for i := 0; i < len(exResult.PokemonEncounters); i++ {
		pokeName := exResult.PokemonEncounters[i].Pokemon.Name
		fmt.Printf(" -  %v\n", pokeName)
	}
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("\nWelcome to the Pokedex!\nUsage:")
	fmt.Println()
	for _, c := range getCommands() {
		fmt.Printf("%v: %v\n", c.name, c.description)
	}
	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("can't Inspect, no pokemon provided. Please try 'explore [pokemon]'\n")
	}

	pr, ok := cfg.caughtPokemon[args[0]]
	if !ok {
		return fmt.Errorf("you haven't caught that pokemon yet, can't provide stats")
	}
	pokemon := newPokemonInspectInfo(pr)
	pokemon.print()
	return nil
}

func commandMap(cfg *config, args ...string) error {
	webResult, err := queryCache(cfg, cfg.nextUrl)
	if err != nil {
		return err
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

func commandMapb(cfg *config, args ...string) error {
	if cfg.prevUrl == "" {
		return fmt.Errorf("no previous location set exists. Please try using \"map\" at least twice before \"mapb\"")
	}

	webResult, err := queryCache(cfg, cfg.prevUrl)
	if err != nil {
		return err
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
