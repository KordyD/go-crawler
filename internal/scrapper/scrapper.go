package scrapper

import (
	"log"
	urlpkg "net/url"
	"strings"

	"github.com/PuerkitoBio/purell"
	"github.com/gocolly/colly"
	"github.com/kordyd/go-crawler/internal/entities"
)

// TODO: Research colly async
// TODO: Add normalization and filtration of urls (like "#")
// TODO: Add parse html to json
func Scrapper(url string, parsedData chan<- entities.Url, parsedUrls chan<- string) {
	log.Println("Scrapper is running")

	var result entities.Url

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		contentType := r.Headers.Get("Content-Type")
		if strings.Contains(contentType, "text/html") {
			result.Content = string(r.Body)
		} else {
			log.Println("Not HTML content")
			result.Content = ""
		}
	})

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		base, err := urlpkg.Parse(url)
		if err != nil {
			log.Println(err)
		}

		u, err := urlpkg.Parse(link)
		if err != nil {
			log.Println(err)
		}

		urlToNormalize := base.ResolveReference(u)

		normilizedUrl := purell.NormalizeURL(urlToNormalize, purell.FlagsUnsafeGreedy)

		parsedUrls <- normilizedUrl
	})

	c.OnError(func(r *colly.Response, err error) {
		result.Error = err.Error()
	})

	c.OnScraped(func(r *colly.Response) {
		result.Parsed = true
		parsedData <- result
	})

	c.Visit(url)
}
