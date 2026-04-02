package models

import "time"

// CalendarEvent represents a single class, party, workshop, or festival event.
type CalendarEvent struct {
	ID          string
	Title       string
	Description string
	Date        time.Time
	EndDate     time.Time
	Location    string
	Type        string // "class" | "party" | "workshop" | "festival"
	Source      string // "firebase" | "gcal"
	URL         string
}

// Category classifies the event relative to now: "past", "today", or "upcoming".
func (e CalendarEvent) Category(now time.Time) string {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	eventDay := time.Date(e.Date.Year(), e.Date.Month(), e.Date.Day(), 0, 0, 0, 0, e.Date.Location())
	switch {
	case eventDay.Before(today):
		return "past"
	case eventDay.Equal(today):
		return "today"
	default:
		return "upcoming"
	}
}

// DedupKey returns a string used to detect duplicate events across sources.
func (e CalendarEvent) DedupKey() string {
	return e.Date.Format("2006-01-02") + "|" + e.Title
}
