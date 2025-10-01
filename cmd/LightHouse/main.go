package main

import (
	"log"

	"github.com/LSariol/LightHouse/internal/cli"
	"github.com/LSariol/LightHouse/internal/scanner"
)

func main() {

	if err := run(); err != nil {
		log.Fatal(err)
	}

}

func run() error {

	if err := scanner.Initilize(); err != nil {
		return err
	}

	go scanner.Run()

	cli.StartCLI()

	// go scanner.Listen()

	// go CLI.Run()

	return nil
}
