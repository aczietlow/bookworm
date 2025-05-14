package openlibraryapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type book struct {
	Title    string
	Subtitle string
	Author   []string
	Summary  string
	ISBN     string
	Genre    []string
	Cover    string
	Source   string
}

func (c *Client) GetBookById(id string) (book, error) {
	id = strings.ToUpper(id)

	lr, err := getBookDetails(id, c.httpClient)
	if err != nil {
		return book{}, err
	}

	b := aggregateLibraryRecord(lr)

	return b, nil

}

func aggregateLibraryRecord(libraryRecord openLibraryBook) book {
	b := book{
		Title: libraryRecord.Work.Title,
	}

	for _, edition := range libraryRecord.Editions.Entries {
		if edition.Subtitle != "" {
			b.Subtitle = edition.Subtitle
			break
		}
	}

	for _, edition := range libraryRecord.Editions.Entries {
		if edition.Isbn13[0] != "" {
			b.ISBN = edition.Isbn13[0]
			break
		}
	}

	return b
}

func getBookDetails(id string, httpClient http.Client) (openLibraryBook, error) {
	libraryRecord := openLibraryBook{}
	w, err := getWorkById(id, httpClient)
	if err != nil {
		return openLibraryBook{}, err
	}
	libraryRecord.Work = w

	e, err := getWorkEditions(id, httpClient)
	if err != nil {
		return openLibraryBook{}, err
	}
	libraryRecord.Editions = e

	return libraryRecord, nil
}

func getWorkById(id string, httpClient http.Client) (work, error) {
	url := baseURL + "/works/" + id + ".json"

	resp, err := httpClient.Get(url)
	if err != nil {
		return work{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return work{}, fmt.Errorf("received a %d reponse from the api\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return work{}, err
	}

	w := work{}
	if err := json.Unmarshal(body, &w); err != nil {
		return work{}, nil
	}
	return w, nil
}

func getWorkEditions(id string, httpClient http.Client) (editions, error) {
	url := baseURL + "/works/" + id + "/editions.json"
	resp, err := httpClient.Get(url)
	if err != nil {
		return editions{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return editions{}, fmt.Errorf("received a %d reponse from the api\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return editions{}, err
	}

	e := editions{}
	if err := json.Unmarshal(body, &e); err != nil {
		return editions{}, err
	}

	// Only return english editions
	e2 := editions{
		Size:    0,
		Entries: []edition{},
	}
	for _, edition := range e.Entries {
		if edition.Languages[0].Key == "/languages/eng" {
			e2.Entries = append(e2.Entries, edition)
			e2.Size++
		}
	}

	return e2, nil
}
