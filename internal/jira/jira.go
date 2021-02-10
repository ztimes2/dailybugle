package jira

import (
	"github.com/andygrunwald/go-jira"
	"github.com/ztimes2/dailybugle/internal/newspaper"
)

// Jira provides communication with Jira API.
type Jira struct {
	client  *jira.Client
	baseURL string
}

// Config holds Jira's configuration.
type Config struct {
	BaseURL  string
	Username string
	APIToken string
}

// New initializes a new Jira.
func New(conf Config) (Jira, error) {
	t := jira.BasicAuthTransport{
		Username: conf.Username,
		Password: conf.APIToken,
	}

	client, err := jira.NewClient(t.Client(), conf.BaseURL)
	if err != nil {
		return Jira{}, err
	}

	return Jira{
		client:  client,
		baseURL: conf.BaseURL,
	}, nil
}

// GetTicketsAwaitingReview implements newspaper.Ticketer interface and fetches
// Jira tickets that are waiting for code review.
func (j Jira) GetTicketsAwaitingReview() ([]newspaper.Ticket, error) {
	var tickets []newspaper.Ticket

	if err := j.client.Issue.SearchPages(
		`project = "Mobile Backend" AND status = "Awaiting Review"`,
		&jira.SearchOptions{
			StartAt:    0,
			MaxResults: 50,
			Expand:     "changelog",
			Fields:     []string{"summary", "status"},
		},
		func(i jira.Issue) error {
			t := newspaper.Ticket{
				ID:            i.Key,
				URL:           j.toURL(i),
				Summary:       i.Fields.Summary,
				CurrentStatus: i.Fields.Status.Name,
			}

			h, _ := getTransitionToCurrentStatus(i)
			t.CurrentStatusSince, _ = h.CreatedTime()

			tickets = append(tickets, t)
			return nil
		},
	); err != nil {
		return nil, err
	}

	return tickets, nil
}

func getTransitionToCurrentStatus(i jira.Issue) (jira.ChangelogHistory, bool) {
	for _, history := range i.Changelog.Histories {
		for _, item := range history.Items {
			if item.Field == "status" && item.ToString == i.Fields.Status.Name {
				return history, true
			}
		}
	}

	return jira.ChangelogHistory{}, false
}

func (j Jira) toURL(i jira.Issue) string {
	return j.baseURL + "/browse/" + i.Key
}
