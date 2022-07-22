package crawler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
)

const maxRequests = 999

type crawler struct {
	queue    chan string
	urls     map[string]bool
	wg       sync.WaitGroup
	startUrl string
}

func Run(startUrl string) []string {
	c := &crawler{
		queue:    make(chan string),
		urls:     map[string]bool{},
		startUrl: startUrl,
	}

	go c.workOnQueue()

	c.queue <- startUrl

	c.wg.Wait()

	links := []string{}
	for link := range c.urls {
		links = append(links, link)
	}

	return links
}

func (c *crawler) workOnQueue() {
	for current := range c.queue {
		requestCount := len(c.urls)
		if requestCount == maxRequests {
			continue
		}

		isAlreadyRequested := c.urls[current]
		if isAlreadyRequested {
			continue
		}

		c.wg.Add(1)
		c.urls[current] = true
		go c.handleUrl(current)
	}
}

func (c *crawler) handleUrl(current string) {
	defer c.wg.Done()
	log.Printf("Working on %q...\n", current)

	body := getBody(current)
	links := extractLinks(body)

	for _, link := range links {
		c.queue <- link
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
	http.DefaultClient.Timeout = time.Second

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
