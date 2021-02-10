package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
)

// Config holds the application's configuration variables.
type Config struct {
	GoogleClientID     string `config:"GOOGLE_CLIENT_ID,required"`
	GoogleClientSecret string `config:"GOOGLE_CLIENT_SECRET,required"`
	GoogleAccessToken  string `config:"GOOGLE_ACCESS_TOKEN,required"`

	CRMDispatchesCalendarID string `config:"CRMDISPATCHES_CALENDAR_ID,required"`
	CampaignsCalendarID     string `config:"CAMPAIGNS_CALENDAR_ID,required"`
	DevMilestonesCalendarID string `config:"DEVMILESTONES_CALENDAR_ID,required"`

	SlackAPIToken  string `config:"SLACK_API_TOKEN,required"`
	SlackChannelID string `config:"SLACK_CHANNEL_ID,required"`

	JiraBaseURL  string `config:"JIRA_BASE_URL,required"`
	JiraUsername string `config:"JIRA_USERNAME,required"`
	JiraAPIToken string `config:"JIRA_API_TOKEN,required"`
}

// Load loads the application's configuration.
func Load() (Config, error) {
	var cfg Config

	if err := confita.NewLoader(
		env.NewBackend(),
		newDotEnv(),
	).Load(context.Background(), &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
