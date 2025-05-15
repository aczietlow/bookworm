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
	Name string `json:"personal_name"`
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

type openLibraryBook struct {
	Work     work
	Editions editions
}

type work struct {
	Type struct {
		Key string `json:"key"`
	} `json:"type"`
	Title       string   `json:"title"`
	Subjects    []string `json:"subjects"`
	Description string   `json:"description"`
	Key         string   `json:"key"`
	Covers      []int    `json:"covers"`
}

type editions struct {
	Size    int       `json:"size"`
	Entries []edition `json:"entries"`
}

type edition struct {
	Type struct {
		Key string `json:"key"`
	} `json:"type"`
	AuthorKeys []struct {
		Key string `json:"key"`
	} `json:"authors"`
	Authors   []string
	Languages []struct {
		Key string `json:"key"`
	} `json:"languages"`
	PublishDate string   `json:"publish_date"`
	Publishers  []string `json:"publishers"`
	Subjects    []string `json:"subjects,omitempty"`
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle"`
	FullTitle   string   `json:"full_title,omitempty"`
	Key         string   `json:"key"`
	Covers      []int    `json:"covers,omitempty"`
	Isbn13      []string `json:"isbn_13"`
}
