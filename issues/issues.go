package issues

import (
	"strconv"
	"strings"
	"time"
)

// Issue represents an issue to be triaged.
type Issue struct {
	Title    string
	Link     string
	User     string
	Comments int
	Labels   []Label
	Created  time.Time
}

// Label represents an issue label.
type Label struct {
	Name   string
	Colour string
}

// TextColour returns a hex colour to use for a label name, either black or white.
// The label colour is used to determine the best high-contrast value.
func (l Label) TextColour() string {
	colour := "#000000"
	h := strings.TrimLeft(l.Colour, "#")

	if len(h) == 6 && (hex(h[0:2])*0.299+hex(h[2:4])*0.587+hex(h[4:6])*0.114) < 127 {
		colour = "#ffffff"
	}

	return colour
}

func hex(hex string) float32 {
	n, _ := strconv.ParseUint(hex, 16, 32)
	return float32(n)
}
