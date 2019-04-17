package main

import (
	"bufio"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/go-github/v18/github"
	"github.com/juju/gnuflag"
	"github.com/mhemmings/triage/client"
	"github.com/mhemmings/triage/issues"
	"github.com/mhemmings/triage/repository"
)

// repo holds a triage.Repo and a list of Isuues associated with it.
type repo struct {
	repository.Repo
	Issues []issues.Issue
	Error  error
}

func main() {
	var port string
	var singleRepo string
	var filters client.IssueFilters
	gnuflag.StringVar(&port, "port", "8080", "HTTP port to use for serving HTML page")
	gnuflag.StringVar(&port, "p", "8080", "HTTP port to use for serving HTML page")
	gnuflag.StringVar(&singleRepo, "repo", "", "An individual repo to check")
	gnuflag.StringVar(&singleRepo, "r", "", "An individual repo to check")
	gnuflag.Var(&filters.Labels, "labels", "List of comma separated label names to filter by. By default, only issues with no labels will be shown")
	gnuflag.Var(&filters.Labels, "l", "List of comma separated label names to filter by. By default, only issues with no labels will be shown")
	gnuflag.Var(&filters.Since, "s", "Only show issues since a certain time, using the following formats such as 1m, 2h, 3d, 4w")
	gnuflag.Var(&filters.Since, "since", "Only show issues since a certain time, using the following formats: 1m, 2h, 3d, 4w")
	gnuflag.BoolVar(&filters.ShowAll, "all", false, "Show all issues. All label filters will be ignored")

	gnuflag.Parse(true)

	var repos []repo

	// Has an individual repo been provided in the command?
	if singleRepo != "" {
		r, err := repository.ParseGithubRepo(singleRepo)
		if err != nil {
			log.Fatalf("%s is not a valid Github repository", singleRepo)
		}
		repos = append(repos, repo{Repo: r})
	}

	// Has a repo list file been passed?
	if gnuflag.NArg() > 0 {
		reposFromFile, err := parseRepoListFile(gnuflag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		repos = append(repos, reposFromFile...)
	}

	ghToken := os.Getenv("TRIAGE_GITHUB_TOKEN")
	if ghToken == "" {
		// Look for legacy token
		ghToken = os.Getenv("GITHUB_TOKEN")
	}

	ghClient := client.NewGithubClient(ghToken)

	log.Printf("Collecting issues for %d repos", len(repos))

	populateIssues(ghClient, &repos, filters)

	var err error
	t := template.New("main")
	t, err = t.Parse(htmltemplate)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, templateData{Repos: repos})
	})

	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		populateIssues(ghClient, &repos, filters)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	log.Printf("Serving issue triage on http://localhost:%s\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// populateIssues takes a slice of repos and uses the provided client to repopulate the issues
// list for each repo.
func populateIssues(client client.Client, repos *[]repo, filter client.IssueFilters) {
	var wg sync.WaitGroup
	for i, repo := range *repos {
		wg.Add(1)
		i := i
		repo := repo
		go func() {
			defer wg.Done()
			var err error
			(*repos)[i].Issues, err = client.GetIssuesForTriage(context.Background(), repo.Owner, repo.Name, filter)
			(*repos)[i].Error = err
			if err != nil {
				log.Printf("Error gettings issues from %s, Error: %v\n", repo.FullName, err)
				if ghErr, ok := err.(*github.ErrorResponse); ok {
					if ghErr.Response.StatusCode == http.StatusNotFound {
						log.Println("The repository could not be found, or you are not allowed to view it")
						log.Println("Have you set TRIAGE_GITHUB_TOKEN?")
					}
				}
				return
			}
		}()
	}

	wg.Wait()
}

// parseRepoListFile takes a file name, and parses the repo list from that file
// returning a repo list. An error is returned if a file is unable to be
// opened from the provided path, or the file is not a valid repo list.
func parseRepoListFile(fileName string) ([]repo, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseRepoList(fileName, file)
}

// parseRepoList parses a repo list read from r, or an error if 1 or more repos are invalid.
// TODO: Currently all repos are assumed to be Github repos.
func parseRepoList(filename string, r io.Reader) ([]repo, error) {
	var repos []repo
	fileScanner := bufio.NewScanner(r)
	lineNum := 0

	for fileScanner.Scan() {
		lineNum++

		// TODO: Infer the correct client here, instead of just using Github.
		r, err := repository.ParseGithubRepo(fileScanner.Text())
		if err != nil {
			return nil, err
		}

		repos = append(repos, repo{Repo: r})
	}

	if err := fileScanner.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}
