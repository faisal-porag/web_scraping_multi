package main

import (
	"bytes"
	"fmt"
	bs "github.com/faisal-porag/web_scraping_multi/bing_scraper"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/index.html")
		_ = t.Execute(w, nil)
	} else {
		_ = r.ParseForm()
		searchText := r.Form["search"]
		log.Println(searchText)
		if len(searchText) > 0 && strings.TrimSpace(searchText[0]) != "" {
			t := time.Now().UnixNano()
			dt := time.Now()
			date := fmt.Sprintf("%s", dt.Format("01.02.2006"))
			filename := fmt.Sprintf("file_%s_%d.txt", date, t)
			file, err := os.Create(filename)
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "http://127.0.0.1:8084", 301)
				return
			}
			defer file.Close()
			scrapingData := ScrapingData(searchText[0])
			data := []byte(scrapingData)
			err = ioutil.WriteFile(filename, data, 0644)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			log.Println("text blank")
		}

		http.Redirect(w, r, "http://127.0.0.1:8084", 301)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	_ = http.ListenAndServe(":"+port, mux)
}

func ScrapingData(searchText string) string {
	var buffer bytes.Buffer
	var resultData string
	res, err := bs.BingScrape(searchText, "us", nil, 5, 10, 5)
	if err == nil || len(res) > 0 {
		for _, data := range res {
			resultData = fmt.Sprintf(
				"URL: %s\n\n",
				data.ResultURL,
			)
			buffer.WriteString(resultData)
		}
		dynamicString := buffer.String()
		//log.Println(dynamicString)
		return dynamicString
	} else {
		log.Println(err)
		return ""
	}
}
