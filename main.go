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
	"path/filepath"

	"github.com/Stinson-Moss/infengine/db/postgres/db"
	"github.com/Stinson-Moss/infengine/obsidian"
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

	vaultPath, ok := os.LookupEnv("VAULT_PATH")
	if !ok {
		log.Fatalln("VAULT_PATH does not exist")
	}
	vaultPath = filepath.Clean(vaultPath)
	vaultPath, err = filepath.Abs(vaultPath)
	if err != nil {
		log.Fatalln("VaultPath must be an absolute path")
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

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, time.Second * 30)
	if err := pool.Ping(timeoutCtx); err != nil {
		log.Fatalf("Error pinging connection: %v", err)
	}
	cancelTimeout()

	queries := db.New(pool)
	newDocuments := []db.Document{}
	
	channel := make(chan string)
	wg := sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for url := range channel {
				if len(url) == 0 {
					continue
				}

				fetcher := gofeed.NewParser()
				feed, err := fetcher.ParseURL(url)
				if err != nil {
					line := fmt.Sprintf("Error parsing url: %v", err)
					
					fmt.Println(line, "\nUrl: ", url)
					continue
				}

				for _, rawDoc := range feed.Items {

					ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
					defer cancel()
					_, err := queries.GetDocumentByGuid(ctx, rawDoc.GUID)
					if err != nil && err != pgx.ErrNoRows {
						fmt.Println("Error finding existing document")
						cancel()
						continue
					}

					if err == pgx.ErrNoRows {
						ctx, cancelTimeoutTx := context.WithTimeout(context.Background(), time.Minute)
						defer cancelTimeoutTx()

						transaction, _ := pool.Begin(ctx)
						defer transaction.Rollback(ctx)

						qtx := queries.WithTx(transaction)
						
						document, err := qtx.CreateDocument(ctx, db.CreateDocumentParams{
							Guid: rawDoc.GUID,
							Title: rawDoc.Title,
							Description: rawDoc.Description,
							Content: rawDoc.Content,
						})
						
						if err != nil {
							fmt.Printf("Error creating new document, %v\n", err)
						} else {
							newDocuments = append(newDocuments, document)
						}

						for _, author := range rawDoc.Authors {
							if author == nil || len(author.Name) == 0 {
								continue
							}

							if _, err := qtx.GetOrCreateAuthor(ctx, author.Name); err != nil {
								fmt.Printf("Error creating/getting author %s: %v\n", author.Name, err)
							}
						}

						for _, tag := range rawDoc.Categories {
							if len(tag) == 0 {
								continue
							}

							if _, err := qtx.GetOrCreateTag(ctx, tag); err != nil {
								fmt.Printf("Error creating/getting tag %s: %v\n", tag, err)
							}
						}

						for _, link := range rawDoc.Links {
							if len(link) == 0 {
								continue
							}

							if _, err := qtx.GetOrCreateLink(ctx, link); err != nil {
								fmt.Printf("Error creating/getting link %s: %v\n", link, err)
							}
						}


						commitCtx, cancelCommit := context.WithTimeout(context.Background(), time.Minute)
						defer cancelCommit()
						transaction.Commit(commitCtx)
		
						if err := obsidian.ExportSourceDocument(document, rawDoc.Categories, vaultPath); err != nil {
							fmt.Printf("%v\n", err)
						}
					}
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

	for _, doc := range newDocuments {
		fmt.Println(doc.Title)
	}

	// run python script to bring in AI to analyze the data
}