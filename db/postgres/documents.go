package postgres

import (
	"github.com/Stinson-Moss/infengine/items"
)

// postgres db handling

type DocumentsDb struct {
	connection string // placeholder
}

func (db DocumentsDb) GetDocumentById(guid string) (*items.Document, error) {

}

func (db DocumentsDb) CreateDocument(doc *items.Document) error {
	
}

func (db DocumentsDb) UpdateDocument(guid string, doc *items.Document) error {
	
}

func (db DocumentsDb) DeleteDocument(guid string) error {
	
}