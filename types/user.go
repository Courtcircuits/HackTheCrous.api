package types

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/Courtcircuits/HackTheCrous.api/graph/model"
	"github.com/Courtcircuits/HackTheCrous.api/util"
)

type User struct {
	ID            sql.NullInt32  `json:"id"`
	Email         sql.NullString `json:"email"`
	Password      sql.NullString `json:"password"`
	Name          sql.NullString `json:"name"`
	IDSchool      sql.NullInt32  `json:"idschool"`
	Ical          sql.NullString `json:"ical"`
	Nonce         sql.NullString `json:"nonce"`
	Name_modified sql.NullString `json:"name_modified"`
	Refresh_token sql.NullString `json:"refresh_token"`
	Auth_token    sql.NullString `json:"auth_token"`
	Salt          sql.NullString `json:"salt"`
}

type Tokens struct {
	Refresh_token string `json:"refresh_token"`
	Auth_token    string `json:"token"`
}

func (u *User) GenAuthToken() string {
	log.Printf("id  :%d, mail: %q\n", u.ID.Int32, u.Email.String)
	expiration := time.Now().Add(time.Hour * 24) // tomorrow
	auth_payload := map[string]any{
		"iduser": int(u.ID.Int32),
		"mail":   u.Email.String,
	}
	return util.GenJWT(expiration, auth_payload)
}

var ErrRefreshTokenNeedUpdate error = errors.New("need update")

func (u *User) GetTokens(mustBeRemembered bool) (Tokens, error) {
	auth_token := u.GenAuthToken()

	log.Printf("auth token generated : %q \n", auth_token)

	if !mustBeRemembered { //refresh token is null
		return Tokens{
			Refresh_token: "",
			Auth_token:    auth_token,
		}, nil
	}

	if !u.Refresh_token.Valid {
		return Tokens{
			Refresh_token: "",
			Auth_token:    auth_token,
		}, ErrRefreshTokenNeedUpdate
	}

	return Tokens{
		Refresh_token: u.Refresh_token.String,
		Auth_token:    auth_token,
	}, nil
}

func ScanUser(row *sql.Row) (User, error) {
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.IDSchool, &user.Nonce, &user.Name_modified, &user.Refresh_token, &user.Ical, &user.Salt)

	return user, err
}

func (u *User) CheckPassword(given_password string) bool {
	return util.CompareHash(u.Password.String, given_password, u.Salt.String)
}

func (u User) ToGraphQL() *model.User {
	ID := int(u.ID.Int32)
	return &model.User{
		Iduser: &ID,
		Name:   &u.Name.String,
		Mail:   &u.Email.String,
		Ical:   &u.Ical.String,
		Nonce:  &u.Nonce.Valid,
	}
}
