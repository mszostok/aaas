package fork

import (
	"github.com/goccy/go-yaml"
	"os"
)

func getKnownForks() (KnownForks, error) {
	data, err := os.ReadFile(knownForksPath)
	if err != nil {
		return KnownForks{}, err
	}

	knownForks := KnownForks{
		Forks: map[string]struct{}{},
	}
	if err = yaml.Unmarshal(data, &knownForks); err != nil {
		return KnownForks{}, err
	}

	for _, f := range knownForks.SortedForks {
		knownForks.Forks[f] = struct{}{}
	}

	return knownForks, nil
}

func storeKnownForks(knownForks KnownForks) error {
	data, err := yaml.Marshal(knownForks)
	if err != nil {
		return err
	}

	return os.WriteFile(knownForksPath, data, 0644)
}

