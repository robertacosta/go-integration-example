package main

import (
	"log"

	"github.com/robertacosta/go-integration-example/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Printf("unable to execute command: %s", err)
	}
}
