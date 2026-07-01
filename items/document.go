package items

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

type Document struct {
	Guid string
	Title string
	
	Authors []string
	Description string
	Content string
	
	Tags []string
	Links []string
	Created *time.Time
}

func GetNames(list []*gofeed.Person) []string {
	names := []string{}
	for _, person := range list {
		if person == nil {
			continue
		}

		names = append(names, person.Name)
	}

	return names
}

func GetCreationTime(feed *gofeed.Item) (*time.Time, error) {
	if feed == nil {
		return nil, fmt.Errorf("feed is nil")
	}

	if feed.UpdatedParsed != nil {
		return feed.UpdatedParsed, nil
	}

	if feed.PublishedParsed != nil {
		return feed.PublishedParsed, nil
	}

	return nil, fmt.Errorf("No valid parsed date")
}