package main

import (
	"lamoid/alpaca"
	"log"

	env "github.com/Netflix/go-env"
)

func main() {
	var llama alpaca.LamoidEnv

	_, err := env.UnmarshalFromEnviron(&llama)
	if err != nil || err == env.ErrUnexportedField {
		log.Fatalf("[ENV-ERR]: There was a problem with one or more expected environment")
	}

	llama.Graze()
}
