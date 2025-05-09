package openlibraryapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type results struct {
	Start    int     `json:"start"`
	Offset   int     `json:"offset"`
	NumFound int     `json:"num_found"`
	Books    []Books `json:"docs"`
}

type Books struct {
	Title            string   `json:"title"`
	AuthorName       []string `json:"author_name"`
	FirstPublishYear int      `json:"first_publish_year"`
	Key              string   `json:"key"`
	AuthorKey        []string `json:"author_key"`
}

func (c *Client) SearchQuery(query string) ([]Books, error) {
	url := baseURL + "/search.json?q=" + query
	url += "&limit=2"
	fmt.Printf("debug url %s\n", url)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return []Books{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []Books{}, fmt.Errorf("Received a %d response from the api", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Books{}, err
	}

	solrResults := results{}
	if err = json.Unmarshal(body, &solrResults); err != nil {
		return []Books{}, err
	}

	return solrResults.Books, nil
}
