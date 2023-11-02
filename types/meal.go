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

type Food struct {
	Name string `json:"name"`
}

type Foods struct {
	Food []string `json:"food"`
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

func (f Foods) ToGraphQL() *model.Food {
	names := make([]*string, len(f.Food))
	for i := 0; i < len(f.Food); i++ {
		names[i] = &f.Food[i]
	}
	return &model.Food{
		Names:    names,
		Category: &f.Type,
	}
}

func (m Meal) ToGraphQL() *model.Meal {
	stringifiedDate := m.Day.Format(time.RFC3339Nano)

	foodies_gql := make([]*model.Food, len(m.Foodies))

	for i := 0; i < len(m.Foodies); i++ {
		foodies_gql = append(foodies_gql, m.Foodies[i].ToGraphQL())
	}

	return &model.Meal{
		Idmeal:   &m.ID,
		Typemeal: (*string)(&m.Type),
		Foodies:  foodies_gql,
		Day:      &stringifiedDate,
	}
}

func (f Food) ToGraphQL() *model.Food {
	fake_cat := "food"
	return &model.Food{
		Names:    []*string{&f.Name},
		Category: &fake_cat,
	}
}
