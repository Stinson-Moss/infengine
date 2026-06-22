package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"github.com/mmcdole/gofeed"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please provide a file path to a .txt list of RSS feed urls")
		return
	}

	path := os.Args[1]
	fileData, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading the file ", path)
		os.Exit(1)
	}

	numWorkers := 10
	if len(os.Args) == 3 {
		numWorkers, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid argument for number of workers. Please input a valid integer (x > 0)")
			os.Exit(1)
		}
	}

	content := string(fileData)
	var channel chan string = make(chan string)
	wg := sync.WaitGroup{}
	
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for url := range channel {
				fetcher := gofeed.NewParser()
				feed, err := fetcher.ParseURL(url)
				if err != nil {
					line := fmt.Sprintf("Error parsing url: %v", err)
					fmt.Println(line)
					continue
				}

				for _, doc := range feed.Items {
					fmt.Println(doc.Title)

					// check the repo for the same doc by guid.
					// if the doc already exists in the db, skip
					// if not, put it in there, with the UNANALYZED tag
				}
			}
		}()
	}

	urls := strings.Split(content, "\n")
	for _, url := range urls {
		channel <- url
	}

	close(channel)

	wg.Wait()

	fmt.Println("Updated database with RSS feeds")
	// fetch data
	// parse data and get tags and relevant headlines
	// form documents from data
	// put documents in database, tag as unanalyzed

	// run python script to bring in AI to analyze the data
}