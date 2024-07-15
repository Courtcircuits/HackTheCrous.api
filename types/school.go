package types

import (
	"github.com/Courtcircuits/HackTheCrous.api/util"
)

type School struct {
	ID     int              `json:"id,omitempty"`
	Name   string           `json:"name,omitempty"`
	Coords util.Coordinates `json:"coords,omitempty"`
}
