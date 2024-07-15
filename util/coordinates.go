package util

import (
	"math"
	"strconv"
	"strings"
)

type Coordinates struct {
	X float64
	Y float64
}

func ParseCoordinates(coordinates_string string) (Coordinates, error) {
	coordinates_arr := strings.Split(coordinates_string[1:len(coordinates_string)-1], ",")

	x, err := strconv.ParseFloat(coordinates_arr[0], 32)

	if err != nil {
		return Coordinates{}, err
	}

	y, err := strconv.ParseFloat(coordinates_arr[1], 32)

	if err != nil {
		return Coordinates{}, err
	}

	return Coordinates{
		X: x,
		Y: y,
	}, nil
}

func ComputeDistance(a, b Coordinates) float64 {
	R := 6371e3
	phi1 := a.X * math.Pi / 180
	phi2 := a.Y * math.Pi / 180
	deltaPhi := (a.X - b.X) * math.Pi / 180
	deltaLambda := (a.Y - b.Y) * math.Pi / 180
	aPrim := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)

	c := 2 * math.Atan2(math.Sqrt(aPrim), math.Sqrt(1-aPrim))

	return R * c
}
