package openlibraryapi

type results struct {
	Start    int               `json:"start"`
	Offset   int               `json:"offset"`
	NumFound int               `json:"num_found"`
	Books    []bookSolrResults `json:"docs"`
}

type bookSolrResults struct {
	Title            string   `json:"title"`
	AuthorName       []string `json:"author_name"`
	FirstPublishYear int      `json:"first_publish_year"`
	Key              string   `json:"key"`
	AuthorKey        []string `json:"author_key"`
}
