package main

import (
	"crypto/elliptic"
	"encoding/gob"
	"os"

	"github.com/nheingit/learnGo/cli"
)

func main() {
	gob.Register(elliptic.P256())
	defer os.Exit(0)
	ls := []string{}
	cmd := cli.CommandLine{}
	print(cmd.CheckPopulation(ls))
}
