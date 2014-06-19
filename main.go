package main

import (
	"flag"
	crowdflow "github.com/abhiyerra/workmachine/crowdflow"
	"log"
	"os"
)

var (
	Program crowdflow.Program
)

func init() {
	log.Println("CrowdFlow Starting...")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Println("Crowdflow program file is missing.")
		os.Exit(1)
	}

	jsonStream, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}

	Program = crowdflow.ParseProgram(jsonStream)
}

func main() {
	Program.Run()
}
