package yaas

import "strings"

// Browser is a type that represents browser identification
type Browser string

var (
	browserMap = map[string]Browser{
		"Chrome":   Chrome,
		"Chromium": Chromium,
		"Firefox":  Firefox,
		"MSIE":     InternetExplorer,
		"Opera":    Opera,
		"Safari":   Safari,
	}
)

// ParseBrowser finds the browser name in the string and returns the matchin type.
func ParseBrowser(v string) Browser {
	for name, browser := range browserMap {
		if strings.Contains(v, name) {
			return browser
		}
	}

	return OtherBrowser
}

const (
	// Chrome represents Google Chrome browser
	Chrome Browser = "Chrome"
	// Chromium represents Chromium browser
	Chromium Browser = "Chromium"
	// Firefox represents Firefox browser
	Firefox Browser = "Firefox"
	// InternetExplorer represents IE browser
	InternetExplorer = "Internet Explorer"
	// Opera represents Opera browser
	Opera Browser = "Opera"
	// Safari represents Safari browser
	Safari Browser = "Safari"
	// OtherBrowser represents other browsers
	OtherBrowser Browser = "Other"
)
