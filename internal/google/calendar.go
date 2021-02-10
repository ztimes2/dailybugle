package google

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/ztimes2/dailybugle/internal/newspaper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const (
	crmDispatchesCalendarID = "zalora.com_vjahlu1fmi8jnmetoclsl75asg@group.calendar.google.com"
	campaignsCalendarID     = "zalora.com_6km7sa5pr2qhhg3k14kc0tgb4c@group.calendar.google.com"
	devMilestonesCalendarID = "zalora.com_n0shb4lud7nfmavj89qq003vu8@group.calendar.google.com"
)

// NewClient returns a new http.Client that uses the given authentication credentials
// when making HTTP requests. If the provided OAuth2 token contains a refresh token,
// then it will automatically be refreshed after its expiry.
func NewClient(clientID, clientSecret string, t *oauth2.Token) *http.Client {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
	}
	return conf.Client(context.Background(), t)
}

// CRMDispatchesCalendar provides communication with Zalora's CRM Dispatches calendar
// where Push Notification events are kept.
type CRMDispatchesCalendar struct {
	service    *calendar.Service
	calendarID string
}

// NewCRMDispatchesCalendar initializes a new CRMDispatchesCalendar.
func NewCRMDispatchesCalendar(c *http.Client, calendarID string,
) (CRMDispatchesCalendar, error) {

	s, err := calendar.NewService(context.Background(), option.WithHTTPClient(c))
	if err != nil {
		return CRMDispatchesCalendar{}, err
	}

	return CRMDispatchesCalendar{
		service:    s,
		calendarID: calendarID,
	}, nil
}

// GetCalendarEventsByDay returns a list of events from the calendar scheduled
// for the given day around the working hours (9AM-7PM SGT).
func (c CRMDispatchesCalendar) GetCalendarEventsByDay(t time.Time,
) ([]newspaper.CalendarEvent, error) {

	start := time.Date(t.Year(), t.Month(), t.Day(), 9, 0, 0, 0, t.Location())
	end := time.Date(t.Year(), t.Month(), t.Day(), 19, 0, 0, 0, t.Location())

	resp, err := c.service.Events.
		List(c.calendarID).
		TimeMin(start.Format(time.RFC3339)).
		TimeMax(end.Format(time.RFC3339)).
		Do()
	if err != nil {
		return nil, err
	}

	var events []newspaper.CalendarEvent

	for _, item := range resp.Items {
		e := newspaper.CalendarEvent{
			Title: item.Summary,
		}

		e.StartsAt, err = time.Parse(time.RFC3339, item.Start.DateTime)
		if err != nil {
			return nil, err
		}

		e.EndsAt, err = time.Parse(time.RFC3339, item.End.DateTime)
		if err != nil {
			return nil, err
		}

		// Push Notification events have a '[PN]' part in their title.
		if strings.Contains(strings.ToUpper(item.Summary), "[PN]") {
			e.Type = newspaper.CalendarEventTypePushNotification
		}

		events = append(events, e)
	}

	return events, nil
}

// CampaignsCalendar provides communication with Zalora's Campaigns calendar where
// campaign events are kept.
type CampaignsCalendar struct {
	service    *calendar.Service
	calendarID string
}

// NewCampaignCalendar initializes a new CampaignCalendar.
func NewCampaignCalendar(c *http.Client, calendarID string,
) (CampaignsCalendar, error) {

	s, err := calendar.NewService(context.Background(), option.WithHTTPClient(c))
	if err != nil {
		return CampaignsCalendar{}, err
	}

	return CampaignsCalendar{
		service:    s,
		calendarID: calendarID,
	}, nil
}

// GetCalendarEventsByDay returns a list of events from the calendar scheduled
// for the given day around the working hours (9AM-7PM SGT).
func (c CampaignsCalendar) GetCalendarEventsByDay(t time.Time,
) ([]newspaper.CalendarEvent, error) {

	start := time.Date(t.Year(), t.Month(), t.Day(), 9, 0, 0, 0, t.Location())
	end := time.Date(t.Year(), t.Month(), t.Day(), 19, 0, 0, 0, t.Location())

	resp, err := c.service.Events.
		List(c.calendarID).
		TimeMin(start.Format(time.RFC3339)).
		TimeMax(end.Format(time.RFC3339)).
		Do()
	if err != nil {
		return nil, err
	}

	var events []newspaper.CalendarEvent

	for _, item := range resp.Items {
		e := newspaper.CalendarEvent{
			Title: item.Summary,
			Type:  newspaper.CalendarEventTypeCampaign,
		}

		e.StartsAt, err = time.Parse(time.RFC3339, item.Start.DateTime)
		if err != nil {
			return nil, err
		}

		e.EndsAt, err = time.Parse(time.RFC3339, item.End.DateTime)
		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}

// DevMilestonesCalendar provides communication with Zalora's Dev Milestones calendar
// where various events for developers are kept.
type DevMilestonesCalendar struct {
	service    *calendar.Service
	calendarID string
}

// NewDevMilestonesCalendar initializes a new DevMilestonesCalendar.
func NewDevMilestonesCalendar(c *http.Client, calendarID string,
) (DevMilestonesCalendar, error) {

	s, err := calendar.NewService(context.Background(), option.WithHTTPClient(c))
	if err != nil {
		return DevMilestonesCalendar{}, err
	}

	return DevMilestonesCalendar{
		service:    s,
		calendarID: calendarID,
	}, nil
}

// GetCalendarEventsByDay returns a list of events from the calendar scheduled
// for the given day around the working hours (9AM-7PM SGT).
func (d DevMilestonesCalendar) GetCalendarEventsByDay(t time.Time,
) ([]newspaper.CalendarEvent, error) {

	start := time.Date(t.Year(), t.Month(), t.Day(), 9, 0, 0, 0, t.Location())
	end := time.Date(t.Year(), t.Month(), t.Day(), 19, 0, 0, 0, t.Location())

	resp, err := d.service.Events.
		List(d.calendarID).
		TimeMin(start.Format(time.RFC3339)).
		TimeMax(end.Format(time.RFC3339)).
		Do()
	if err != nil {
		return nil, err
	}

	var events []newspaper.CalendarEvent

	for _, item := range resp.Items {
		e := newspaper.CalendarEvent{
			Title: item.Summary,
		}

		e.StartsAt, err = time.Parse(time.RFC3339, item.Start.DateTime)
		if err != nil {
			return nil, err
		}

		e.EndsAt, err = time.Parse(time.RFC3339, item.End.DateTime)
		if err != nil {
			return nil, err
		}

		if strings.Contains(strings.ToLower(item.Summary), "code freeze") {
			e.Type = newspaper.CalendarEventTypeCodeFreeze
		}

		events = append(events, e)
	}

	return events, nil
}
