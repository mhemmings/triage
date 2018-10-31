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

	"github.com/juju/gnuflag"
	"github.com/mhemmings/triage"
)

// repo holds a triage.Repo and a list of Isuues associated with it.
type repo struct {
	triage.Repo
	Issues []triage.Issue
}

func main() {
	var port string
	var singleRepo string
	gnuflag.StringVar(&port, "port", "8080", "HTTP port to use for serving HTML page")
	gnuflag.StringVar(&port, "p", "8080", "HTTP port to use for serving HTML page")
	gnuflag.StringVar(&singleRepo, "repo", "", "An individual repo to check")
	gnuflag.StringVar(&singleRepo, "r", "", "An individual repo to check")

	gnuflag.Parse(true)

	var repos []repo

	// Has an individual repo been provided in the command?
	if singleRepo != "" {
		r, err := triage.ParseGithubRepo(singleRepo)
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

	client := triage.NewGithubClient(os.Getenv("GITHUB_TOKEN"))

	var wg sync.WaitGroup
	for i, repo := range repos {
		wg.Add(1)
		i := i
		repo := repo
		log.Println("Searching for issues in:", repo.FullName)
		go func() {
			defer wg.Done()
			var err error
			repos[i].Issues, err = client.GetIssuesForTriage(context.Background(), repo.Owner, repo.Name)
			if err != nil {
				log.Errorf("Error gettings issurs from %s, Error: %v", repo.FullName, err)
				return
			}
			log.Printf("Found %d issues in %s", len(repos[i].Issues), repos[i].FullName)
		}()
	}

	wg.Wait()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.New("main")    //name of the template is main
		t, _ = t.Parse(htmltemplate) // parsing of template string
		t.Execute(w, templateData{Repos: repos})
	})

	log.Printf("Serving issue triage on http://localhost:%s\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
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
		r, err := triage.ParseGithubRepo(fileScanner.Text())
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
