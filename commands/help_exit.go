package commands

import (
	"fmt"
	"os"
)

func CommandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(commands map[string]CLICommand) error {
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
	}
	return nil
}
