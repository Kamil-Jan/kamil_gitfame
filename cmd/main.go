package main

import (
	"log"
	"os"

	"gitfame/pkg/cli_scanner"
	"gitfame/pkg/parser"
)

func main() {
	args := os.Args[1:]
	settings := cli_scanner.Scan(args)

	parser := parser.NewParser(settings)
	err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}
}
