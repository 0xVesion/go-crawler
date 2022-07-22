package main

import (
	"fmt"
	"os"

	"github.com/0xVesion/go-crawler"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("💥 Start url not supplied!")
		printUsage()

		return
	}

	startUrl := os.Args[1]
	links := crawler.Run(startUrl)

	printResult(links)
}

func printResult(links []string) {
	fmt.Println("Found the following links:")
	for _, link := range links {
		fmt.Printf("\t- %s\n", link)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcrawl [start-url]")
	fmt.Println("Example:")
	fmt.Println("\tcrawl https://vincentengel.io")
}
