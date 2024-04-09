package scrapper

import (
	"log"

	"github.com/gocolly/colly"
)

// TODO: Research colly async
// TODO: Add normalization and filtration of urls (like "#")
// TODO: Add parse html to json
func Scrapper(url string, parsedUrls chan<- string, fetchedBody chan<- string) {
	// , wg *sync.WaitGroup
	log.Println("Scrapper is running")

	// defer wg.Done()

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		fetchedBody <- string(r.Body)

	})

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		parsedUrls <- link
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println(err)
	})

	c.Visit(url)
}
