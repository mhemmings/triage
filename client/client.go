package client

import (
	"context"

	"github.com/mhemmings/triage/issues"
)

// Client is the interface that groups methods for interacring with an issue provider
type Client interface {
	// GetIssuesForTriage gets a list of "untriaged" issues from the Client for the given
	// repository name and owner. The labels slice filters by issue label. If showAll is
	// true, then all issues will be shown regardless of label.
	GetIssuesForTriage(ctx context.Context, owner string, repo string, labels []string, showAll bool) ([]issues.Issue, error)
}
