package repository

import (
	"fmt"
	"net/url"
	"strings"
)

// Repo represents a source of issues.
type Repo struct {
	Owner      string
	Name       string
	FullName   string
	IssuesLink string
}

// ParseGithubRepo takes a string representation of a repository, and attempts to build a Repo.
// Acceptable formats are:
//  "owner/name"
//  "github.com/owner/name"
//  "https://github.com/owner/name"
func ParseGithubRepo(repoStr string) (Repo, error) {
	var repo Repo
	u, err := url.Parse(repoStr)
	if err != nil {
		return repo, err
	}
	var parts []string
	if u.Host == "github.com" {
		// Assume we were given a valid Github URL
		path := strings.Trim(u.Path, "/")
		parts = strings.Split(path, "/")
	} else if u.Scheme == "" {
		// Assume we have the format "owner/name" or "github.com/owner/name"
		ghI := strings.Index(u.Path, "github.com/")
		if ghI < 0 {
			ghI = 0
		} else {
			ghI += 11
		}
		parts = strings.Split(repoStr[ghI:], "/")
	}
	if len(parts) < 2 {
		return repo, fmt.Errorf("%s not a valid repo", repoStr)
	}
	repo.Owner = parts[0]
	repo.Name = parts[1]
	repo.FullName = fmt.Sprintf("%s/%s", parts[0], parts[1])

	repo.IssuesLink = fmt.Sprintf("https://github.com/%s/issues?q=is%%3Aopen+is%%3Aissue+no%%3Alabel", repo.FullName)
	return repo, nil
}
