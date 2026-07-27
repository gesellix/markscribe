package main

import (
	"testing"
	"time"

	"github.com/shurcooL/githubv4"
)

func TestGistFromQL(t *testing.T) {
	createdAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	in := qlGist{
		Name:        "my-gist",
		Description: "a gist",
		URL:         "https://gist.github.com/my-gist",
		CreatedAt:   githubv4.DateTime{Time: createdAt},
	}

	got := gistFromQL(in)

	want := Gist{
		Name:        "my-gist",
		Description: "a gist",
		URL:         "https://gist.github.com/my-gist",
		CreatedAt:   createdAt,
	}
	if got != want {
		t.Errorf("gistFromQL() = %+v, want %+v", got, want)
	}
}

func TestPullRequestFromQL(t *testing.T) {
	createdAt := time.Date(2023, 6, 1, 12, 30, 0, 0, time.UTC)
	in := qlPullRequest{
		URL:       "https://github.com/o/r/pull/1",
		Title:     "Fix bug",
		State:     githubv4.PullRequestStateMerged,
		CreatedAt: githubv4.DateTime{Time: createdAt},
		Repository: qlRepository{
			NameWithOwner: "o/r",
			URL:           "https://github.com/o/r",
			Description:   "repo description",
			IsPrivate:     false,
			Stargazers: struct {
				TotalCount githubv4.Int
			}{TotalCount: 42},
		},
	}

	got := pullRequestFromQL(in)

	want := PullRequest{
		Title:     "Fix bug",
		URL:       "https://github.com/o/r/pull/1",
		State:     "MERGED",
		CreatedAt: createdAt,
		Repo: Repo{
			Name:        "o/r",
			URL:         "https://github.com/o/r",
			Description: "repo description",
			IsPrivate:   false,
			Stargazers:  42,
		},
	}
	if got != want {
		t.Errorf("pullRequestFromQL() = %+v, want %+v", got, want)
	}
}

func TestReleaseFromQL(t *testing.T) {
	publishedAt := time.Date(2022, 3, 10, 0, 0, 0, 0, time.UTC)
	in := qlRelease{
		Nodes: []struct {
			Name         githubv4.String
			TagName      githubv4.String
			PublishedAt  githubv4.DateTime
			URL          githubv4.String
			IsPrerelease githubv4.Boolean
			IsDraft      githubv4.Boolean
		}{
			{
				Name:        "v1.0.0",
				TagName:     "v1.0.0",
				PublishedAt: githubv4.DateTime{Time: publishedAt},
				URL:         "https://github.com/o/r/releases/v1.0.0",
			},
			{
				Name:    "v0.9.0",
				TagName: "v0.9.0",
			},
		},
	}

	got := releaseFromQL(in)

	want := Release{
		Name:        "v1.0.0",
		TagName:     "v1.0.0",
		PublishedAt: publishedAt,
		URL:         "https://github.com/o/r/releases/v1.0.0",
	}
	if got != want {
		t.Errorf("releaseFromQL() = %+v, want %+v", got, want)
	}
}

func TestRepoFromQL(t *testing.T) {
	in := qlRepository{
		NameWithOwner: "o/r",
		URL:           "https://github.com/o/r",
		Description:   "repo description",
		IsPrivate:     true,
		Stargazers: struct {
			TotalCount githubv4.Int
		}{TotalCount: 7},
	}

	got := repoFromQL(in)

	want := Repo{
		Name:        "o/r",
		URL:         "https://github.com/o/r",
		Description: "repo description",
		IsPrivate:   true,
		Stargazers:  7,
		// LastRelease is never populated by repoFromQL, stays zero-value.
	}
	if got != want {
		t.Errorf("repoFromQL() = %+v, want %+v", got, want)
	}
}

func TestUserFromQL(t *testing.T) {
	in := qlUser{
		Login:     "octocat",
		Name:      "The Octocat",
		AvatarURL: "https://github.com/octocat.png",
		URL:       "https://github.com/octocat",
	}

	got := userFromQL(in)

	want := User{
		Login:     "octocat",
		Name:      "The Octocat",
		AvatarURL: "https://github.com/octocat.png",
		URL:       "https://github.com/octocat",
	}
	if got != want {
		t.Errorf("userFromQL() = %+v, want %+v", got, want)
	}
}
