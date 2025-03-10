package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aczietlow/addToBookshelf/api"
	"github.com/aczietlow/addToBookshelf/net"
)

type book struct {
	Title    string
	Subtitle string
	Author   []string
	ISBN     []string
	cover    string
	// categories []string
}

func main() {
	// book := "The minority report"
	bookISBN := "9780756419189"
	bookData := getBookByID(bookISBN)
	prettyPrint(bookData)
}

func getBookByID(id string) api.BookResponse {
	var bookData api.BookResponse

	var libraryURL string = "https://openlibrary.org/api/volumes/brief/isbn/" + id + ".json"
	req, err := http.NewRequest("GET", libraryURL, nil)
	if err != nil {
		fmt.Printf("Request failed: %w", err)
	}

	client := &http.Client{
		Transport: &net.OpenLibraryTransport{
			UserAgent: "Add To Bookshelf/0.1 (aczietlow@gmail.com)",
			Transport: http.DefaultTransport,
		},
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Response failed: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&bookData); err != nil {
		fmt.Printf("%s", err)
		log.Fatal("Failed to decode json on test")
	}
	return bookData
}

func prettyPrint(data any) {
	encodedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal("failed to pretty print map")
	}
	fmt.Println(string(encodedData))
}
