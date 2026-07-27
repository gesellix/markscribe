package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func rssFixture(items int) string {
	body := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
<channel>
<title>Test Feed</title>
<link>http://example.com</link>
<description>desc</description>
`
	pubDates := []string{
		"Mon, 01 Jan 2024 10:00:00 GMT",
		"Tue, 02 Jan 2024 10:00:00 GMT",
		"Wed, 03 Jan 2024 10:00:00 GMT",
		"Thu, 04 Jan 2024 10:00:00 GMT",
		"Fri, 05 Jan 2024 10:00:00 GMT",
	}
	for i := 0; i < items; i++ {
		body += fmt.Sprintf(`<item>
<title>Item %d</title>
<link>http://example.com/%d</link>
<pubDate>%s</pubDate>
</item>
`, i+1, i+1, pubDates[i])
	}
	body += `</channel>
</rss>`
	return body
}

func newRSSServer(t *testing.T, items int) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml; charset=UTF-8")
		if _, err := w.Write([]byte(rssFixture(items))); err != nil {
			t.Fatal(err)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func TestRSSFeed(t *testing.T) {
	t.Run("maps entries in feed order", func(t *testing.T) {
		server := newRSSServer(t, 3)

		got := rssFeed(server.URL, 3)

		if len(got) != 3 {
			t.Fatalf("len(rssFeed()) = %d, want 3", len(got))
		}

		wantTitles := []string{"Item 1", "Item 2", "Item 3"}
		wantURLs := []string{"http://example.com/1", "http://example.com/2", "http://example.com/3"}
		for i, entry := range got {
			if entry.Title != wantTitles[i] {
				t.Errorf("entry[%d].Title = %q, want %q", i, entry.Title, wantTitles[i])
			}
			if entry.URL != wantURLs[i] {
				t.Errorf("entry[%d].URL = %q, want %q", i, entry.URL, wantURLs[i])
			}
			if entry.PublishedAt.IsZero() {
				t.Errorf("entry[%d].PublishedAt is zero, want a parsed date", i)
			}
		}

		wantFirstDate := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		if !got[0].PublishedAt.Equal(wantFirstDate) {
			t.Errorf("entry[0].PublishedAt = %v, want %v", got[0].PublishedAt, wantFirstDate)
		}
	})

	t.Run("truncates to count", func(t *testing.T) {
		server := newRSSServer(t, 5)

		got := rssFeed(server.URL, 2)

		if len(got) != 2 {
			t.Fatalf("len(rssFeed()) = %d, want 2", len(got))
		}
	})

	t.Run("count larger than available items returns all items", func(t *testing.T) {
		server := newRSSServer(t, 3)

		got := rssFeed(server.URL, 10)

		if len(got) != 3 {
			t.Fatalf("len(rssFeed()) = %d, want 3", len(got))
		}
	})
}
