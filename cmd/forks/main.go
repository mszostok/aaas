package main

import (
	"context"
	"github.com/google/go-github/v38/github"
	"log"
	"strings"
)

var knownForks = map[string]struct{}{
	".github-capact": {},
	"addons": {},
	"brokerapi": {},
	"capact": {},
	"cli": {},
	"client-go": {},
	"community": {},
	"community-1": {},
	"controller-runtime": {},
	"gcp-service-broker": {},
	"go-zglob": {},
	"helm": {},
	"helm-docs": {},
	"helm-operator": {},
	"hub-manifests": {},
	"k3d": {},
	"kubernetes": {},
	"kyma": {},
	"markdownfmt": {},
	"mdsh": {},
	"node-env-configuration": {},
	"octopus": {},
	"org": {},
	"podpreset-crd": {},
	"prometheus-operator": {},
	"Quick-Start-Big-Bang": {},
	"quicktype": {},
	"rafter": {},
	"service-broker-plugins": {},
	"service-catalog": {},
	"service-catalog-tester": {},
	"shelldoc": {},
	"slides": {},
	"utils": {},
	"vimium": {},
	"website": {},
}

func main() {
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

	var unknown []string
	for _, f := range forks {
		_, isKnown := knownForks[f]
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
