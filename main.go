package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/amalrajan30/pokedexcli/inernal/pokeapi"
	"github.com/amalrajan30/pokedexcli/inernal/pokecache"
	"github.com/amalrajan30/pokedexcli/inernal/pokedex"
)

const DEBUG_CATCHALL = true

func printPrompt() {
	fmt.Print("Pokedex > ")
}

func commandHelp(_ *config, _ ...string) error {
	fmt.Printf("\n Welcome to the Pokedex! \n Usage: \n\n")
	println("help: Displays a help message")
	println("exit: Exit the Pokedex")
	println("map: Displays the names of next 20 location areas")
	println("mapb: Displays the names of previous 20 location areas")
	fmt.Printf("\n\n")

	return nil
}

func commandExit(_ *config, _ ...string) error {
	os.Exit(0)

	return nil
}

func commandMap(config *config, _ ...string) error {
	if config.next == "" {
		return errors.New("you have reached the end")
	}

	data, ok := pokeapi.GetLocationArea(config.next, config.cache)

	if !ok {
		return errors.New("something went wrong")
	}

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	config.next = data.Next
	if data.Previous != nil {
		config.previous = *data.Previous
	} else {
		config.previous = ""
	}

	return nil
}

func commandMapB(config *config, _ ...string) error {
	if config.previous == "" {
		return errors.New("you are at the beginning")
	}

	data, ok := pokeapi.GetLocationArea(config.previous, config.cache)

	if !ok {
		return errors.New("something went wrong")
	}

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	config.next = data.Next
	if data.Previous != nil {
		config.previous = *data.Previous
	} else {
		config.previous = ""
	}

	return nil
}

func commandExplore(config *config, cmd ...string) error {

	if (len(cmd) < 2) {
		return errors.New("please provide an area name")
	}

	data, ok := pokeapi.ExplorePokemon(cmd[1], config.cache)

	if !ok {
		return errors.New("something went wrong")
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range data.PokemonEncounters {
		fmt.Printf(" - %s \n",pokemon.Pokemon.Name)
	}

	return nil
}

func commandPokedex(config *config, _ ...string) error {

	dex := config.dex.Dex

	fmt.Println("Your Pokedex:")
	for _, pokemon := range dex {
		fmt.Printf(" - %s \n", pokemon.Name)
	}

	return nil
}

func calculateChance(baseExp float64) int {
	// Normalize x to a probability range
	prob := 1.0 / (1.0 + baseExp)

	// Generate a random number between 0 and 1
	r := rand.Float64()

	// Return 1 if the random number is less than the probability
	if r < prob {
		return 1
	}
	return 0
	
}

func commandCatch(config *config, cmd ...string,) error {

	if (len(cmd) < 2) {
		return errors.New("please give a pokemon to catch")
	}

	data, ok := pokeapi.GetPokemon(cmd[1], config.cache)

	if !ok {
		return errors.New("something went wrong")
	}


	fmt.Printf("Throwing a Pokeball at %s \n", cmd[1])

	chance := calculateChance(float64(data.BaseExperience))

	if DEBUG_CATCHALL {
		chance = 1
	}

	if chance == 1 {
		fmt.Printf("%s was caught! \n", cmd[1])
		config.dex.AddPokemon(data, cmd[1])
		fmt.Println("You may now inspect it with the inspect command")
	} else {
		fmt.Printf("%s escaped! \n", cmd[1])
	}

	return nil
}

func commandInspect(config *config, cmd ...string,) error {

	if (len(cmd) < 2) {
		return errors.New("please give a pokemon to inspect")
	}

	pokemon, ok := config.dex.Dex[cmd[1]]

	if !ok {
		return errors.New("you have not caught that pokemon")
	}

	statInfo := ""
	typesInfo := ""

	for _, stat := range pokemon.Stats {
		statInfo = strings.Join([]string{statInfo, fmt.Sprintf("  - %s: %v", stat.Stat.Name, stat.BaseStat)}, "\n")
	}

	for _, types := range pokemon.Types {
		typesInfo = strings.Join([]string{typesInfo, fmt.Sprintf("  - %s", types.Type.Name,)}, "\n")
	}

	pokemonInfo := fmt.Sprintf("Name: %s \nHeight: %v \nWeight: %v \nStats: %s \n Types:%s \n ", pokemon.Name, pokemon.Height, pokemon.Weight, statInfo, typesInfo)

	fmt.Println(pokemonInfo)

	return nil
}

type cliCommand struct {
	name 		string
	description string
	callback 	func(config *config, params ...string) error
}

type config struct {
	next string
	previous string
	cache *pokecache.Cache
	dex pokedex.PokeDex
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"exit": {
			name: "exit",
			description: "Exit the pokedex",
			callback: commandExit,
		},
		"map": {
			name: "map",
			description: "Displays the names of next 20 location areas",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Displays the names of previous 20 location areas",
			callback: commandMapB,
		},
		"explore": {
			name: "explore",
			description: "",
			callback: commandExplore,
		},
		"catch": {
			name: "explore",
			description: "",
			callback: commandCatch,
		},
		"inspect": {
			name: "explore",
			description: "",
			callback: commandInspect,
		},
		"pokedex": {
			name: "explore",
			description: "",
			callback: commandPokedex,
		},
	}
}

func main(){
	reader := bufio.NewScanner(os.Stdin)
	printPrompt()

	duration, _ := time.ParseDuration("7s")

	cache := pokecache.NewCache(duration)
	userDex := pokedex.NewDex()

	locationAreaConfig := &config{
		next: "https://pokeapi.co/api/v2/location-area/",
		previous: "",
		cache: cache,
		dex: userDex,
	}

	commands := getCommands()

	for reader.Scan() {
		text := reader.Text()
		words := strings.Fields(text)
		if cmd, exists := commands[words[0]]; exists {
			cmdErr := cmd.callback(locationAreaConfig, words...)

			if cmdErr != nil {
				fmt.Println(cmdErr.Error())
			}
		}

		printPrompt()
	}
}