package internet

import (
	"strings"
	"testing"
)

func TestParseDuckHTML(t *testing.T) {
	page := `<html><body>
<a class="result__a" href="https://duckduckgo.com/l/?uddg=https%3A%2F%2Fweather.example%2Fjuneau">Juneau Weather</a>
<a class="result__snippet">Current conditions in Juneau.</a>
</body></html>`
	results := ParseDuckHTML(page)
	if len(results) != 1 {
		t.Fatalf("results %#v", results)
	}
	if results[0].URL != "https://weather.example/juneau" || !strings.Contains(results[0].Snippet, "Current") {
		t.Fatalf("result %#v", results[0])
	}
}

func TestExtractTextSkipsNavigation(t *testing.T) {
	page := `<html><head><title>Docs</title></head><body><nav>Home</nav><main><p>Install the package.</p></main><script>junk()</script></body></html>`
	title, content := extractText(page)
	if title != "Docs" || !strings.Contains(content, "Install the package") || strings.Contains(content, "junk") || strings.Contains(content, "Home") {
		t.Fatalf("title=%q content=%q", title, content)
	}
}
