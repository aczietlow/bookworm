package openlibraryapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) DebugQuery() (any, error) {
	authorOID := "OL2830895A"
	// https://openlibrary.org/authors/OL23919A.json
	url := baseURL + "/authors/" + authorOID + ".json"

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var authorData any
	if err = json.Unmarshal(body, &authorData); err != nil {
		return nil, err
	}

	return authorData, nil
}

func (c *Client) DebugQueryJson() (any, error) {
	authorOID := "OL2830895A"
	// https://openlibrary.org/authors/OL23919A.json
	url := baseURL + "/authors/" + authorOID + ".json"

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var authorData any
	if err = json.Unmarshal(body, &authorData); err != nil {
		return nil, err
	}

	return authorData, nil
}
