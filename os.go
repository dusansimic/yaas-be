package yaas

import "strings"

// OS is a type that represents an operating system
type OS string

var (
	osMap = map[string]OS{
		"Android":   Android,
		"Linux":     Linux,
		"Mac OS X":  MacOS,
		"Windows":   Windows,
		"iPhone OS": iOS,
	}
)

func GetOS(v string) OS {
	for name, os := range osMap {
		if strings.Contains(v, name) {
			return os
		}
	}

	return OtherOS
}

const (
	// Android represents Android OS
	Android OS = "Android"
	// Linux represents Linux based OS
	Linux OS = "Linux"
	// MacOS represents Mac OS X
	MacOS OS = "macOS"
	// Windows represents Windows os
	Windows OS = "Windows"
	// iOS represents iOS
	iOS OS = "iOS"
	// OtherOS represents other operation systems
	OtherOS OS = "Other"
)
