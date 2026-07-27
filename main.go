// Package main implements markscribe, a template-driven markdown
// generator with GitHub, RSS, and reading-list data sources.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/KyleBanks/goodreads"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// debugTransport logs the outgoing GraphQL request and the full response
// (status, headers, body) for any request whose response body contains a
// GraphQL "errors" field. Temporary debugging aid, not for production use.
type debugTransport struct {
	base http.RoundTripper
}

func (d *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(reqBody))
	}

	resp, err := d.base.RoundTrip(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "--- TRANSPORT ERROR for %s %s: %v\n", req.Method, req.URL, err)
		return resp, err
	}

	var respBody []byte
	if resp.Body != nil {
		respBody, _ = io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
	}

	if bytes.Contains(respBody, []byte(`"errors"`)) {
		fmt.Fprintf(os.Stderr, "\n=== FAILING REQUEST %s %s ===\n", req.Method, req.URL)
		fmt.Fprintf(os.Stderr, "--- request headers ---\n")
		for k, v := range req.Header {
			if strings.EqualFold(k, "Authorization") {
				fmt.Fprintf(os.Stderr, "%s: [redacted]\n", k)
				continue
			}
			fmt.Fprintf(os.Stderr, "%s: %s\n", k, strings.Join(v, ", "))
		}
		fmt.Fprintf(os.Stderr, "--- request body ---\n%s\n", string(reqBody))
		fmt.Fprintf(os.Stderr, "--- response status: %s ---\n", resp.Status)
		fmt.Fprintf(os.Stderr, "--- response headers ---\n")
		for k, v := range resp.Header {
			fmt.Fprintf(os.Stderr, "%s: %s\n", k, strings.Join(v, ", "))
		}
		fmt.Fprintf(os.Stderr, "--- response body ---\n%s\n=== END ===\n\n", string(respBody))
	}

	return resp, err
}

var (
	gitHubClient    *githubv4.Client
	goodReadsClient *goodreads.Client
	goodReadsID     string
	username        string

	write = flag.String("write", "", "write output to")
)

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: markscribe [template]")
		os.Exit(1)
	}

	tplIn, err := os.ReadFile(flag.Args()[0])
	if err != nil {
		fmt.Println("Can't read file:", err)
		os.Exit(1)
	}

	tpl, err := template.New("tpl").Funcs(template.FuncMap{
		/* GitHub */
		"recentContributions": recentContributions,
		"recentPullRequests":  recentPullRequests,
		"recentRepos":         recentRepos,
		"recentForks":         recentForks,
		"recentReleases":      recentReleases,
		"followers":           recentFollowers,
		"recentStars":         recentStars,
		"gists":               gists,
		"sponsors":            sponsors,
		"repo":                repo,
		/* RSS */
		"rss": rssFeed,
		/* GoodReads */
		"goodReadsReviews":          goodReadsReviews,
		"goodReadsCurrentlyReading": goodReadsCurrentlyReading,
		/* Literal.club */
		"literalClubCurrentlyReading": literalClubCurrentlyReading,
		/* Utils */
		"humanize": humanized,
		"reverse":  reverse,
		"now":      time.Now,
		"contains": strings.Contains,
		"toLower":  strings.ToLower,
	}).Parse(string(tplIn))
	if err != nil {
		fmt.Println("Can't parse template:", err)
		os.Exit(1)
	}

	var httpClient *http.Client
	gitHubToken := os.Getenv("GITHUB_TOKEN")
	goodReadsToken := os.Getenv("GOODREADS_TOKEN")
	goodReadsID = os.Getenv("GOODREADS_USER_ID")
	if len(gitHubToken) > 0 {
		httpClient = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: gitHubToken},
		))
		httpClient.Transport = &debugTransport{base: httpClient.Transport}
	}

	gitHubClient = githubv4.NewClient(httpClient)
	goodReadsClient = goodreads.NewClient(goodReadsToken)

	if len(gitHubToken) > 0 {
		username, err = getUsername()
		if err != nil {
			fmt.Println("Can't retrieve GitHub profile:", err)
			os.Exit(1)
		}
	}

	w := os.Stdout
	if len(*write) > 0 {
		f, err := os.Create(*write)
		if err != nil {
			fmt.Println("Can't create:", err)
			os.Exit(1)
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
		w = f
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Can't render template:", err)
		os.Exit(1)
	}
}
