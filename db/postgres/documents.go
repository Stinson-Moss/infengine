package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Stinson-Moss/infengine/items"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config = pgxpool.Config

// postgres db handling

type DocumentsDb struct {
	pool *pgxpool.Pool
}

func (db *DocumentsDb) GetDocumentById(guid string) (*items.Document, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second * 10)
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to get connection: %v", err)
	}

	defer conn.Release()

	doc := items.Document{}
	row := conn.QueryRow(ctx, `SELECT (Guid, Title, Authors, Description, Content, Tags, Links, Created) 
	FROM Documents 
	WHERE guid = $1`, guid)

	if err := row.Scan(&doc.Guid, &doc.Title, &doc.Authors, &doc.Description,
	&doc.Content, &doc.Tags, &doc.Links, &doc.Created); err != nil {
		log.Fatal(err)
	}


}

func (db *DocumentsDb) CreateDocument(doc *items.Document) error {
	
}

func (db *DocumentsDb) UpdateDocument(guid string, doc *items.Document) error {
	
}

func (db *DocumentsDb) DeleteDocument(guid string) error {
	
}

func CreateDB(config *Config) (*DocumentsDb, error) {
	if config == nil {
		return nil, fmt.Errorf("Config pointer is nil")
	}

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating pgx pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("Error making connection: %v", ctx.Err())
	}

	return &DocumentsDb{pool: pool}, nil
}