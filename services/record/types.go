package record

import "time"

type Record struct {
	DomainID   int
	ReferrerID int
	OSID       int
	DeviceID   int
	CountryID  int
	BrowserID  int
	Timestamp  time.Time
	Url        string
	Path       string
}
