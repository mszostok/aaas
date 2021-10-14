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
	Forks       map[string]struct{} `yaml:"-"`
	SortedForks []string            `yaml:"forks"`
}

func (k KnownForks) Contains(in string) bool {
	_, found := k.Forks[in]
	return found
}

func (k *KnownForks) AddUnique(in string) {
	if !k.Contains(in) {
		k.Forks[in] = struct{}{}
		k.SortedForks = append(k.SortedForks, in)
	}
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
	flag.CommandLine.Parse(os.Args[2:])
	if forksLit == nil || *forksLit == "" {
		exitOnError(errors.New("-forks-list flag is required"))
	}

	forks := getKnownForks()
	for _, f := range strings.SplitN(*forksLit, ",", 100) {
		forks.AddUnique(f)
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
		if knownForks.Contains(f) {
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

	knownForks := KnownForks{
		Forks: map[string]struct{}{},
	}
	exitOnError(yaml.Unmarshal(data, &knownForks))

	for _, f := range knownForks.SortedForks {
		knownForks.Forks[f] = struct{}{}
	}

	return knownForks
}

func storeKnownForks(knownForks KnownForks) {
	data, err := yaml.Marshal(knownForks)
	exitOnError(err)

	exitOnError(os.WriteFile(knownForksPath, data, 0644))
}
