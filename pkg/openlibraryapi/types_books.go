package openlibraryapi

// type author struct {
// 	Name string `json:"personal_name"`
// }

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
	Description struct {
		Value string `json:"value"`
	} `json:"description"`
	Key    string `json:"key"`
	Covers []int  `json:"covers"`
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
