package handlers

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"swing-society-website/server/internal/api/models"
	"swing-society-website/server/internal/gcal"
	"swing-society-website/server/internal/storage"
)

// CalendarHandler serves past / today / upcoming event fragments for HTMX.
type CalendarHandler struct {
	calStorage *storage.CalendarStorage
	gcalSvc    *gcal.Service // nil when Google Calendar is not configured
}

func NewCalendarHandler(calStorage *storage.CalendarStorage, gcalSvc *gcal.Service) *CalendarHandler {
	return &CalendarHandler{calStorage: calStorage, gcalSvc: gcalSvc}
}

// eventListTmpl renders a list of CalendarEvents as an HTMX HTML fragment.
var eventListTmpl = template.Must(template.New("event-list").Funcs(template.FuncMap{
	"bgMonth": func(t time.Time) string {
		months := [13]string{"", "Яну", "Фев", "Мар", "Апр", "Май", "Юни", "Юли", "Авг", "Сеп", "Окт", "Ное", "Дек"}
		return months[t.Month()]
	},
	"typeLabel": func(typ string) string {
		switch typ {
		case "class":
			return "Клас"
		case "party":
			return "Парти"
		case "workshop":
			return "Уъркшоп"
		case "festival":
			return "Фестивал"
		default:
			return typ
		}
	},
	"typeBadge": func(typ string) string {
		switch typ {
		case "class":
			return "is-info"
		case "party":
			return "is-warning"
		case "workshop":
			return "is-success"
		case "festival":
			return "is-danger"
		default:
			return "is-light"
		}
	},
	"hasTime": func(t time.Time) bool {
		return t.Hour() != 0 || t.Minute() != 0
	},
	"fmtTime": func(t time.Time) string {
		return t.Format("15:04")
	},
}).Parse(`
{{if .}}
{{range .}}
<div class="event-card event-type-{{.Type}}">
  <div class="event-date-badge">
    <span class="event-day">{{.Date.Day}}</span>
    <span class="event-month">{{bgMonth .Date}}</span>
  </div>
  <div class="event-info">
    {{if .Type}}<span class="tag {{typeBadge .Type}} is-small mb-1">{{typeLabel .Type}}</span>{{end}}
    <p class="event-title">{{.Title}}</p>
    {{if .Location}}<p class="event-location">📍 {{.Location}}</p>{{end}}
    {{if hasTime .Date}}<p class="event-time">🕐 {{fmtTime .Date}}</p>{{end}}
  </div>
  {{if .URL}}<a href="{{.URL}}" target="_blank" rel="noopener" class="button is-small is-light event-link">Детайли</a>{{end}}
</div>
{{end}}
{{else}}
<p class="event-empty">Няма записи</p>
{{end}}
`))

// HandleCalendar serves GET /api/calendar/{past|today|upcoming}
func (h *CalendarHandler) HandleCalendar(w http.ResponseWriter, r *http.Request) {
	category := strings.TrimPrefix(r.URL.Path, "/api/calendar/")
	category = strings.TrimSuffix(category, "/")

	if category != "past" && category != "today" && category != "upcoming" {
		http.NotFound(w, r)
		return
	}

	now := time.Now()
	from := now.AddDate(0, 0, -30)
	to := now.AddDate(0, 0, 60)

	events := h.fetchAll(from, to)

	var filtered []models.CalendarEvent
	for _, e := range events {
		if e.Category(now) == category {
			filtered = append(filtered, e)
		}
	}

	// Past events: most recent first. Today/upcoming: soonest first.
	if category == "past" {
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].Date.After(filtered[j].Date)
		})
	} else {
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].Date.Before(filtered[j].Date)
		})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := eventListTmpl.Execute(w, filtered); err != nil {
		log.Printf("CalendarHandler: template error: %v", err)
	}
}

// fetchAll pulls events from Firebase and Google Calendar, deduplicating by date+title.
func (h *CalendarHandler) fetchAll(from, to time.Time) []models.CalendarEvent {
	seen := make(map[string]bool)
	var events []models.CalendarEvent

	if fbEvents, err := h.calStorage.GetEvents(from, to); err != nil {
		log.Printf("CalendarHandler: Firebase fetch error: %v", err)
	} else {
		for _, e := range fbEvents {
			k := e.DedupKey()
			if !seen[k] {
				seen[k] = true
				events = append(events, e)
			}
		}
	}

	if h.gcalSvc != nil {
		if gcalEvents, err := h.gcalSvc.GetEvents(from, to); err != nil {
			log.Printf("CalendarHandler: GCal fetch error: %v", err)
		} else {
			for _, e := range gcalEvents {
				k := e.DedupKey()
				if !seen[k] {
					seen[k] = true
					events = append(events, e)
				}
			}
		}
	}

	return events
}
