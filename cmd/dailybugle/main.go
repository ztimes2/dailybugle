package main

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/ztimes2/dailybugle/internal/config"
	"github.com/ztimes2/dailybugle/internal/google"
	"github.com/ztimes2/dailybugle/internal/jira"
	"github.com/ztimes2/dailybugle/internal/newspaper"
	"github.com/ztimes2/dailybugle/internal/slack"
	"golang.org/x/oauth2"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		handleError(errors.Wrap(err, "could not load configuration"))
		return
	}

	jiraClient, err := jira.New(jira.Config{
		BaseURL:  cfg.JiraBaseURL,
		Username: cfg.JiraUsername,
		APIToken: cfg.JiraAPIToken,
	})
	if err != nil {
		handleError(errors.Wrap(err, "could not init Jira client"))
		return
	}

	calendars, err := initCalendars(cfg)
	if err != nil {
		handleError(err)
		return
	}

	if err := newspaper.EditAndPublish(
		slack.NewChannel(cfg.SlackAPIToken, cfg.SlackChannelID),
		newspaper.NewCodeReviewMarket(jiraClient),
		newspaper.NewReleaseForecast(calendars...),
	); err != nil {
		handleError(err)
		return
	}
}

func handleError(err error) {
	panic(err)
}

func initCalendars(cfg config.Config) ([]newspaper.Calendar, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(cfg.GoogleAccessToken), &token); err != nil {
		return nil, errors.Wrap(err, "could not parse Google's access token")
	}

	client := google.NewClient(cfg.GoogleClientID, cfg.GoogleClientSecret, &token)

	crmDispatches, err := google.NewCRMDispatchesCalendar(
		client, cfg.CRMDispatchesCalendarID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not init CRM dispatches calendar")
	}

	campaigns, err := google.NewCampaignCalendar(
		client, cfg.CampaignsCalendarID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not init campaigns calendar")
	}

	devMilestones, err := google.NewDevMilestonesCalendar(
		client, cfg.DevMilestonesCalendarID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not init dev milestones calendar")
	}

	return []newspaper.Calendar{
		crmDispatches, campaigns, devMilestones,
	}, nil
}
