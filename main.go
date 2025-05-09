package main

import (
	"time"

	"github.com/aczietlow/bookworm/pkg/openlibraryapi"
)

func main() {
	// book := "The minority report"
	// bookISBN := "9780756419189"
	// bookData := getBookByID(bookISBN)
	// prettyPrint(bookData)

	openLibClient := openlibraryapi.NewClient(time.Minute * 5)
	conf := &config{
		apiClient: openLibClient,
	}
	startCli(conf)

}
