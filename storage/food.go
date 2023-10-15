package storage

import (
	"log"

	"github.com/Courtcircuits/HackTheCrous.api/storage/pg_util"
	"github.com/Courtcircuits/HackTheCrous.api/types"
)

func (db *PostgresDatabase) GetFoods() ([]types.Food, error) {
	query := `SELECT DISTINCT sr.keyword from suggestions_restaurant sr JOIN cat_suggestions cs ON sr.idcat = cs.idcat WHERE cs.namecat='food' LIMIT 12`
	client, err := db.Connect()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	rows, err := client.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foods []types.Food
	for rows.Next() {
		var food types.Food
		err := rows.Scan(&food.Name)
		if err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}
	return foods, nil
}

func (db *PostgresDatabase) GetRestaurantsFromFood(food string) ([]types.Restaurant, error) {
	query := `SELECT r.idrestaurant, r.name, r.url, r.gpscoord FROM restaurant r JOIN suggestions_restaurant sr ON r.idrestaurant = sr.idrestaurant WHERE sr.keyword=$1`
	client, err := db.Connect()
	if err != nil {
		log.Fatal(err)
		return []types.Restaurant{}, err
	}
	defer client.Close()

	return pg_util.QueryRestaurants(client, query, food)
}
