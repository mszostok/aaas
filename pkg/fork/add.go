package fork

import (
	"errors"
	"flag"
	"os"
	"strings"
)

func AddKnownFork() error {
	forksLit := flag.String("forks-list", "", "Comma separated list of known forks to add")
	if err := flag.CommandLine.Parse(os.Args[2:]); err != nil {
		return err
	}
	if forksLit == nil || *forksLit == "" {
		return errors.New("-forks-list flag is required")
	}

	forks, err := getKnownForks()
	if err != nil {
		return err
	}

	for _, f := range strings.SplitN(*forksLit, ",", 100) {
		forks.AddUnique(f)
	}

	return storeKnownForks(forks)
}
