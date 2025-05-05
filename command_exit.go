package main

import (
	"fmt"
	"os"
)

func commandExit(conf *config, args ...string) error {
	fmt.Print("Shutting down...Goodbye.\n")
	os.Exit(0)
	return nil
}
