package main

import (
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/nheingit/learnGo/cli"
)

func main() {
	gob.Register(elliptic.P256())
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Populate()

	inputls := []string{}
	inputls = append(inputls, "PREVENTIS SALUD", "SONRI DENT", "ozonature", "CLINICA JESUS MARIA")
	outputls := cmd.CheckPopulation(inputls)
	fmt.Print(outputls)

}
