// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Attachment struct {
	Filename    *string `json:"filename,omitempty"`
	ContentType *string `json:"contentType,omitempty"`
	Size        *int    `json:"size,omitempty"`
}

type Coordinates struct {
	X *float64 `json:"x,omitempty"`
	Y *float64 `json:"y,omitempty"`
}

type Food struct {
	Names    []*string `json:"names,omitempty"`
	Category *string   `json:"category,omitempty"`
}

type Mail struct {
	From    *string   `json:"from,omitempty"`
	To      *string   `json:"to,omitempty"`
	Cc      *string   `json:"cc,omitempty"`
	Subject *string   `json:"subject,omitempty"`
	Date    *string   `json:"date,omitempty"`
	Tags    []*string `json:"tags,omitempty"`
	Text    *string   `json:"text,omitempty"`
	HTML    *string   `json:"html,omitempty"`
}

type Meal struct {
	Idmeal   *int    `json:"idmeal,omitempty"`
	Typemeal *string `json:"typemeal,omitempty"`
	Foodies  []*Food `json:"foodies,omitempty"`
	Day      *string `json:"day,omitempty"`
}

type PlanningDay struct {
	Start       *string `json:"start,omitempty"`
	End         *string `json:"end,omitempty"`
	Summary     *string `json:"summary,omitempty"`
	Location    *string `json:"location,omitempty"`
	Description *string `json:"description,omitempty"`
}

type Restaurant struct {
	Idrestaurant *int         `json:"idrestaurant,omitempty"`
	URL          *string      `json:"url,omitempty"`
	Name         *string      `json:"name,omitempty"`
	Meals        []*Meal      `json:"meals,omitempty"`
	Coords       *Coordinates `json:"coords,omitempty"`
	Distance     *float64     `json:"distance,omitempty"`
}

type School struct {
	Idschool *int         `json:"idschool,omitempty"`
	Name     *string      `json:"name,omitempty"`
	Coords   *Coordinates `json:"coords,omitempty"`
}

type User struct {
	Iduser    *int          `json:"iduser,omitempty"`
	Name      *string       `json:"name,omitempty"`
	Mail      *string       `json:"mail,omitempty"`
	Nonce     *bool         `json:"nonce,omitempty"`
	Ical      *string       `json:"ical,omitempty"`
	School    *School       `json:"school,omitempty"`
	Favorites []*Restaurant `json:"favorites,omitempty"`
}
