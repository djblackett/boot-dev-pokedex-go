package commands

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/djblackett/pokedex-go/internal/pokecache"
	"github.com/djblackett/pokedex-go/internal/pokedex"
)

func CommandMap(config *Config, cache *pokecache.Cache) error {

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

func CommandMapb(config *Config, cache *pokecache.Cache) error {
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

func CommandExplore(name string, config *Config, cache *pokecache.Cache) error {
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
