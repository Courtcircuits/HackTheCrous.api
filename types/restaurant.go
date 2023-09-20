package types

import (
	"database/sql"

	"github.com/Courtcircuits/HackTheCrous.api/graph/model"
	"github.com/Courtcircuits/HackTheCrous.api/util"
)

type Restaurant struct {
	ID        int              `json:"id,omitempty"`
	Url       string           `json:"url,omitempty"`
	Name      string           `json:"name,omitempty"`
	Gps_coord util.Coordinates `json:"gps_coord,omitempty"`
}

func ScanRestaurants(row *sql.Rows) (Restaurant, error) {
	type SQL_restaurant struct {
		Idrestaurant sql.NullInt32  `json:"idrestaurant,omitempty"`
		Url          sql.NullString `json:"url,omitempty"`
		Name         sql.NullString `json:"name,omitempty"`
		Gpscoord     sql.NullString `json:"gpscoord,omitempty"`
	}

	var sql_restaurant SQL_restaurant
	err := row.Scan(&sql_restaurant.Idrestaurant, &sql_restaurant.Name, &sql_restaurant.Url, &sql_restaurant.Gpscoord)

	if err != nil {
		return Restaurant{}, err
	}

	coords, err := util.Parse_coordinates(sql_restaurant.Gpscoord.String)

	return Restaurant{
		ID:        int(sql_restaurant.Idrestaurant.Int32),
		Name:      sql_restaurant.Name.String,
		Url:       sql_restaurant.Url.String,
		Gps_coord: coords,
	}, err
}

func (r Restaurant) ToGraphQL() *model.Restaurant {
	return &model.Restaurant{
		Idrestaurant: &r.ID,
		URL:          &r.Url,
		Name:         &r.Name,
		Coords:       r.Gps_coord.ToGraphQL(),
	}
}
