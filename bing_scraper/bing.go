package bing_scraper

import (
	"fmt"
	"github.com/oklog/ulid/v2"
	"io"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func buildBingUrls(searchTerm, country string, pages, count int) ([]string, error) {
	var toScrape []string
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	uniqueId := NewULidString()
	if countryCode, found := bingDomains[country]; found {
		for i := 0; i < pages; i++ {
			first, formPERE := firstParameter(i, count)
			redigId := NewULidString()
			scrapeURL := fmt.Sprintf("https://bing.com/search?q=%s&qs=n&sp=-1&lq=0&pq=%s&sc=0-23&sk=&cvid=%s&count=%d%s&ghsh=0&ghacc=0&ghpl=&toWww=1&redig=%s&first=%d%s", searchTerm, searchTerm, uniqueId, count, countryCode, redigId, first, formPERE)
			log.Println(scrapeURL)
			toScrape = append(toScrape, scrapeURL)
		}
	} else {
		err := fmt.Errorf("country (%s) is currently not supported", country)
		return nil, err
	}
	return toScrape, nil
}

func firstParameter(number, count int) (int, string) {
	if number == 0 {
		return number + 1, ""
	} else if number == 1 {
		return number*count + 1, "&FORM=PERE"
	}
	return number*count + 1, fmt.Sprintf("&FORM=PERE%d", number-1)
}

func bingResultParser(response io.Reader, rank int) ([]SearchResult, error) {
	doc, err := goquery.NewDocumentFromReader(response)
	if err != nil {
		return nil, err
	}
	var results []SearchResult
	sel := doc.Find("main #b_results li.b_algo")
	rank++
	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("h2")
		descTag := item.Find("div.b_caption p")
		desc := descTag.Text()
		title := titleTag.Text()
		link = strings.Trim(link, " ")
		if link != "" && link != "#" && !strings.HasPrefix(link, "/") {
			result := SearchResult{
				rank,
				link,
				title,
				desc,
			}
			results = append(results, result)
			rank++
		}
	}
	return results, err
}

// BingScrape scrapes bing.com with desired parameters
func BingScrape(searchTerm, country string, proxyString interface{}, pageSize, noOfDataPerPage, sleepTime int) ([]SearchResult, error) {
	var results []SearchResult
	bingPages, err := buildBingUrls(searchTerm, country, pageSize, noOfDataPerPage)
	if err != nil {
		return nil, err
	}
	for _, page := range bingPages {
		rank := len(results)
		resIoReader, errRequest := scrapeClientRequest(page, proxyString)
		if errRequest != nil {
			return nil, errRequest
		}
		data, err := bingResultParser(resIoReader, rank)
		if err != nil {
			return nil, err
		}
		for _, result := range data {
			results = append(results, result)
		}
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
	return results, nil
}

func NewULidString() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}
