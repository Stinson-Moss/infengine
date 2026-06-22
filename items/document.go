package items

import (
	"time"
)

type Document struct {
	GUID string
	Title string
	
	Authors []string
	Tags []string
	Links []string
	Description string
	Content string
	Created *time.Time
}

type DocumentsRepository interface {

	// CRUD
	GetDocumentById(guid string) (*Document, error)
	CreateDocument(*Document) error
	UpdateDocument(guid string, doc *Document) error
	DeleteDocument(guid string) error

	// 
}