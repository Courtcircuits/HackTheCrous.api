package storage

import (
	"database/sql"
	"errors"
	"regexp"

	"github.com/Courtcircuits/HackTheCrous.api/types"
	"github.com/Courtcircuits/HackTheCrous.api/util"
	_ "github.com/lib/pq"
)

type PostgresDatabase struct {
	user     string
	password string
	host     string
	port     string
	database string
}

func NewPostgresDatabase() *PostgresDatabase {
	return &PostgresDatabase{
		util.Get("PG_USER"),
		util.Get("PG_PASSWORD"),
		util.Get("PG_HOST"),
		util.Get("PG_PORT"),
		util.Get("PG_DATABASE"),
	}
}

func (db *PostgresDatabase) Connect() (*sql.DB, error) {
	connStr := "user=" + db.user + " password=" + db.password + " host=" + db.host + " port=" + db.port + " dbname=" + db.database + " sslmode=require"
	client, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return client, nil
}

var ErrUserNotFound = errors.New("user not found")

// Search a user in the DB by ID and returns an object of type User
func (db *PostgresDatabase) GetUserByID(id int) (types.User, error) {
	client, err := db.Connect()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	query := `SELECT iduser, mail, password, name, idschool, nonce, name_modified, token, ical, salt FROM users WHERE iduser = $1;`
	user, err := types.ScanUser(client.QueryRow(query, id))

	if err != nil {
		if err == sql.ErrNoRows {
			return types.User{}, ErrUserNotFound
		} else {
			panic(err)
		}
	}
	return user, nil
}

// get a user from his mail, if not found throw an error
func (db *PostgresDatabase) GetUserByEmail(mail string) (types.User, error) {
	client, err := db.Connect()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	var user types.User

	query := `SELECT iduser, mail, password, name, idschool, nonce, name_modified, token, ical, salt FROM users WHERE mail=$1;`

	err = client.QueryRow(query, mail).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.IDSchool, &user.Nonce, &user.Name_modified, &user.Refresh_token, &user.Ical, &user.Salt)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.User{}, ErrUserNotFound
		} else {
			panic(err)
		}
	}
	return user, nil
}

func (db *PostgresDatabase) UpdateRefreshToken(id_user int) string {
	client, err := db.Connect()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	refresh_token := util.GenRefreshToken()

	query := `UPDATE users SET token=$1 WHERE iduser=$2;`

	client.QueryRow(query, refresh_token, id_user)

	return refresh_token
}

func (db *PostgresDatabase) CreateUser(email string, password string) (types.User, error) {

	if len(password) < 6 {
		return types.User{}, ErrShortPassword
	}

	if match, _ := regexp.MatchString("^.*@etu\\.umontpellier\\.fr$", email); !match {
		return types.User{}, ErrWrongEmailFormat
	}

	activation_code := util.GenActivationCode()
	hashed_password, salt := util.HashAndSalted(password)
	query := `INSERT INTO users(mail, password, nonce, salt) VALUES ($1, $2, $3, $4) RETURNING iduser, mail, password, name, idschool, nonce, name_modified, token, ical, salt;`

	client, err := db.Connect()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	user, err := types.ScanUser(client.QueryRow(query, email, hashed_password, activation_code, salt))

	if err != nil {
		return types.User{}, err
	}

	return user, nil
}

func (db *PostgresDatabase) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE iduser=$1`
	client, err := db.Connect()

	if err != nil {
		panic(err)
	}

	defer client.Close()
	err = client.QueryRow(query, id).Scan()

	return err
}

func (db *PostgresDatabase) DeleteUserByMail(email string) error {
	query := `DELETE FROM users WHERE mail=$1`
	client, err := db.Connect()

	if err != nil {
		panic(err)
	}

	defer client.Close()
	err = client.QueryRow(query, email).Scan()

	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

var ErrWrongEmailFormat = errors.New("email must finished by @etu.umontpellier.fr")
var ErrShortPassword = errors.New("password must be 6 characters long")
