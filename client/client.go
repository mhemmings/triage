package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mhemmings/triage/issues"
)

// IssueFilters contains the possible filters to use when searching for issues.
type IssueFilters struct {
	// By default, only unlabeled issues are shown. A slice of label names can be provided to override this
	// and match only those label names.
	Labels Labels
	// If true, all issues will be gathered regardless of label name.
	ShowAll bool
	// Filter by time since created.
	Since Since
}

// Labels is a slice of strings that implements the gnuflag.Value interface.
type Labels []string

func (l Labels) String() string {
	return strings.Join(l, ",")
}

func (l *Labels) Set(str string) error {
	*l = strings.Split(str, ",")
	return nil
}

// Since is a time.Duration that implements the gnuflag.Value interface.
type Since time.Duration

func (s Since) String() string {
	dur := time.Duration(s)
	if dur == 0 {
		return ""
	}
	return dur.String()
}

// Set takes a string and assigns the value to the given Since. It supports
// the time formats from time.ParseDuration, as well as 'd' for day, and
// 'w' for week.
func (s *Since) Set(str string) error {
	var v float64
	var d time.Duration
	var err error

	switch {
	case strings.HasSuffix(str, "d"):
		_, err = fmt.Sscanf(str, "%fd", &v)
		d = time.Duration(v * float64(24*time.Hour))
	case strings.HasSuffix(str, "w"):
		_, err = fmt.Sscanf(str, "%fw", &v)
		d = time.Duration(v * float64(24*time.Hour*7))
	default:
		d, err = time.ParseDuration(str)
	}

	*s = Since(d)

	return err
}

// AsTime returns the current time.Time minus the Since duration.
func (s Since) AsTime() time.Time {
	return time.Now().Add(-time.Duration(s))
}

// Client is the interface that groups methods for interacring with an issue provider
type Client interface {
	// GetIssuesForTriage gets a list of "untriaged" issues from the Client for the given
	// repository name and owner, filtered by the provided IssueFilters.
	GetIssuesForTriage(ctx context.Context, owner string, repo string, filters IssueFilters) ([]issues.Issue, error)
}
