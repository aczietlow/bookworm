package openlibraryapi

type author struct {
	PersonalName   string   `json:"personal_name"`
	Key            string   `json:"key"`
	EntityType     string   `json:"entity_type"`
	BirthDate      string   `json:"birth_date"`
	Links          []link   `json:"links"`
	AlternateNames []string `json:"alternate_names"`
	Name           string   `json:"name"`
	Title          string   `json:"title"`
	Bio            string   `json:"bio"`
	FullerName     string   `json:"fuller_name"`
	SourceRecords  []string `json:"source_records"`
	Photos         []int    `json:"photos"`
}

type link struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}
