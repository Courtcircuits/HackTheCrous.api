package restaurants

import (
	"database/sql"
	"log"
	"strings"

	"github.com/Courtcircuits/HackTheCrous.api/domains/meals"
	"github.com/Courtcircuits/HackTheCrous.api/storage"
	"github.com/Courtcircuits/HackTheCrous.api/types"
	"github.com/Courtcircuits/HackTheCrous.api/util"
)

type RestaurantSQL struct {
	ID       sql.NullInt64  `json:"id,omitempty"`
	Url      sql.NullString `json:"url,omitempty"`
	Name     sql.NullString `json:"name,omitempty"`
	GpsCoord sql.NullString `json:"gps_coord,omitempty"`
	Hours    sql.NullString `json:"hours,omitempty"`
}

type IRestaurantService interface {
	GetRestaurants() ([]types.Restaurant, error)
	GetRestaurant(id int) (types.Restaurant, error)
	GetMeals(id string) ([]types.Meal, error)
	Search(query string) ([]types.Restaurant, error)
}

type RestaurantServiceImpl struct {
	mealService meals.IMealService
	database    storage.IDatabase
}

func (service RestaurantServiceImpl) GetRestaurants() ([]types.Restaurant, error) {
	query := "SELECT idrestaurant, url, name, gpscoord, hours FROM restaurant"
	rows, err := service.database.Select(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	restaurants := make([]types.Restaurant, 0)
	for rows.Next() {
		restaurant, err := scanRestaurant(rows.Scan)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

func (service RestaurantServiceImpl) GetRestaurant(id int) (types.Restaurant, error) {
	query := "SELECT idrestaurant, url, name, gpscoord, hours FROM restaurant WHERE idrestaurant = $1"

	rows, err := service.database.Select(query, id)
	if err != nil {
		return types.Restaurant{}, err
	}
	defer rows.Close()

	if rows.Next() {
		restaurant, err := scanRestaurant(rows.Scan)
		if err != nil {
			return types.Restaurant{}, err
		}
		return restaurant, nil
	}

	return types.Restaurant{}, nil
}

func (service RestaurantServiceImpl) GetMeals(restaurant_id string) ([]types.Meal, error) {
	return service.mealService.GetMeals(restaurant_id)
}

func (service RestaurantServiceImpl) Search(query string) ([]types.Restaurant, error) {
	sql_query := `SELECT idrestaurant, url, name, gpscoord, hours FROM restaurant
WHERE idrestaurant IN (SELECT r.idrestaurant FROM restaurant r JOIN suggestions_restaurant sr ON sr.idrestaurant=r.idrestaurant WHERE UPPER(sr.keyword) LIKE $1)`

	query = "%" + strings.ToUpper(query) + "%"

	rows, err := service.database.Select(sql_query, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	restaurants := make([]types.Restaurant, 0)
	for rows.Next() {
		restaurant, err := scanRestaurant(rows.Scan)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

func scanRestaurant(scan func(dest ...any) error) (types.Restaurant, error) {
	var restaurant RestaurantSQL
	err := scan(&restaurant.ID, &restaurant.Url, &restaurant.Name, &restaurant.GpsCoord, &restaurant.Hours)
	if err != nil {
		return types.Restaurant{}, err
	}
	coords, _ := util.ParseCoordinates(restaurant.GpsCoord.String)
	return types.Restaurant{
		ID:        int(restaurant.ID.Int64),
		Url:       restaurant.Url.String,
		Name:      restaurant.Name.String,
		Gps_coord: coords,
		Hours:     restaurant.Hours.String,
	}, nil
}

var RestaurantService IRestaurantService = RestaurantServiceImpl{
	mealService: meals.MealService,
	database:    storage.Database,
}
