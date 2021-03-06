package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v18/github"
	"github.com/mhemmings/triage/issues"
	"golang.org/x/oauth2"
)

// GithubClient implements Client and interacts with the Github API
type GithubClient struct {
	client *github.Client
}

// NewGithubClient returns Github client created from the provided auth
// token.
// If no token is provided, the client with interact with the Github API
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

// GetIssuesForTriage returns a list of untriaged issues for the given repository filtered by the provided IssueFilters, or an error.
func (gh GithubClient) GetIssuesForTriage(ctx context.Context, owner string, repo string, filters IssueFilters) ([]issues.Issue, error) {
	var allIssues []issues.Issue

	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	if !filters.Since.AsTime().IsZero() {
		opt.Since = filters.Since.AsTime()
	}

	if !filters.ShowAll {
		opt.Labels = filters.Labels
	}

	for {
		ghIssues, resp, err := gh.client.Issues.ListByRepo(ctx, owner, repo, opt)
		if err != nil {
			return allIssues, err
		}

		for _, issue := range ghIssues {
			if issue.PullRequestLinks != nil {
				// If the issue is actually a PR, skip.
				continue
			}

			if !filters.ShowAll && len(filters.Labels) == 0 && len(issue.Labels) > 0 {
				// Currenly no way to get just issues without a label, so we get them all and filter now.
				continue
			}

			iss := issues.Issue{
				Title:    *issue.Title,
				Link:     *issue.HTMLURL,
				User:     *issue.User.Login,
				Comments: *issue.Comments,
				Created:  *issue.CreatedAt,
			}
			for _, label := range issue.Labels {
				iss.Labels = append(iss.Labels, issues.Label{
					Name:   *label.Name,
					Colour: fmt.Sprintf("#%s", *label.Color),
				})
			}
			allIssues = append(allIssues, iss)
		}

		if resp.NextPage == 0 {
			// Last page.
			break
		}

		opt.Page = resp.NextPage
	}

	return allIssues, nil
}
