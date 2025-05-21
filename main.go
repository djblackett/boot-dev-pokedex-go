package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/djblackett/pokedex-go/commands"
	"github.com/djblackett/pokedex-go/internal/pokecache"
	"github.com/djblackett/pokedex-go/internal/pokedex"
)

var commandsMap map[string]commands.CLICommand

var cache *pokecache.Cache
var pokedexMap map[string]pokedex.PokemonResult

func main() {
	bufio := bufio.NewScanner(os.Stdin)
	config := commands.Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		Count:    0,
	}

	cache = pokecache.NewCache(10 * time.Second)
	pokedexMap = make(map[string]pokedex.PokemonResult)

	commandsMap = map[string]commands.CLICommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback: func(args []string) error {
				return commands.CommandExit()
			},
		},
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback: func(args []string) error {
				return commands.CommandHelp(commandsMap)
			},
		},
		"map": {
			Name:        "map",
			Description: "Displays next 20 locations",
			Callback: func(args []string) error {
				return commands.CommandMap(&config, cache)
			},
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays previous 20 locations",
			Callback: func(args []string) error {
				return commands.CommandMapb(&config, cache)
			},
		},
		"explore": {
			Name:        "explore",
			Description: "Lists pokemon in area",
			Callback: func(args []string) error {
				if len(args) == 2 {
					return commands.CommandExplore(args[1], &config, cache)
				}
				fmt.Println("Command explore requires an argument")
				return nil
			},
		},
		"catch": {
			Name:        "catch",
			Description: "Try to catch a pokemon",
			Callback: func(args []string) error {
				if len(args) >= 2 {
					return commands.CommandCatch(args[1], pokedexMap)
				}
				fmt.Println("Command catch requires an argument")
				return nil
			},
		},
		"inspect": {
			Name:        "inspect",
			Description: "View a pokemon's information",
			Callback: func(args []string) error {
				if len(args) >= 2 {
					return commands.CommandInspect(args[1], pokedexMap)
				}
				fmt.Println("Command inspect requires an argument")
				return nil
			},
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Lists all of your pokemon",
			Callback: func(args []string) error {
				return commands.CommandPokedex(pokedexMap)
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
		if cmd, ok := commandsMap[cmdName]; ok {
			err := cmd.Callback(input)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command:")
		}
	}
}
