package types

import (
	"errors"
	"time"

	"github.com/Courtcircuits/HackTheCrous.api/graph/model"
	"github.com/Courtcircuits/HackTheCrous.api/util"
	ics "github.com/arran4/golang-ical"
)

type Calendar struct {
	Url      string
	Calendar *ics.Calendar
}

func NewCalendar(url string) (Calendar, error) {
	cal, err := util.ParseCalendarFromUrl(url)
	return Calendar{
		Url:      url,
		Calendar: cal,
	}, err
}

func (cal *Calendar) ToMap() (map[time.Time][]*model.PlanningDay, error) {
	evs_timestamped := make(map[time.Time][]*model.PlanningDay)
	evs := cal.Calendar.Events()
	for _, ev := range evs {
		start, err := ev.GetStartAt()
		if err != nil {
			return nil, err
		}
		formated_date := start.Format("2006-Jan-02")
		// start at midnight
		start_date, err := time.Parse("2006-Jan-02", formated_date)
		if err != nil {
			return nil, err
		}
		ev_planningday, err := VEventToPlanningDay(*ev)
		if err != nil {
			return nil, err
		}
		evs_timestamped[start_date] = append(evs_timestamped[start_date], ev_planningday)
	}
	return evs_timestamped, nil
}

var ErrUrlIcalUndefined = errors.New("your ical ur is not defined")

func VEventToPlanningDay(vevent ics.VEvent) (*model.PlanningDay, error) {
	start, err := vevent.GetStartAt()
	if err != nil {
		return nil, err
	}
	end, err := vevent.GetEndAt()
	if err != nil {
		return nil, err
	}
	start_stringified := start.Format(time.RFC3339Nano)
	end_stringified := end.Format(time.RFC3339Nano)

	return &model.PlanningDay{
		Start:       &start_stringified,
		End:         &end_stringified,
		Summary:     &vevent.GetProperty(ics.ComponentPropertySummary).BaseProperty.Value,
		Location:    &vevent.GetProperty(ics.ComponentPropertyLocation).BaseProperty.Value,
		Description: &vevent.GetProperty(ics.ComponentPropertyDescription).BaseProperty.Value,
	}, nil
}

func (cal *Calendar) GetPeriod(start time.Time, end time.Time) ([]*model.PlanningDay, error) {
	if cal.Url == "" {
		return nil, ErrUrlIcalUndefined
	}
	evs, err := cal.ToMap()
	if err != nil {
		return nil, err
	}

	var planning_days []*model.PlanningDay

	for i := start; i.Before(end); i = i.Add(time.Hour * 24) {
		planning_days = append(planning_days, evs[i]...)
	}

	return planning_days, nil
}
