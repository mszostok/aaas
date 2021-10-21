# aaas

This repository contains various checks/functions executed via GitHub Action.

### Workflows:
- [add_know_forks.yml](.github/workflows/add_know_forks.yml) - Manually executed action that takes a comma separated list of forks that should be accepted and not reported as unwanted.
- [check_forks.yml](.github/workflows/check_forks.yml) - Periodic job that detects if I have unwanted forks under my profile. It makes easier to detect that some repository was forked by a mistake. For example, when you click the edit button on some file in another repository, GitHub automatically forks it. Personally, I like to have a short list of forked repositories and decide each time if I really need it or not.
