package yaas

type Refers struct {
	Referrer string `json:"name"`
	Reqs     int    `json:"reqs"`
}

type Devices struct {
	Device Device `json:"name"`
	Reqs   int    `json:"reqs"`
}

type Paths struct {
	Path string `json:"name"`
	Reqs int    `json:"reqs"`
}

type OSs struct {
	OS   OS  `json:"name"` // TODO: should use a separate type
	Reqs int `json:"reqs"`
}

type Browsers struct {
	Browser Browser `json:"name"` // TODO: should use a separate type
	Reqs    int     `json:"reqs"`
}

type Record struct {
	Name string `json:"name"` // Name of the specificy group in entity
	Reqs int    `json:"reqs"` // Number of requests for specified group
}

// Summary stores data of a single summary
type Summary struct {
	Domain    string   `json:"domain"`
	Type      string   `json:"type"`
	Referrers []Record `json:"referrers"`
	Devices   []Record `json:"devices"`
	Paths     []Record `json:"paths"`
	OS        []Record `json:"os"`
	Browsers  []Record `json:"browsers"`
}
