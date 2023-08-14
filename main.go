package main

import (
	"os"

	"github.com/Sahil-4555/Golang_Chain/cli" // Import the "cli" package
)

func main() {
	defer os.Exit(0) // Ensure that the program exits gracefully
	cmd := cli.CommandLine{} // Create an instance of the CommandLine struct from the "cli" package
	cmd.Run() // Run the command-line interface using the Run() method
}
