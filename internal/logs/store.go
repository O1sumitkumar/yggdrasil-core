package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Entry describes a log file available for diagnostics.
type Entry struct {
	Name      string    `json:"name"`
	Kind      string    `json:"kind"` // daemon | runtime | other
	SizeBytes int64     `json:"size_bytes"`
	ModTime   time.Time `json:"modified_at"`
	Label     string    `json:"label"`
}

// Content is a log file payload returned to the UI.
type Content struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Label     string `json:"label"`
	Content   string `json:"content"`
	Truncated bool   `json:"truncated"`
	SizeBytes int64  `json:"size_bytes"`
}

// Store reads log files from the application logs directory.
type Store struct {
	Dir string
}

// List returns safe-to-show log files, newest first.
func (s *Store) List() ([]Entry, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, err
	}
	var out []Entry
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !isAllowedLogName(name) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		out = append(out, Entry{
			Name:      name,
			Kind:      classify(name),
			SizeBytes: info.Size(),
			ModTime:   info.ModTime().UTC(),
			Label:     labelFor(name),
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ModTime.After(out[j].ModTime)
	})
	return out, nil
}

// ReadTail returns the last maxBytes of a log file.
func (s *Store) ReadTail(name string, maxBytes int64) (Content, error) {
	if !isAllowedLogName(name) {
		return Content{}, fmt.Errorf("log not found")
	}
	path, err := s.resolve(name)
	if err != nil {
		return Content{}, err
	}
	st, err := os.Stat(path)
	if err != nil {
		return Content{}, err
	}
	if maxBytes <= 0 {
		maxBytes = 256 * 1024
	}
	if maxBytes > 2<<20 {
		maxBytes = 2 << 20
	}

	f, err := os.Open(path)
	if err != nil {
		return Content{}, err
	}
	defer f.Close()

	size := st.Size()
	truncated := false
	offset := int64(0)
	readSize := size
	if size > maxBytes {
		offset = size - maxBytes
		readSize = maxBytes
		truncated = true
	}
	if _, err := f.Seek(offset, 0); err != nil {
		return Content{}, err
	}
	buf := make([]byte, readSize)
	n, err := f.Read(buf)
	if err != nil && n == 0 {
		return Content{}, err
	}
	text := string(buf[:n])
	if truncated {
		// Avoid starting mid-line when possible.
		if i := strings.IndexByte(text, '\n'); i >= 0 && i+1 < len(text) {
			text = text[i+1:]
		}
	}
	return Content{
		Name:      name,
		Kind:      classify(name),
		Label:     labelFor(name),
		Content:   text,
		Truncated: truncated,
		SizeBytes: size,
	}, nil
}

func (s *Store) resolve(name string) (string, error) {
	base := filepath.Clean(s.Dir)
	target := filepath.Clean(filepath.Join(base, name))
	if !strings.HasPrefix(target, base+string(os.PathSeparator)) && target != base {
		return "", fmt.Errorf("invalid log path")
	}
	return target, nil
}

func isAllowedLogName(name string) bool {
	if name == "" || strings.Contains(name, "..") || strings.ContainsAny(name, `/\`) {
		return false
	}
	lower := strings.ToLower(name)
	if strings.Contains(lower, "secret") || strings.Contains(lower, "key") || strings.HasSuffix(lower, ".pem") {
		return false
	}
	return strings.HasSuffix(lower, ".log") || strings.HasSuffix(lower, ".txt")
}

func classify(name string) string {
	lower := strings.ToLower(name)
	switch {
	case strings.HasPrefix(lower, "daemon"):
		return "daemon"
	case strings.HasPrefix(lower, "llamacpp") || strings.Contains(lower, "runtime"):
		return "runtime"
	default:
		return "other"
	}
}

func labelFor(name string) string {
	lower := strings.ToLower(name)
	switch {
	case strings.HasPrefix(lower, "daemon"):
		return "Application log"
	case strings.HasPrefix(lower, "llamacpp"):
		return "Model runtime log"
	default:
		return name
	}
}
