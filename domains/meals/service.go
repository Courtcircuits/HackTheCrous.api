package meals

import (
	"database/sql"

	"github.com/Courtcircuits/HackTheCrous.api/storage"
	"github.com/Courtcircuits/HackTheCrous.api/types"
)

type IMealService interface {
	GetMeals(restaurant_id string) ([]types.Meal, error)
}

type MealServiceImpl struct {
	database storage.IDatabase
}

type MealSQL struct {
	ID       sql.NullInt32  `json:"ID,omitempty"`
	Typemeal sql.NullString `json:"Typemeal,omitempty"`
	Foodies  sql.NullString `json:"Foodies,omitempty"`
	Day      sql.NullTime   `json:"Day,omitempty"`
}

func (service MealServiceImpl) GetMeals(restaurant_id string) ([]types.Meal, error) {
	query := "SELECT idmeal, typemeal, foodies, day FROM meal WHERE idrestaurant = $1"

	rows, err := service.database.Select(query, restaurant_id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meals := make([]types.Meal, 0)
	for rows.Next() {
		meal, err := scanMeal(rows.Scan)
		if err != nil {
			return nil, err
		}
		meals = append(meals, meal)
	}

	return meals, nil
}

func scanMeal(scan func(args ...any) error) (types.Meal, error) {
	var meal MealSQL
	err := scan(&meal.ID, &meal.Typemeal, &meal.Foodies, &meal.Day)
	if err != nil {
		return types.Meal{}, err
	}

	foods, err := types.ParseFoods(meal.Foodies.String)

	if err != nil {
		return types.Meal{}, err
	}

	var period types.Period

	switch meal.Typemeal.String {
	case "Déjeuner":
		period = types.DEJEUNER
	case "Petit déjeuner":
		period = types.PETITDEJEUNER
	default:
		period = types.DINER
	}

	return types.Meal{
		ID:      int(meal.ID.Int32),
		Type:    period,
		Foodies: foods,
		Day:     meal.Day.Time,
	}, nil

}

var MealService IMealService = MealServiceImpl{
	database: storage.Database,
}
