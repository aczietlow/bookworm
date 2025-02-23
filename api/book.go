package api

type BookResponse struct {
	// will match JSON key "/books/<id>"
	Records map[string]Book `json:"records"`
}

type Book struct {
	Data Data `json:"data"`
}

type Data struct {
	Key       string    `json:"key"`
	Author    []author  `json:"authors"`
	Title     string    `json:"title"`
	SubTitle  string    `json:"subtitle"`
	Identifer identifer `json:"identifiers"`
	Cover     cover     `json:"cover"`
}

type author struct {
	Name string `json:"name"`
}

type identifer struct {
	ISBN10      []string `json:"isbn_10"`
	ISBN13      []string `json:"isbn_13"`
	OpenLibrary []string `json:"openlibrary"`
	// Goodreads   []string `json:"goodreads"`
}

type cover struct {
	Large string `json:"large"`
}
