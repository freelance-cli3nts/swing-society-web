// Package gcal wraps the Google Calendar API for fetching public calendar events.
// The service account used must have the Google Calendar API enabled and must be
// granted read access to the target calendar (share the calendar with the service
// account email address).
package gcal

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"swing-society-website/server/internal/api/models"
)

// Service wraps a Google Calendar API client for a specific calendar.
type Service struct {
	svc        *calendar.Service
	calendarID string
}

// NewService initialises a Service using the same credentials as Firebase.
// Returns (nil, nil) when calendarID is empty so callers can skip GCal gracefully.
func NewService(calendarID string) (*Service, error) {
	if calendarID == "" {
		return nil, nil
	}

	ctx := context.Background()

	var (
		svc *calendar.Service
		err error
	)

	// Reuse the same credentials already used for Firebase.
	if creds := os.Getenv("GOOGLE_CREDENTIALS"); creds != "" {
		svc, err = calendar.NewService(ctx,
			option.WithCredentialsJSON([]byte(creds)),
			option.WithScopes(calendar.CalendarReadonlyScope),
		)
	} else {
		credFile := os.Getenv("FIREBASE_CREDENTIALS_FILE")
		if credFile == "" {
			credFile = "swing-society-realtime-firebase-adminsdk.json"
		}
		svc, err = calendar.NewService(ctx,
			option.WithCredentialsFile(credFile),
			option.WithScopes(calendar.CalendarReadonlyScope),
		)
	}
	if err != nil {
		return nil, err
	}

	log.Printf("Google Calendar service initialised for calendar: %s", calendarID)
	return &Service{svc: svc, calendarID: calendarID}, nil
}

// GetEvents returns all events in [from, to] from the configured calendar.
func (g *Service) GetEvents(from, to time.Time) ([]models.CalendarEvent, error) {
	result, err := g.svc.Events.List(g.calendarID).
		TimeMin(from.UTC().Format(time.RFC3339)).
		TimeMax(to.UTC().Format(time.RFC3339)).
		SingleEvents(true).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, err
	}

	events := make([]models.CalendarEvent, 0, len(result.Items))
	for _, item := range result.Items {
		events = append(events, mapEvent(item))
	}
	return events, nil
}

func mapEvent(item *calendar.Event) models.CalendarEvent {
	var start, end time.Time
	if item.Start.DateTime != "" {
		start, _ = time.Parse(time.RFC3339, item.Start.DateTime)
		end, _ = time.Parse(time.RFC3339, item.End.DateTime)
	} else {
		// All-day event — date only, no time
		start, _ = time.Parse("2006-01-02", item.Start.Date)
		end, _ = time.Parse("2006-01-02", item.End.Date)
	}
	return models.CalendarEvent{
		ID:          item.Id,
		Title:       item.Summary,
		Description: item.Description,
		Date:        start,
		EndDate:     end,
		Location:    item.Location,
		Source:      "gcal",
		URL:         item.HtmlLink,
	}
}
