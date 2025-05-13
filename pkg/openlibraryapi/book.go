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
	ISBN     []string
	Cover    string
	// categories []string
}

func (c *Client) GetBookById(id string) (book, error) {
	id = strings.ToUpper(id)

	lr, err := getBookDetails(id, c.httpClient)
	if err != nil {
		return book{}, err
	}

	b := book{
		Title: lr.Work.Title,
	}

	for _, edition := range lr.Editions.Entries {
		if edition.Subtitle != "" {
			b.Subtitle = edition.Subtitle
			break
		}
	}

	return b, nil

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

	// fmt.Printf("t:\n%v\n", libraryRecord)

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
