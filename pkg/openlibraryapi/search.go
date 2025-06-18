package openlibraryapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) SearchQuery(query string) ([]bookSolrResults, error) {
	url := baseURL + "/search.json?q=" + query
	url += "&limit=5"

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return []bookSolrResults{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []bookSolrResults{}, fmt.Errorf("Received a %d response from the api", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []bookSolrResults{}, err
	}

	solrResults := results{}
	if err = json.Unmarshal(body, &solrResults); err != nil {
		return []bookSolrResults{}, err
	}

	return solrResults.Books, nil
}
