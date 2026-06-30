package items

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

type Document struct {
	Id int64
	Guid string
	Title string
	
	Authors []string
	Description string
	Content string
	
	Tags []string
	Links []string
	Created *time.Time
}

func getNames(list []*gofeed.Person) []string {
	names := []string{}
	for _, person := range list {
		if person == nil {
			continue
		}

		names = append(names, person.Name)
	}

	return names
}

func FromFeed(feed *gofeed.Item) (*Document, error) {
	if feed == nil {
		return nil, fmt.Errorf("Feed is nil")
	}

	doc := Document{
		Id: -1,
		Guid: feed.GUID,
		Title: feed.Title,

		Authors: getNames(feed.Authors),
		Description: feed.Description,
		Content: feed.Content,
		
		Tags: feed.Categories,
		Links: feed.Links,
		Created: feed.PublishedParsed,
	}

	if doc.Created == nil {
		doc.Created = feed.UpdatedParsed
	}

	return &doc, nil
} 