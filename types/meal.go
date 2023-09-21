package types

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/Courtcircuits/HackTheCrous.api/graph/model"
)

type Period string

const (
	DEJEUNER      Period = "Déjeuner"
	PETITDEJEUNER Period = "Petit déjeuner"
	DINER         Period = "Dîner"
)

type Foods []struct {
	Food []string `json:"food"`
	Type string   `json:"type"`
}

type Meal struct {
	ID      int
	Type    Period
	Foodies Foods
	Day     time.Time
}

func ParseFoods(foodies string) (Foods, error) {
	var foods Foods
	err := json.Unmarshal([]byte(foodies), &foods)
	return foods, err
}

func (f Foods) Stringify() (string, error) {
	bytesSlice, err := json.Marshal(f)
	return string(bytesSlice), err
}

func ScanMeals(rows *sql.Rows) (Meal, error) {
	type SQL_Meals struct {
		Idmeal   sql.NullInt32  `json:"idmeal,omitempty"`
		Typemeal sql.NullString `json:"typemeal,omitempty"`
		Foodies  sql.NullString `json:"foodies,omitempty"`
		Day      sql.NullTime   `json:"day,omitempty"`
	}

	var sql_meals SQL_Meals
	err := rows.Scan(&sql_meals.Idmeal, &sql_meals.Typemeal, &sql_meals.Foodies, &sql_meals.Day)

	if err != nil {
		return Meal{}, err
	}

	foods, err := ParseFoods(sql_meals.Foodies.String)

	if err != nil {
		log.Println(err.Error())
	}

	var period Period

	switch sql_meals.Typemeal.String {
	case "Déjeuner":
		period = DEJEUNER
	case "Petit déjeuner":
		period = PETITDEJEUNER
	default:
		period = DINER
	}

	return Meal{
		ID:      int(sql_meals.Idmeal.Int32),
		Type:    period,
		Foodies: foods,
		Day:     sql_meals.Day.Time,
	}, err
}

func (m Meal) ToGraphQL() *model.Meal {
	stringifiedFoodies, err := m.Foodies.Stringify()
	if err != nil {
		log.Fatal(err)
	}
	stringifiedDate := m.Day.Format(time.RFC3339Nano)

	return &model.Meal{
		Idmeal:   &m.ID,
		Typemeal: (*string)(&m.Type),
		Foodies:  &stringifiedFoodies,
		Day:      &stringifiedDate,
	}
}
