package triage

import (
	"context"
)

// Client is the interface that groups methods for interacring with an issue provider
type Client interface {
	// GetIssuesForTriage gets a list of "untriaged" issues from the Client for the given
	// repository name and owner
	GetIssuesForTriage(ctx context.Context, owner string, repo string) ([]Issue, error)
}
