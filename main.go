package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Stinson-Moss/infengine/db/postgres/db"
	"github.com/Stinson-Moss/infengine/items"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmcdole/gofeed"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("Please provide a file path to a .txt list of RSS feed urls")
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

	pgUrl, ok := os.LookupEnv("POSTGRES_URL")
	if !ok {
		log.Fatal("POSTGRES_URL env variable is null")
	}

	config, err := pgxpool.ParseConfig(pgUrl)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 10 * time.Minute

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		log.Fatalf("Unable to instantiate connection pool: %v", err)
	}

	timeoutCtx, _ := context.WithTimeout(ctx, time.Minute * 1)
	if err := pool.Ping(timeoutCtx); err != nil {
		log.Fatalf("Error pinging connection: %v", err)
	}

	db := db.New(pool)
	
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

					ctx, _ := context.WithTimeout(context.Background(), 5 * time.Minute)
					_, err := db.GetDocumentByGuid(ctx, doc.Guid)
					if err != nil && err != pgx.ErrNoRows {
						fmt.Println("Error finding existing document")
						continue
					}

					if err == pgx.ErrNoRows {
						db
						db.CreateDocument()
					}

					if existingDoc != nil

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