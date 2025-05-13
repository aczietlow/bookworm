package main

import "fmt"

func commandHelp(conf *config, args ...string) error {
	fmt.Print("Bookworm usage:\n\n")
	for _, c := range registry {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	fmt.Println()
	return nil
}
