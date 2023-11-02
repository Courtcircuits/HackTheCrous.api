package test

import (
	"testing"

	"github.com/Courtcircuits/HackTheCrous.api/storage"
	"github.com/Courtcircuits/HackTheCrous.api/types"
)

func TestCreateUser(t *testing.T) {
	email := "testtest@etu.umontpellier.fr"
	password := "12341234"
	pg_storage := storage.NewPostgresDatabase()

	err_delete := pg_storage.DeleteUserByMail(email)

	if err_delete != nil {
		t.Errorf("error when delete : %q", err_delete)
	}

	user, err := pg_storage.CreateLocalUser(email, password)

	if err != nil {
		t.Errorf("shouldn't throw error %q", err)
		return
	}

	if email != user.Email.String {
		t.Errorf("got %q different than expected %q", user.Email.String, email)
		return
	}

	user_searched, err := pg_storage.GetUserByEmail(email)

	if err != nil {
		t.Errorf("shouldn't throw error %q", err)
		return
	}

	if user_searched.ID.Int32 != user.ID.Int32 {
		t.Errorf("got %d different ID than expected %d", user_searched.ID.Int32, user.ID.Int32)
	}
}

func TestGetRestaurants(t *testing.T) {
	restaurants, err := storage.NewPostgresDatabase().GetRestaurants()

	if err != nil {
		t.Errorf("error while getting restaurants : %q\n", err)
	}

	if len(restaurants) == 0 {
		t.Errorf("restaurants shouldn't be empty\n")
	}

	t.Log(restaurants)
}

func TestGetFood(t *testing.T) {
	foods, err := storage.NewPostgresDatabase().GetFoods()

	if err != nil {
		t.Errorf("error while getting foods : %q\n", err)
	}

	if len(foods) == 0 {
		t.Errorf("foods shouldn't be empty\n")
	}

	t.Log(foods)
}

func TestGetRestaurantsByFood(t *testing.T) {
	foods, err := storage.NewPostgresDatabase().GetFoods()

	if err != nil {
		t.Errorf("error while getting foods : %q\n", err)
	}

	restaurants, err := storage.NewPostgresDatabase().GetRestaurantsFromFood(foods[0].Name)

	if err != nil {
		t.Errorf("error while getting restaurants : %q\n", err)
	}

	if len(restaurants) == 0 {
		t.Errorf("restaurants shouldn't be empty\n")
	}

	t.Log(restaurants)
}

func TestGetMealFromRestaurant(t *testing.T) {
	restaurants, err := storage.NewPostgresDatabase().GetRestaurants()

	if err != nil {
		t.Errorf("error while getting restaurants : %q\n", err)
	}

	var meals []types.Meal
	only_empty := true
	for _, restaurant := range restaurants {
		meals, err = storage.NewPostgresDatabase().GetMealsFromRestaurant(restaurant.ID)
		t.Log(meals)
		if len(meals) != 0 {
			only_empty = false
			break
		}
	}

	if err != nil {
		t.Errorf("error while getting meals : %q\n", err)
	}

	if only_empty {
		t.Errorf("meals shouldn't be empty\n")
	}

	t.Log(meals)
}
