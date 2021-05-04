package yaas

import (
	"regexp"
	"time"
)

// RawEvent stores information from the event api request
type RawEvent struct {
	Code     string `json:"c"`
	URL      string `json:"u"`
	Referrer string `json:"r"`
	Width    int    `json:"w"`
}

// GetEvent returns internal event type that has all data converted and parsed
func (e RawEvent) GetEvent() Event {
	var ye Event

	ye.Code = e.Code
	// Get current timestamp
	ye.Time = time.Now()
	// Get requiested url
	ye.URL = e.URL
	ye.Path = findRequestPath(e.URL)
	// Get request referrer
	ye.Referrer = e.Referrer
	// Get device type
	ye.Device = specifyDeviceType(e.Width)

	return ye
}

func findRequestPath(u string) string {
	// Regexp for scraping path from url
	re, err := regexp.Compile("(http|https)://.*?(/.*)")
	if err != nil {
		panic(err)
	}
	// Get path from url
	m := re.FindStringSubmatch(u)
	if len(m) != 3 {
		return "/"
		// panic(errors.New("event: cant find path in url"))
	}
	return m[2]
}

func specifyDeviceType(w int) Device {
	if w >= 1281 {
		return Desktop
	} else if w >= 1025 && w <= 1280 {
		return Laptop
	} else if w >= 481 && w <= 1024 {
		return Tablet
	} else if w >= 320 && w <= 480 {
		return Mobile
	} else {
		return Unkown
	}
}
