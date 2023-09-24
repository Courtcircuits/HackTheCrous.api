package util

import (
	"net/http"

	ics "github.com/arran4/golang-ical"
)

func ParseCalendarFromUrl(url string) (*ics.Calendar, error) {
	http_client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http_client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ics.ParseCalendar(resp.Body)
}
