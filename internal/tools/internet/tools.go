package internet

import (
	"context"
	"fmt"
	"strings"
)

// SearchTool is internet.search.
type SearchTool struct {
	Provider SearchProvider
}

func NewSearch(provider SearchProvider) *SearchTool {
	if provider == nil {
		provider = DuckDuckGo{}
	}
	return &SearchTool{Provider: provider}
}

func (t *SearchTool) ID() string          { return "internet.search" }
func (t *SearchTool) DisplayName() string { return "Web Search" }
func (t *SearchTool) Description() string { return "Search the public internet" }

func (t *SearchTool) Execute(ctx context.Context, args map[string]any) (map[string]any, error) {
	query, _ := args["query"].(string)
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("query required")
	}
	results, err := t.Provider.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	return map[string]any{"results": results}, nil
}

// OpenTool is internet.open.
type OpenTool struct {
	Fetcher PageFetcher
}

func NewOpen(fetcher PageFetcher) *OpenTool {
	if fetcher == nil {
		fetcher = HTTPFetcher{}
	}
	return &OpenTool{Fetcher: fetcher}
}

func (t *OpenTool) ID() string          { return "internet.open" }
func (t *OpenTool) DisplayName() string { return "Open Web Page" }
func (t *OpenTool) Description() string { return "Open a web page and return readable text" }

func (t *OpenTool) Execute(ctx context.Context, args map[string]any) (map[string]any, error) {
	rawURL, _ := args["url"].(string)
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return nil, fmt.Errorf("url required")
	}
	page, err := t.Fetcher.Open(ctx, rawURL)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"title":    page.Title,
		"url":      page.URL,
		"content":  page.Content,
		"metadata": page.Metadata,
	}, nil
}
