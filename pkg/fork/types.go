package fork

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
	if k.Contains(in) {
		return
	}

	k.Forks[in] = struct{}{}
	k.SortedForks = append(k.SortedForks, in)
}
