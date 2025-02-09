package main

import (
	"MP2/executor"
	"MP2/introducer"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	args := os.Args
	typeArg := ""
	if len(args) == 2 {
		typeArg = args[1]
	}
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	if typeArg == "introducer" {
		go introducer.StartIntroducer()
		executor.ExecuteCommand("1") // Join Group
	}
	executor.InputCommand()
}
