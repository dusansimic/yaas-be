package yaas

import "time"

// Event stores info about a single event
type Event struct {
	Time        time.Time `json:"timestamp" db:"timestamp"`
	URL         string    `json:"url" db:"url"`
	Code        string    `json:"code"`
	Path        string    `json:"path" db:"path"`
	Referrer    string    `json:"referrer" db:"referrer"`
	Device      Device    `json:"device" db:"device"`
	CountryCode string    `json:"country_code" db:"country_code"`
	Browser     Browser   `json:"browser" db:"browser"`
	OS          OS        `json:"os" db:"os"`
}
