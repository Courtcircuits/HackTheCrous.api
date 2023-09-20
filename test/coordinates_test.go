package test

import (
	"math"
	"testing"

	"github.com/Courtcircuits/HackTheCrous.api/util"
)

func EqualsFloat(a, b, delta float64) bool {
	return math.Abs(a-b)-delta <= 0
}

func TestCoordinatesParsing(t *testing.T) {
	coords, err := util.Parse_coordinates("(43.6349531,3.870764)")

	if err != nil {
		t.Fatalf("shouldn't throw error %q\n", err)
	}

	if !EqualsFloat(coords.X, 43.6349531, 0.00001) {
		t.Fatalf("expected %f got %f\n", float64(43.6349531), coords.X)
	}

	if !EqualsFloat(coords.Y, 3.870764, 0.00001) {
		t.Fatalf("expected %f got %f\n", float64(3.870764), coords.Y)
	}
}
