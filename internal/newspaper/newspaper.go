package newspaper

import (
	"errors"

	"github.com/slack-go/slack"
)

// EditAndPublish prepares and edits an issue of the newspaper containing pages
// supplied by the writers and publishes it to the publisher.
func EditAndPublish(p Publisher, writers ...Writer) error {
	var issue Issue

	// IDEA the following loop can be rewritten to leverage concurrency
	for _, w := range writers {
		page, err := w.Write()
		if err != nil {
			if errors.Is(err, ErrWriterHasNoInspiration) {
				continue
			}
			return err
		}

		issue = append(issue, page)
	}

	return p.Publish(issue)
}

// ErrWriterHasNoInspiration is used to differentiate a case when a writer has
// nothing to produce although no unexpected errors happened throughout the process.
var ErrWriterHasNoInspiration = errors.New("writer does not have enough inspiration to write")

// Publisher abstracts functionality for publishing the newspaper's issues.
type Publisher interface {
	Publish(Issue) error
}

// Issue represents a collection of pages that form an issue of the newspaper.
type Issue []Page

// Writer abstracts functionality for writing the newspaper's pages.
type Writer interface {
	Write() (Page, error)
}

// Page represents a page of the newspaper and is essentially its main building
// block.
type Page struct {
	HeadlineEmojiName string
	HeadlineText      string
	AuthorName        string
	ContentElements   []slack.Block
}
