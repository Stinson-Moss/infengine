package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"github.com/mmcdole/gofeed"
	"github.com/Stinson-Moss/infengine/items"
	"github.com/Stinson-Moss/infengine/db/postgres"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("Please provide a file path to a .txt list of RSS feed urls")
		return
	}

	path := os.Args[1]
	fileData, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("Error reading the file ", path)
	}

	numWorkers, err := strconv.Atoi(os.Getenv("WORKER_COUNT"))
	if err != nil {
		log.Fatalln("WORKER_COUNT environment variable has not been set correctly. Please put in an integer value")
	}

	db, err := postgres.CreateDB()
	
	channel := make(chan string)
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

				for _, rawDoc := range feed.Items {
					doc, err := items.FromFeed(rawDoc)
					if err != nil {
						fmt.Println("Unable to process feed", rawDoc.Title)
						continue
					}

					existingDoc, err := db
					

					// check the repo for the same doc by guid.
					// if the doc already exists in the db, skip
					// if not, put it in there, with the UNANALYZED tag
				}
			}
		}()
	}

	urls := strings.Split(string(fileData), "\n")
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