package newspaper

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize/english"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"github.com/ztimes2/dailybugle/internal/mrkdwn"
)

const (
	defaultAuthorName = "J. Jonah Jameson"
)

// Ticketer abstracts functionality for accessing tickets.
type Ticketer interface {
	GetTicketsAwaitingReview() ([]Ticket, error)
}

// Ticket represents a ticket.
type Ticket struct {
	ID                 string
	URL                string
	Summary            string
	CurrentStatus      string
	CurrentStatusSince time.Time
}

func (t Ticket) daysSinceTransitionToCurrentStatus() float64 {
	return TimeNowFunc().Sub(t.CurrentStatusSince).Hours() / 24
}

// TimeNowFunc used for mocking time.Now() from outside of the package.
var TimeNowFunc = func() time.Time {
	location, _ := time.LoadLocation("Asia/Singapore")
	return time.Now().In(location)
}

// CodeReviewMarket provides functionality for writing pages related to the
// newspaper's Code Review Market topic.
type CodeReviewMarket struct {
	ticketer Ticketer
}

// NewCodeReviewMarket initializes a new CodeReviewMarket.
func NewCodeReviewMarket(t Ticketer) CodeReviewMarket {
	return CodeReviewMarket{
		ticketer: t,
	}
}

// Write implements Writer interface and generates a page containing latest
// information related to the newspaper's Code Review Market topic.
func (c CodeReviewMarket) Write() (Page, error) {
	tickets, err := c.ticketer.GetTicketsAwaitingReview()
	if err != nil {
		return Page{}, errors.Wrap(err, "could not fetch tickets")
	}

	sort.SliceStable(tickets, func(i, j int) bool {
		return tickets[i].daysSinceTransitionToCurrentStatus() >
			tickets[j].daysSinceTransitionToCurrentStatus()
	})

	p := Page{
		HeadlineEmojiName: "chart_with_upwards_trend",
		HeadlineText:      "Code Review Market",
		AuthorName:        defaultAuthorName,
	}

	if len(tickets) == 0 {
		p.ContentElements = append(p.ContentElements, slack.NewSectionBlock(
			slack.NewTextBlockObject(
				slack.MarkdownType,
				"Looks like there is no demand for code reviews today.",
				false,
				false,
			), nil, nil,
		))

		return p, nil
	}

	lines := []string{
		"Here is a list of hot tickets which index of waiting for code review is " +
			"trending up. Hurry up before someone else reviews them ahead of you!",
		"",
	}

	for _, t := range tickets {
		lines = append(lines, mrkdwn.Bold(fmt.Sprintf("   %s   +%s",
			mrkdwn.Link(t.ID, t.URL),
			english.Plural(
				int(t.daysSinceTransitionToCurrentStatus()), "day", "days",
			),
		)))
	}

	p.ContentElements = append(p.ContentElements, slack.NewSectionBlock(
		slack.NewTextBlockObject(
			slack.MarkdownType,
			strings.Join(lines, "\n"),
			false,
			false,
		), nil, nil,
	))

	return p, nil
}
