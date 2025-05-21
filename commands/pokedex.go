package commands

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"

	"github.com/djblackett/pokedex-go/internal/pokedex"
)

func CommandCatch(name string, pokedexMap map[string]pokedex.PokemonResult) error {
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
func CommandInspect(name string, pokedexMap map[string]pokedex.PokemonResult) error {

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
		fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}
	return nil
}

func CommandPokedex(pokedexMap map[string]pokedex.PokemonResult) error {
	fmt.Println("Your Pokedex:")
	for name := range pokedexMap {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}
