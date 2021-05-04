package yaas

// Device is a type that represents device identificator
type Device string

const (
	// Desktop represents a desktop device
	Desktop Device = "Desktop"
	// Laptop represents a laptop device
	Laptop Device = "Laptop"
	// Tablet represents a tablet device
	Tablet Device = "Tablet"
	// Mobile represents a mobile device
	Mobile Device = "Mobile"
	// Unkown represents an unkown device
	Unkown Device = "Other"
)

func ParseDevice(w int) Device {
	if w > 1224 {
		return Desktop
	} else if w > 1024 {
		return Laptop
	} else if w > 768 {
		return Tablet
	} else if w > 320 {
		return Mobile
	}
	return Unkown
}
