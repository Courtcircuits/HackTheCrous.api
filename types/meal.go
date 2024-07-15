package types

import (
	"encoding/json"
	"time"
)

type Period string

const (
	DEJEUNER      Period = "Déjeuner"
	PETITDEJEUNER Period = "Petit déjeuner"
	DINER         Period = "Dîner"
)

type Food struct {
	Name string `json:"name"`
}

type Foods struct {
	Food []string `json:"content"`
	Type string   `json:"type"`
}

type Meal struct {
	ID      int
	Type    Period
	Foodies []Foods
	Day     time.Time
}

func ParseFoods(foodies string) ([]Foods, error) {
	var foods []Foods
	err := json.Unmarshal([]byte(foodies), &foods)
	return foods, err
}

func (f Foods) Stringify() (string, error) {
	bytesSlice, err := json.Marshal(f)
	return string(bytesSlice), err
}
