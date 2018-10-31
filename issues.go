package triage

import "time"

// Issue represents an issue to be triaged.
type Issue struct {
	Title    string
	Link     string
	User     string
	Comments int
	Created  time.Time
}
