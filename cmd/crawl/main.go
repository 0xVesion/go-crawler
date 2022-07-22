package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

var queue chan string = make(chan string)
var requestedLinksSet map[string]bool = map[string]bool{}
var wg sync.WaitGroup = sync.WaitGroup{}
var lock sync.Mutex = sync.Mutex{}

const maxRequests = 100

func main() {
	go func() {
		for current := range queue {
			requestCount := len(requestedLinksSet)
			if requestCount == maxRequests {
				continue
			}

			isAlreadyRequested := requestedLinksSet[current]
			if isAlreadyRequested {
				continue
			}

			wg.Add(1)
			requestedLinksSet[current] = true
			go doWork(current)
		}
	}()

	queue <- "https://vincentengel.io/"

	wg.Wait()
	fmt.Println("Found the following links:")
	for link := range requestedLinksSet {
		fmt.Printf("\t- %s\n", link)
	}
}

func doWork(current string) {
	defer wg.Done()
	fmt.Printf("Working on %q...\n", current)

	body := getBody(current)
	links := extractLinks(body)

	for _, link := range links {
		queue <- link
	}
}

func extractLinks(str []byte) []string {
	linksPattern, err := regexp.Compile(`(https?):\/\/([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:\/~+#-]*[\w@?^=%&\/~+#-])`)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	linksBytes := linksPattern.FindAll(str, -1)

	links := []string{}
	for _, link := range linksBytes {
		links = append(links, string(link))
	}

	return links
}

func getBody(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}

	return body
}
