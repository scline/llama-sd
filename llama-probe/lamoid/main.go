package main

import (
	"errors"
	"lamoid/alpaca"
	"log"

	env "github.com/Netflix/go-env"
)

//TODO: Pass HTTP Client to methods

func main() {

	var llama alpaca.LamoidEnv

	_, err := env.UnmarshalFromEnviron(&llama)
	if err != nil || errors.Is(err, env.ErrUnexportedField) {
		log.Fatalf("[ENV-ERR]: There was a problem with one or more expected environment: %s", err)
	}

	log.Print(llama)

	llama.Graze()
}
