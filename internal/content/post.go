package content

import "time"

type Post struct {
	Title           string
	Slug            string
	Date            time.Time
	Updated         time.Time
	Summary         string
	Description     string
	Tags            []string
	Draft           bool
	Cover           string
	CanonicalURL    string
	ReadingDuration int
	ContentHTML     string
	SourcePath      string
	AssetBase       string
}
