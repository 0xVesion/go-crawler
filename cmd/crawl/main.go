package main

import (
	"fmt"

	"github.com/0xVesion/go-crawler"
)

func main() {
	links := crawler.Run("https://vincentengel.io/")

	fmt.Println("Found the following links:")
	for _, link := range links {
		fmt.Printf("\t- %s\n", link)
	}
}
