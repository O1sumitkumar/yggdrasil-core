package internet

import "context"

// Result is one search hit. Providers should return structured rows, not a text blob.
type Result struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Snippet     string `json:"snippet"`
	Source      string `json:"source,omitempty"`
	PublishedAt string `json:"publishedAt,omitempty"`
}

// SearchProvider can be replaced without changing tool registration.
type SearchProvider interface {
	Search(ctx context.Context, query string) ([]Result, error)
}

// Page is readable text extracted from a URL.
type Page struct {
	Title    string         `json:"title"`
	URL      string         `json:"url"`
	Content  string         `json:"content"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// PageFetcher opens a URL and returns text.
type PageFetcher interface {
	Open(ctx context.Context, rawURL string) (Page, error)
}
