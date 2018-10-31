package triage

import (
	"context"
	"net/http"

	"github.com/google/go-github/v18/github"
	"golang.org/x/oauth2"
)

// GithubClient implements Client and interacts with the Github API
type GithubClient struct {
	client *github.Client
}

// NewGithubClient returns Guthub client created from the provided auth
// token.
// If no token is provided, the client with interact with the Guthub API
// without authentication. Doing this imposes much stricter API rate limits.
// See https://developer.github.com/v3/#rate-limiting.
func NewGithubClient(token string) Client {
	var httpClient *http.Client

	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		httpClient = oauth2.NewClient(context.Background(), ts)
	}

	return &GithubClient{
		client: github.NewClient(httpClient),
	}
}

// GetIssuesForTriage returns a list of untriaged issues for the given repository, or an error.
// We treat any issue without a label as being untriaged.
func (gh GithubClient) GetIssuesForTriage(ctx context.Context, owner string, repo string) ([]Issue, error) {
	var allIssues []Issue

	// Currenly no way to get just issues without a label, so we get them all and filter later.
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		issues, resp, err := gh.client.Issues.ListByRepo(ctx, owner, repo, opt)
		if err != nil {
			return allIssues, err
		}

		for _, issue := range issues {
			if issue.PullRequestLinks != nil || len(issue.Labels) > 0 {
				// If the issue is actually a PR, or if it has labels, we're not interested.
				continue
			}

			allIssues = append(allIssues, Issue{
				Title:    *issue.Title,
				Link:     *issue.HTMLURL,
				User:     *issue.User.Login,
				Comments: *issue.Comments,
				Created:  *issue.CreatedAt,
			})
		}

		if resp.NextPage == 0 {
			// Last page.
			break
		}

		opt.Page = resp.NextPage
	}

	return allIssues, nil
}
