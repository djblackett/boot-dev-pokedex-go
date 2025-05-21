package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/djblackett/pokedex-go/internal/pokecache"
	"github.com/djblackett/pokedex-go/internal/pokedex"
)

var commands map[string]cliCommand

var cache *pokecache.Cache
var pokedexMap map[string]pokedex.PokemonResult

func main() {
	bufio := bufio.NewScanner(os.Stdin)
	config := Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		Count:    0,
	}

	cache = pokecache.NewCache(10 * time.Second)
	pokedexMap = make(map[string]pokedex.PokemonResult)

	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback: func(args []string) error {
				return commandExit()
			},
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback: func(args []string) error {
				return commandHelp()
			},
		},
		"map": {
			name:        "map",
			description: "Displays next 20 locations",
			callback: func(args []string) error {
				return commandMap(&config)
			},
		},
		"mapb": {
			name:        "map",
			description: "Displays previous 20 locations",
			callback: func(args []string) error {
				return commandMapb(&config)
			},
		},
		"explore": {
			name:        "explore",
			description: "Lists pokemon in area",
			callback: func(args []string) error {
				if len(args) == 2 {
					return commandExplore(args[1])
				}
				fmt.Println("Command explore requires an argument")
				return nil
			},
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon",
			callback: func(args []string) error {
				if len(args) >= 2 {
					return commandCatch(args[1])
				}
				fmt.Println("Command catch requires an argument")
				return nil
			},
		},
		"inspect": {
			name:        "inspect",
			description: "View a pokemon's information",
			callback: func(args []string) error {
				if len(args) >= 2 {
					return commandInspect(args[1])
				}
				fmt.Println("Command inspect requires an argument")
				return nil
			},
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all of your pokemon",
			callback: func(args []string) error {
				return commandPokedex()
			},
		},
	}

	fmt.Println("Welcome to the Pokedex!")
	for {
		fmt.Print("Pokedex > ")
		if !bufio.Scan() {
			break
		}
		input := strings.Fields(strings.TrimSpace(bufio.Text()))
		if len(input) == 0 {
			continue
		}
		cmdName := input[0]
		if cmd, ok := commands[cmdName]; ok {
			err := cmd.callback(input)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command:")
		}
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

type Config struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Count    int    `json:"count"`
}

func CleanInput(text string) []string {
	loweredText := strings.ToLower(text)
	trimmedText := strings.TrimSpace(loweredText)
	words := strings.Fields(trimmedText)
	return words
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(config *Config) error {

	var res *http.Response
	var err error
	var locationAreas pokedex.LocationAreaSmall

	i, ok := cache.Get(config.Next)
	if !ok {
		res, err = http.Get(config.Next)

		if err != nil {
			return err
		}

		if err := json.NewDecoder(res.Body).Decode(&locationAreas); err != nil {
			return err
		}

		defer res.Body.Close()
	} else {
		fmt.Println("Reading from cache...")
		locationAreas = pokedex.LocationAreaSmall{}
		if err := json.Unmarshal(i.Val, &locationAreas); err != nil {
			return err
		}
	}

	locationBytes, err := json.Marshal(locationAreas)
	if err != nil {
		fmt.Println("Invalid data")
	}
	cache.Add(config.Next, locationBytes)
	config.Next = locationAreas.Next

	config.Previous = locationAreas.Previous

	config.Count++

	for _, location := range locationAreas.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(config *Config) error {
	var res *http.Response
	var err error
	var locationAreas pokedex.LocationAreaSmall

	if config.Previous == "" || config.Count == 0 {
		fmt.Println("you're on the first page")
		return nil
	}

	i, ok := cache.Get(config.Previous)
	if !ok {

		res, err = http.Get(config.Previous)

		if err != nil {
			return err
		}

		if err := json.NewDecoder(res.Body).Decode(&locationAreas); err != nil {
			return err
		}

		defer res.Body.Close()
	} else {
		fmt.Println("Reading from cache...")
		locationAreas = pokedex.LocationAreaSmall{}
		if err := json.Unmarshal(i.Val, &locationAreas); err != nil {
			return err
		}
	}
	config.Next = locationAreas.Next

	config.Previous = locationAreas.Previous

	config.Count--

	for _, location := range locationAreas.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(name string) error {
	var res *http.Response
	var err error
	baseURL := "https://pokeapi.co/api/v2/location-area/"
	fullURL := baseURL + name
	var pokemonNames []string
	var locationArea pokedex.LocationArea

	nameList, ok := cache.Get(fullURL)
	if !ok {

		res, err = http.Get(fullURL)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if err := json.NewDecoder(res.Body).Decode(&locationArea); err != nil {
			return err
		}

		pokemonNames = make([]string, 0, len(locationArea.PokemonEncounters))

		for _, encounter := range locationArea.PokemonEncounters {
			pokemonNames = append(pokemonNames, encounter.Pokemon.Name)
		}

		pokeBytes, err := json.Marshal(pokemonNames)
		if err != nil {
			fmt.Println("error marshaling json", err)
			return err
		}

		cache.Add(fullURL, pokeBytes)

		fmt.Printf("Exploring %s...\n", name)
		fmt.Println("Found Pokemon:")
		for _, pokemon := range locationArea.PokemonEncounters {
			fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
		}

	} else {
		fmt.Println("Reading pokemon from cache")
		var nameBytes []string
		err = json.Unmarshal(nameList.Val, &nameBytes)
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, pokemon := range nameBytes {
			fmt.Printf(" - %s\n", pokemon)
		}
	}
	return nil
}

func commandCatch(name string) error {
	baseURL := "https://pokeapi.co/api/v2/pokemon/"
	fullURL := baseURL + name

	var res *http.Response
	var err error

	res, err = http.Get(fullURL)
	if err != nil {
		return err
	}
	var pokemon pokedex.PokemonResult
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&pokemon)
	if err != nil {
		fmt.Println("Error decoding json:", err)
		return err
	}

	baseExp := pokemon.BaseExperience

	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	randomInt := rand.IntN(baseExp)
	if randomInt%3 == 0 {
		// catch pokemon
		fmt.Printf("%s was caught!\n", name)
		pokedexMap[pokemon.Species.Name] = pokemon
		return nil
	}
	fmt.Printf("%s escaped!\n", name)
	return nil
}

func commandInspect(name string) error {
	pokemon, ok := pokedexMap[name]
	if !ok {
		fmt.Println("You do not have this pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")

	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, pType := range pokemon.Types {
		fmt.Printf(" - %s\n", pType.Type.Name)
	}
	return nil
}

func commandPokedex() error {
	fmt.Println("Your Pokedex:")
	for key := range pokedexMap {
		fmt.Printf(" - %s\n", key)
	}
	return nil
}
