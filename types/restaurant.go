package types

import (
	"github.com/Courtcircuits/HackTheCrous.api/util"
)

type Restaurant struct {
	ID        int              `json:"id,omitempty"`
	Url       string           `json:"url,omitempty"`
	Name      string           `json:"name,omitempty"`
	Gps_coord util.Coordinates `json:"gps_coord,omitempty"`
	Hours     string           `json:"hours,omitempty"`
}
