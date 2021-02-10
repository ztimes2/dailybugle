package newspaper

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"github.com/ztimes2/dailybugle/internal/mrkdwn"
)

// Calendar abstract functionality of retrieving calendar events from any source.
type Calendar interface {
	GetCalendarEventsByDay(time.Time) ([]CalendarEvent, error)
}

// CalendarEvent represents a calendar event.
type CalendarEvent struct {
	Title    string
	Type     CalendarEventType
	StartsAt time.Time
	EndsAt   time.Time
}

// CalendarEventType is a type of a calendar type.
type CalendarEventType int

const (
	// CalendarEventTypeUndefined is a calendar event type that is not supported
	// by the application.
	CalendarEventTypeUndefined CalendarEventType = iota

	// CalendarEventTypePushNotification is a calendar event type related to push
	// notifications.
	CalendarEventTypePushNotification

	// CalendarEventTypeCampaign is a calendar event type related to campaigns.
	CalendarEventTypeCampaign

	// CalendarEventTypeCodeFreeze is a calendar event type related to code freezes.
	CalendarEventTypeCodeFreeze
)

// ReleaseForecast provides functionality for writing pages for the newspaper's
// Release Forecast topic.
type ReleaseForecast struct {
	calendars []Calendar
}

// NewReleaseForecast initializes a new ReleaseForecast.
func NewReleaseForecast(calendars ...Calendar) ReleaseForecast {
	return ReleaseForecast{
		calendars: calendars,
	}
}

// Write implements Writer interface and generates a page containing latest
// information related to the newspaper's Release Forecast topic.
func (r ReleaseForecast) Write() (Page, error) {
	var events []CalendarEvent

	for _, c := range r.calendars {
		e, err := c.GetCalendarEventsByDay(TimeNowFunc())
		if err != nil {
			return Page{}, err
		}
		events = append(events, e...)
	}

	var pushNotifications, campaigns, codeFreezes []CalendarEvent

	for _, e := range events {
		switch e.Type {
		case CalendarEventTypePushNotification:
			pushNotifications = append(pushNotifications, e)
		case CalendarEventTypeCampaign:
			campaigns = append(campaigns, e)
		case CalendarEventTypeCodeFreeze:
			codeFreezes = append(codeFreezes, e)
		}
	}

	p := Page{
		HeadlineEmojiName: "sun_behind_rain_cloud",
		HeadlineText:      "Release Forecast",
		AuthorName:        defaultAuthorName,
	}

	lines := []string{
		getReleaseForecastSummary(pushNotifications, campaigns, codeFreezes) + " " +
			getReleaseForecastRecommendation(pushNotifications, campaigns, codeFreezes),
	}

	if len(codeFreezes) > 0 {
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

	if len(pushNotifications) > 0 {
		breakdown := []string{
			"",
			"Thunderstorm of Push Notifications is expected during the " +
				"following hours (SGT):",
		}

		for _, pn := range mergeOverlappingCalendarEvents(pushNotifications) {
			breakdown = append(breakdown, mrkdwn.Bold(
				"    "+pn.StartsAt.Format(time.Kitchen),
			))
		}

		lines = append(lines, breakdown...)
	}

	if len(campaigns) > 0 {
		breakdown := []string{
			"",
			"Heavy rain of Campaigns is expected during the following hours (SGT):",
		}

		for _, c := range mergeOverlappingCalendarEvents(campaigns) {
			breakdown = append(breakdown, mrkdwn.Bold(fmt.Sprintf(
				"    %s - %s",
				c.StartsAt.Format(time.Kitchen),
				c.EndsAt.Format(time.Kitchen),
			)))
		}

		lines = append(lines, breakdown...)
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

func getReleaseForecastSummary(pushNotifications, campaigns, codeFreezes []CalendarEvent,
) string {

	if len(codeFreezes) > 0 {
		return mrkdwn.Emoji("snowflake") +
			" The day is freezingly cold due to the Code Freeze."
	}

	if len(pushNotifications) > 0 && len(campaigns) > 0 {
		return mrkdwn.Emoji("thunder_cloud_and_rain") +
			" The day is cloudy due to Push Notifications and Campaigns."
	}

	if len(pushNotifications) > 0 {
		return mrkdwn.Emoji("thunder_cloud_and_rain") +
			" The day is cloudy due to Push Notifications."
	}

	if len(campaigns) > 0 {
		return mrkdwn.Emoji("thunder_cloud_and_rain") +
			" The day is cloudy due to Campaigns."
	}

	return mrkdwn.Emoji("sunny") + " The day is sunny and the sky is clear."
}

func getReleaseForecastRecommendation(
	pushNotifications, campaigns, codeFreezes []CalendarEvent,
) string {

	if len(codeFreezes) > 0 {
		return "Totally bad day for a release!"
	}

	if len(pushNotifications) > 0 || len(campaigns) > 0 {
		return "Be careful with a release today!"
	}

	return "Looks like a good day for a release!"
}

func mergeOverlappingCalendarEvents(events []CalendarEvent) []CalendarEvent {
	if len(events) <= 1 {
		return events
	}

	merged := append([]CalendarEvent(nil), events...)

	// Sorts events by time intervals in ascending order.
	sort.SliceStable(merged, func(i, j int) bool {
		if merged[i].StartsAt.Before(merged[j].StartsAt) {
			return true
		}
		if merged[i].StartsAt.Equal(merged[j].StartsAt) &&
			merged[i].EndsAt.Before(merged[j].EndsAt) {
			return true
		}
		return false
	})

	j := 0
	for i := 1; i < len(merged); i++ {
		if merged[i].StartsAt.Before(merged[j].EndsAt) ||
			merged[i].StartsAt.Equal(merged[j].EndsAt) {
			if merged[j].EndsAt.Before(merged[i].EndsAt) {
				merged[j].EndsAt = merged[i].EndsAt
			}
		} else {
			j++
			merged[j] = merged[i]
		}

	}

	return append([]CalendarEvent(nil), merged[:j+1]...)
}
