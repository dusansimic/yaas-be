package stats

import "time"

type Record struct {
	Name string `json:"name"` // Name of the specificy group in entity
	Reqs int    `json:"reqs"` // Number of requests for specified group
}

type TimeRecord struct {
	Time time.Time `json:"time"`
	Reqs int       `json:"reqs"`
}

// Summary stores records of a single summary
type Summary struct {
	Domain    string   `json:"domain"`
	Type      string   `json:"type"`
	Referrers []Record `json:"referrers"`
	Devices   []Record `json:"devices"`
	Paths     []Record `json:"paths"`
	OS        []Record `json:"os"`
	Browsers  []Record `json:"browsers"`
}
