package test

import (
	"testing"
	"time"

	"github.com/Courtcircuits/HackTheCrous.api/types"
	"github.com/Courtcircuits/HackTheCrous.api/util"
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
