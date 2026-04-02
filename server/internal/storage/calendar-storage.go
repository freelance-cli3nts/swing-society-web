package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"swing-society-website/server/internal/api/models"
)

// CalendarStorage handles reading and writing calendar events in Firebase.
type CalendarStorage struct {
	firebase *FirebaseClient
}

// NewCalendarStorage creates a CalendarStorage, degrading gracefully if Firebase is unavailable.
func NewCalendarStorage() *CalendarStorage {
	firebase, err := NewFirebaseClient()
	if err != nil {
		log.Printf("Warning: CalendarStorage Firebase init failed: %v", err)
	}
	return &CalendarStorage{firebase: firebase}
}

// firebaseEvent is the raw JSON shape stored under events/{id} in Firebase.
type firebaseEvent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`     // RFC3339
	EndDate     string `json:"end_date"` // RFC3339
	Location    string `json:"location"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

// GetEvents fetches all events from Firebase that fall within [from, to].
func (s *CalendarStorage) GetEvents(from, to time.Time) ([]models.CalendarEvent, error) {
	if s.firebase == nil {
		return nil, nil
	}

	ctx := context.Background()
	ref := s.firebase.db.NewRef("events")

	var raw map[string]firebaseEvent
	if err := ref.Get(ctx, &raw); err != nil {
		return nil, fmt.Errorf("firebase get events: %w", err)
	}

	var events []models.CalendarEvent
	for id, fe := range raw {
		date, err := time.Parse(time.RFC3339, fe.Date)
		if err != nil {
			log.Printf("CalendarStorage: skipping event %s, bad date %q: %v", id, fe.Date, err)
			continue
		}
		if date.Before(from) || date.After(to) {
			continue
		}
		end, _ := time.Parse(time.RFC3339, fe.EndDate)
		events = append(events, models.CalendarEvent{
			ID:          id,
			Title:       fe.Title,
			Description: fe.Description,
			Date:        date,
			EndDate:     end,
			Location:    fe.Location,
			Type:        fe.Type,
			Source:      "firebase",
			URL:         fe.URL,
		})
	}
	return events, nil
}

// SaveEvent writes a new event to Firebase and returns the generated key.
func (s *CalendarStorage) SaveEvent(event models.CalendarEvent) (string, error) {
	if s.firebase == nil {
		return "", fmt.Errorf("firebase not initialized")
	}

	ctx := context.Background()
	ref := s.firebase.db.NewRef("events")

	data := map[string]interface{}{
		"title":       event.Title,
		"description": event.Description,
		"date":        event.Date.Format(time.RFC3339),
		"end_date":    event.EndDate.Format(time.RFC3339),
		"location":    event.Location,
		"type":        event.Type,
		"url":         event.URL,
	}

	newRef, err := ref.Push(ctx, data)
	if err != nil {
		return "", fmt.Errorf("firebase push event: %w", err)
	}
	return newRef.Key, nil
}
