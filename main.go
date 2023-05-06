package main

import (
	"bufio"
	"bytes"
	"fmt"
	bs "github.com/faisal-porag/web_scraping_multi/bing_scraper"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

func saveData(keyWord string) error {
	t := time.Now().UnixNano()
	dt := time.Now()
	date := fmt.Sprintf("%s", dt.Format("01.02.2006"))
	filename := fmt.Sprintf("file_(%s)_%s_%d.txt", keyWord, date, t)
	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()
	scrapingData := ScrapingData(keyWord)
	data := []byte(scrapingData)
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func main() {
	f, err := os.Open("test_data.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)
	var keyWordList []string

	for scanner.Scan() {
		strKeyWord := scanner.Text()
		keyWordList = append(keyWordList, strKeyWord)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sliceLength := len(keyWordList)
	if sliceLength > 0 {
		var wg sync.WaitGroup
		wg.Add(sliceLength)
		fmt.Println("Running")
		for i := 0; i < sliceLength; i++ {
			go func(i int) {
				defer wg.Done()
				val := keyWordList[i]
				//fmt.Printf("i: %v, val: %v\n", i, val)
				err := saveData(val)
				log.Println(err)
			}(i)
		}
		wg.Wait()
		fmt.Println("Finished")
	} else {
		log.Println("no key word found")
	}
}

func ScrapingData(searchText string) string {
	//fmt.Println("Scrapping Start")
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
		//fmt.Println("Scrapping End")
		return dynamicString
	} else {
		log.Println(err)
		//fmt.Println("Scrapping End")
		return ""
	}
}
