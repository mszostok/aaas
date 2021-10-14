package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/google/go-github/v38/github"
)

const knownForksPath = "./known-forks.yaml"

type KnownForks struct {
	Forks map[string]struct{} `yaml:"forks"`
}

func main() {
	switch cmd := os.Args[1]; cmd {
	case "add":
		addKnownFork()
	case "check":
		checkForks()
	default:
		exitOnError(fmt.Errorf("unkown command %v: allowed %s, %s", cmd, "add", "check"))
	}
}

func addKnownFork() {
	forksLit := flag.String("forks-list", "", "Comma separated list of known forks to add")
	flag.Parse()
	if forksLit == nil || *forksLit == "" {
		exitOnError(errors.New("-forks-list flag is required"))
	}

	forks := getKnownForks()
	for _, f := range strings.SplitN(*forksLit, ",", 100) {
		forks.Forks[strings.TrimSpace(f)] = struct{}{}
	}

	storeKnownForks(forks)
}

func checkForks() {
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
		exitOnError(err)
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

	knownForks := getKnownForks()
	var unknown []string
	for _, f := range forks {
		_, isKnown := knownForks.Forks[f]
		if isKnown {
			continue
		}
		unknown = append(unknown, f)
	}

	if len(unknown) > 0 {
		log.Fatalf("Found unknown forks: %s", strings.Join(unknown, ", "))
	}

	log.Println("All looks good.")
}

func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getKnownForks() KnownForks {
	data, err := os.ReadFile(knownForksPath)
	exitOnError(err)

	knownForks := KnownForks{}
	exitOnError(yaml.Unmarshal(data, &knownForks))
	return knownForks
}

func storeKnownForks(knownForks KnownForks) {
	data, err := yaml.Marshal(knownForks)
	exitOnError(err)

	exitOnError(os.WriteFile(knownForksPath, data, 0644))
}
