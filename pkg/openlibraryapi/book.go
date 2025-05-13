package openlibraryapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	url := baseURL + "/works/" + id + ".json"

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return book{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return book{}, fmt.Errorf("received a %d reponse from the api\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return book{}, err
	}

	work := work{}
	if err := json.Unmarshal(body, &work); err != nil {
		return book{}, nil
	}

	b := book{
		Title: work.Title,
		// Subtitle: source.Data.SubTitle,
		// Cover:    source.Data.Cover.Large,
	}

	return b, nil

}

func prettyPrint(data any) {
	encodedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal("failed to pretty print map")
	}
	fmt.Println(string(encodedData))
}
