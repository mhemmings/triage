package triage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v18/github"
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

// GetIssuesForTriage returns a list of untriaged issues for the given repository, or an error.
// If labels is empty, only issues with no labels will be shown.
// If showAll is true, all issues will be returned, regardless of label.
func (gh GithubClient) GetIssuesForTriage(ctx context.Context, owner string, repo string, labels []string, showAll bool) ([]Issue, error) {
	var allIssues []Issue

	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	if !showAll {
		opt.Labels = labels
	}

	for {
		issues, resp, err := gh.client.Issues.ListByRepo(ctx, owner, repo, opt)
		if err != nil {
			return allIssues, err
		}

		for _, issue := range issues {
			if issue.PullRequestLinks != nil {
				// If the issue is actually a PR, skip.
				continue
			}

			if !showAll && len(labels) == 0 && len(issue.Labels) > 0 {
				// Currenly no way to get just issues without a label, so we get them all and filter now.
				continue
			}

			iss := Issue{
				Title:    *issue.Title,
				Link:     *issue.HTMLURL,
				User:     *issue.User.Login,
				Comments: *issue.Comments,
				Created:  *issue.CreatedAt,
			}
			for _, label := range issue.Labels {
				iss.Labels = append(iss.Labels, Label{
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
