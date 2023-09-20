package test

import (
	"testing"

	"github.com/Courtcircuits/HackTheCrous.api/storage"
)

func TestGetUserById(t *testing.T) {
	id := 50
	user, err := storage.NewPostgresDatabase().GetUserByID(id)

	if err != nil {
		t.Errorf("shouldn't be throwing error : %q", err)
	}

	if user.ID.Int32 != 50 {
		t.Errorf("ID is %q but must be 50", user.ID.Int32)
	}
	if user.Email.String != "test@test.com" {
		t.Errorf("Mail is %q but must be test@test.com", user.Email.String)
	}
}

func TestCreateUser(t *testing.T) {
	email := "testtest@etu.umontpellier.fr"
	password := "12341234"
	pg_storage := storage.NewPostgresDatabase()

	err_delete := pg_storage.DeleteUserByMail(email)

	if err_delete != nil {
		t.Errorf("error when delete : %q", err_delete)
	}

	user, err := pg_storage.CreateUser(email, password)

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
