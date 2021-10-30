package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mszostok/aaas/pkg/fork"
)

func main() {
	switch cmd := os.Args[1]; cmd {
	case "add":
		exitOnError(fork.AddKnownFork())
	case "check":
		exitOnError(fork.Check())
	default:
		exitOnError(fmt.Errorf("unkown command %v: allowed %s, %s", cmd, "add", "check"))
	}
}

func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
