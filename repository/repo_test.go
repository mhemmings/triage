package repository

import (
	"errors"
	"testing"
)

var githubTests = []struct {
	Input    string
	Expected Repo
	Error    error
}{{
	Input: "foo/bar",
	Expected: Repo{
		Owner:      "foo",
		Name:       "bar",
		FullName:   "foo/bar",
		IssuesLink: "https://github.com/foo/bar/issues?q=is%3Aopen+is%3Aissue+no%3Alabel",
	},
}, {
	Input: "github.com/foo/bar",
	Expected: Repo{
		Owner:      "foo",
		Name:       "bar",
		FullName:   "foo/bar",
		IssuesLink: "https://github.com/foo/bar/issues?q=is%3Aopen+is%3Aissue+no%3Alabel",
	},
}, {
	Input: "https://github.com/foo/bar",
	Expected: Repo{
		Owner:      "foo",
		Name:       "bar",
		FullName:   "foo/bar",
		IssuesLink: "https://github.com/foo/bar/issues?q=is%3Aopen+is%3Aissue+no%3Alabel",
	},
}, {
	Input: "https://github.com/foo/bar/baz",
	Expected: Repo{
		Owner:      "foo",
		Name:       "bar",
		FullName:   "foo/bar",
		IssuesLink: "https://github.com/foo/bar/issues?q=is%3Aopen+is%3Aissue+no%3Alabel",
	},
},
	{
		Input: "https://notgithub.com/foo/bar",
		Error: errors.New("https://notgithub.com/foo/bar not a valid repo"),
	},
	{
		Input: "notarepo",
		Error: errors.New("notarepo not a valid repo"),
	},
}

func TestParseGithubRepo(t *testing.T) {
	for _, test := range githubTests {
		result, err := ParseGithubRepo(test.Input)
		if test.Error != nil {
			if err == nil || test.Error.Error() != err.Error() {
				t.Errorf("Expected %v, got, %v", test.Error, err)
			}
		}
		if result != test.Expected {
			t.Errorf("Expected %v, got, %v", test.Expected, result)
		}
	}
}
