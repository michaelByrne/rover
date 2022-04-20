package main

import (
	"fmt"
	"log"
	"os"
	"rover/cli"
	"rover/getrover"
	"rover/provider"
)

func main() {
	fmt.Println("we in space whoa")

	apiKey := os.Getenv("NASA_API_KEY")

	prov := provider.New(provider.BaseURL, provider.WithAPIKey(apiKey))
	provWithCache := provider.NewCachedProvider(prov)
	roverGetter := getrover.New(provWithCache)

	err := cli.CLI(os.Args[1:], roverGetter)
	if err != nil {
		log.Fatal(err)
	}
}
