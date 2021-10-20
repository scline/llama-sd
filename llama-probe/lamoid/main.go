package main

import (
	"lamoid/alpaca"

	env "github.com/Netflix/go-env"
)

func main() {
	var environment alpaca.GrazingEnv

	es, err := env.UnmarshalFromEnviron(&environment)

}
