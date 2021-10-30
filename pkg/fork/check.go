package fork

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v38/github"
)

func Check() error {
	client := github.NewClient(nil)

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: 1000,
		},
		Sort: "name",
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(context.Background(), "mszostok", opt)
		if err != nil {
			return err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var forks []string
	for _, repo := range allRepos {
		if !repo.GetFork() {
			continue
		}

		forks = append(forks, repo.GetName())
	}

	knownForks, err := getKnownForks()
	if err != nil {
		return err
	}
	var unknown []string

	for _, f := range forks {
		if knownForks.Contains(f) {
			continue
		}
		unknown = append(unknown, f)
	}

	if len(unknown) > 0 {
		return fmt.Errorf("Found unknown forks: %s", strings.Join(unknown, ", "))
	}

	if len(knownForks.SortedForks) != len(forks) {
		// remove orphan forks
		knownForks.SortedForks = forks
		if err = storeKnownForks(knownForks); err != nil {
			return err
		}
	}

	log.Println("All looks good.")
	return nil
}
