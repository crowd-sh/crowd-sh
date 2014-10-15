package main

import (
	"flag"
	crowdflow "github.com/workmachine/workmachine/crowdflow"
	"log"
	"os"
)

var (
	Program crowdflow.Program
)

func init() {
	flag.Parse()
}

func main() {
	log.Println("CrowdFlow Starting...")

	args := flag.Args()
	if len(args) < 2 {
		log.Println("Crowdflow program file and data csv are missing.")
		os.Exit(1)
	}

	jsonProgram, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}

	csvFile, err := os.Open(args[1])
	if err != nil {
		log.Fatal(err)
	}

	Program = crowdflow.ParseProgram(jsonProgram)
	Program.LoadJobs(csvFile)
	Program.Run()
}
