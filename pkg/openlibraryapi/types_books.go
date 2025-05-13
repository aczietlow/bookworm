package openlibraryapi

type bookResponse struct {
	// will match JSON key "/books/<id>"
	Records map[string]bookRecord `json:"records"`
}

type bookRecord struct {
	Data bookData `json:"data"`
}

type bookData struct {
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

type work struct {
	Type struct {
		Key string `json:"key"`
	} `json:"type"`
	Title    string   `json:"title"`
	Subjects []string `json:"subjects"`
	// Authors  []struct {
	// 	Type struct {
	// 		Key string `json:"key"`
	// 	} `json:"type"`
	// 	Author struct {
	// 		Key string `json:"key"`
	// 	} `json:"author"`
	// } `json:"authors"`
	Key string `json:"key"`
	// Covers         []int  `json:"covers"`
}
