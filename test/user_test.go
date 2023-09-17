package test

import (
	"database/sql"
	"testing"

	"github.com/Courtcircuits/HackTheCrous.api/types"
)

func TestGenAuthToken(t *testing.T) {
	u := types.User{
		ID: sql.NullInt32{
			Int32: 3,
			Valid: true,
		},
		Email: sql.NullString{
			String: "testest@etu.umontpellier.fr",
			Valid:  true,
		},
	}

	auth_token := u.GenAuthToken()

	t.Log(auth_token)

	if auth_token == "" {
		t.Errorf("empty token \n")
	}
}
