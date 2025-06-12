package main

import (
	"time"

	"github.com/aczietlow/bookworm/pkg/openlibraryapi"
)

func main() {
	openLibClient := openlibraryapi.NewClient(time.Second*10, time.Minute*5)
	conf := &config{
		apiClient: openLibClient,
		registry:  registerCommands(),
	}
	cli := NewTui(conf)

	cli.Run()
}
