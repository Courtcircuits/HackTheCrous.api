package test

import (
	"testing"
	"time"

	"github.com/Courtcircuits/HackTheCrous.api/types"
	"github.com/Courtcircuits/HackTheCrous.api/util"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestGetPeriod(t *testing.T) {
	cal, err := types.NewCalendar(util.Get("CALENDAR_LINK"))
	if err != nil {
		t.Fatalf("shouldn't throw %v\n", err)
	}
	today, _ := time.Parse("2006-Jan-02", "2023-Sep-17")
	inAWeek, _ := time.Parse("2006-Jan-02", "2023-Sep-22")
	period, err := cal.GetPeriod(today, inAWeek)
	if err != nil {
		t.Fatalf("shouldn't throw %v\n", err)
	}
	if len(period) != 8 {
		t.Fatalf("period should contain %d, but contains %d\n", 8, len(period))
	}
}

func TestSerializationMap(t *testing.T) {
	whatever := "whatever"
	now := time.Now()
	test_map := types.JsonCalendar{
		now: {{
			Start:       &whatever,
			End:         &whatever,
			Summary:     &whatever,
			Location:    &whatever,
			Description: &whatever,
		}},
	}

	val, err := types.JsonCalendarToString(&test_map)
	if err != nil {
		t.Errorf("caught an unexpected error : %v\n", err)
	}
	if val == "" {
		t.Errorf("returned empty string")
	}

	map_parsed, err := types.ParseJsonCalendar(val)

	if err != nil {
		t.Errorf("caught an unexpected error : %v\n", err)
	}

	test_map_day := *maps.Values(test_map)[0][0]
	map_parsed_day := *maps.Values(map_parsed)[0][0]

	assert.Equal(t, test_map_day, map_parsed_day)
}
