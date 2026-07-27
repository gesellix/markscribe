package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
)

func TestHumanized(t *testing.T) {
	t.Run("within last 24h returns today", func(t *testing.T) {
		got := humanized(time.Now())
		if got != "today" {
			t.Errorf("humanized(now) = %q, want %q", got, "today")
		}
	})

	t.Run("older than 24h matches humanize.Time", func(t *testing.T) {
		old := time.Now().AddDate(0, 0, -10)
		truncated := time.Date(old.Year(), old.Month(), old.Day(), 0, 0, 0, 0, old.Location())
		want := humanize.Time(truncated)

		got := humanized(old)
		if got != want {
			t.Errorf("humanized(old) = %q, want %q", got, want)
		}
	})

	t.Run("non-time value falls back to Sprintf", func(t *testing.T) {
		cases := []interface{}{42, "hello", 3.14}
		for _, c := range cases {
			got := humanized(c)
			want := fmt.Sprintf("%v", c)
			if got != want {
				t.Errorf("humanized(%v) = %q, want %q", c, got, want)
			}
		}
	})
}

func TestReverse(t *testing.T) {
	t.Run("reverses []int", func(t *testing.T) {
		in := []int{1, 2, 3, 4, 5}
		got := reverse(in).([]int)
		want := []int{5, 4, 3, 2, 1}
		if !equalInts(got, want) {
			t.Errorf("reverse() = %v, want %v", got, want)
		}
	})

	t.Run("reverses []string", func(t *testing.T) {
		in := []string{"a", "b", "c"}
		got := reverse(in).([]string)
		want := []string{"c", "b", "a"}
		if !equalStrings(got, want) {
			t.Errorf("reverse() = %v, want %v", got, want)
		}
	})

	t.Run("reverses []Repo", func(t *testing.T) {
		in := []Repo{{Name: "a"}, {Name: "b"}, {Name: "c"}}
		got := reverse(in).([]Repo)
		want := []Repo{{Name: "c"}, {Name: "b"}, {Name: "a"}}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("reverse() = %+v, want %+v", got, want)
				break
			}
		}
	})

	t.Run("single element unchanged", func(t *testing.T) {
		in := []int{1}
		got := reverse(in).([]int)
		if len(got) != 1 || got[0] != 1 {
			t.Errorf("reverse() = %v, want %v", got, in)
		}
	})

	t.Run("empty slice unchanged", func(t *testing.T) {
		in := []int{}
		got := reverse(in).([]int)
		if len(got) != 0 {
			t.Errorf("reverse() = %v, want empty", got)
		}
	})
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
