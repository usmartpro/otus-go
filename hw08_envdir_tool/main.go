package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Wrong arguments count")
		// os.Exit(1)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalf("Failed read %s: %s", os.Args[1], err.Error())
		// os.Exit(1)
	}

	code := RunCmd(os.Args[2:], env)

	os.Exit(code)
}
