package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"net/http"
	"strings"
	"time"
)

const (
	alertsURL    = "http://portal.bisbank.com.ua/kyiv_alerts.json"
	pollInterval = time.Minute
	httpTimeout  = 15 * time.Second

	// true: play a sound on every check (every minute).
	// false: play a sound only when the status or level changes.
	playSoundEveryCheck = false
)

// Icon colors.
var (
	colGreen  = color.RGBA{R: 40, G: 180, B: 60, A: 255}
	colYellow = color.RGBA{R: 255, G: 200, B: 0, A: 255}
	colOrange = color.RGBA{R: 255, G: 130, B: 0, A: 255}
	colRed    = color.RGBA{R: 220, G: 30, B: 30, A: 255}
	colPurple = color.RGBA{R: 150, G: 40, B: 200, A: 255}
	colGray   = color.RGBA{R: 140, G: 140, B: 140, A: 255} // fetch error
)

// levelColors maps the "level" field to an icon color.
// Add entries here if the server starts sending other values.
var levelColors = map[string]color.RGBA{
	"green":  colGreen,
	"yellow": colYellow,
	"orange": colOrange,
	"red":    colRed,
	"purple": colPurple,
}

// colorForUnknownLevel is used when "level" is empty or not in the table.
var colorForUnknownLevel = colRed

// alertStatus matches the server JSON, for example:
//
//	{"slug":"kyiv","type":"state","name_uk":"м. Київ","name_en":"Kyiv",
//	 "active":true,"level":"yellow",
//	 "started_at":"2026-11-02T15:09:07.000000Z",
//	 "updated_at":"2026-11-02T15:42:50.705565Z",
//	 "districts_active":[],"active_anywhere":true}
type alertStatus struct {
	Slug            string    `json:"slug"`
	Type            string    `json:"type"`
	NameUK          string    `json:"name_uk"`
	NameEN          string    `json:"name_en"`
	Active          bool      `json:"active"`
	Level           string    `json:"level"`
	StartedAt       time.Time `json:"started_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DistrictsActive []string  `json:"districts_active"`
	ActiveAnywhere  bool      `json:"active_anywhere"`
}

func (s alertStatus) level() string {
	return strings.ToLower(strings.TrimSpace(s.Level))
}

// color returns the icon color: green when there is no alert,
// otherwise the color for the alert level.
func (s alertStatus) color() color.RGBA {
	if !s.Active {
		return colGreen
	}
	if c, ok := levelColors[s.level()]; ok {
		return c
	}
	return colorForUnknownLevel
}

// key identifies the state, used to detect changes between checks.
func (s alertStatus) key() string {
	if !s.Active {
		return "off"
	}
	return "on:" + s.level()
}

// name returns the region name, preferring Ukrainian.
func (s alertStatus) name() string {
	if s.NameUK != "" {
		return s.NameUK
	}
	if s.NameEN != "" {
		return s.NameEN
	}
	return "Kyiv"
}

// tooltip is the hover text for the tray icon (Windows limit: 127 chars).
func (s alertStatus) tooltip(checked time.Time) string {
	if s.Active {
		since := ""
		if !s.StartedAt.IsZero() {
			since = " з " + s.StartedAt.Local().Format("15:04")
		}
		return fmt.Sprintf("%s: ТРИВОГА (%s)%s\nПеревірено %s",
			s.name(), s.level(), since, checked.Format("15:04:05"))
	}
	return fmt.Sprintf("%s: тривоги немає\nПеревірено %s",
		s.name(), checked.Format("15:04:05"))
}

var httpClient = &http.Client{Timeout: httpTimeout}

// fetchAlert downloads and parses the alert JSON.
func fetchAlert() (alertStatus, error) {
	var st alertStatus

	req, err := http.NewRequest(http.MethodGet, alertsURL, nil)
	if err != nil {
		return st, err
	}
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := httpClient.Do(req)
	if err != nil {
		return st, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return st, fmt.Errorf("HTTP %s", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&st); err != nil {
		return st, fmt.Errorf("bad JSON: %w", err)
	}
	return st, nil
}
