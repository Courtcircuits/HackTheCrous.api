package pg_util

import (
	"database/sql"
	"log"

	"github.com/Courtcircuits/HackTheCrous.api/types"
)

func QueryRestaurants(client *sql.DB, query string, args ...any) ([]types.Restaurant, error) {
	var rows *sql.Rows
	var err error
	if len(args) == 0 {
		rows, err = client.Query(query)
	} else {
		rows, err = client.Query(query, args...)
	}

	if err != nil {
		log.Fatalf("caught database err when querying : %q\n", err)
		return []types.Restaurant{}, err
	}
	defer rows.Close()
	var restaurants []types.Restaurant

	for rows.Next() {
		restaurant, err := types.ScanRestaurants(rows)
		if err != nil {
			log.Fatalf("caught database err when iterating through restaurants : %q\n", err)
			return []types.Restaurant{}, err
		}
		restaurants = append(restaurants, restaurant)
		log.Println(restaurant)
	}

	return restaurants, nil
}

// might be useless though
func QueryRestaurant(client *sql.DB, query string, args ...any) (types.Restaurant, error) {
	var row *sql.Row
	if len(args) == 0 {
		row = client.QueryRow(query)
	} else {
		row = client.QueryRow(query, args...)
	}

	return types.ScanRestaurant(row)
}
