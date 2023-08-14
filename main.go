package main

import (
	"os"

	"github.com/Sahil-4555/Golang_Chain/cli"
)

func main() {
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()
}