package commands

type CLICommand struct {
	Name        string
	Description string
	Callback    func(args []string) error
}

type Config struct {
	Next     string
	Previous string
	Count    int
}

// var commands map[string]CLICommand
// var cache *pokecache.Cache

// var config *Config

// func GetCommands(cfg *Config, c *pokecache.Cache, p map[string]pokedex.PokemonResult) map[string]CLICommand {
// 	config = cfg
// 	cache = c
// 	pokedexMap = p

// 	return map[string]CLICommand{
// 		"exit":    {"exit", "Exit the Pokedex", CommandExit},
// 		"help":    {"help", "Displays a help message", CommandHelp},
// 		"map":     {"map", "Displays next 20 locations", func(args []string) error { return commandMap(config) }},
// 		"mapb":    {"mapb", "Displays previous 20 locations", func(args []string) error { return commandMapb(config) }},
// 		"explore": {"explore", "Lists pokemon in area", CommandExplore},
// 		"catch":   {"catch", "Try to catch a pokemon", commandCatch},
// 		"inspect": {"inspect", "View a pokemon's information", commandInspect},
// 		"pokedex": {"pokedex", "Lists all of your pokemon", commandPokedex},
// 	}
// }
