package storage

import (
	"database/sql"
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/Courtcircuits/HackTheCrous.api/storage/pg_util"
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
	connStr := "user=" + db.user + " password=" + db.password + " host=" + db.host + " port=" + db.port + " dbname=" + db.database
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

func (db *PostgresDatabase) GetRestaurants() ([]types.Restaurant, error) {
	query := `SELECT idrestaurant, name, url, gpscoord FROM restaurant`

	client, err := db.Connect()
	if err != nil {
		log.Fatalf("caught database err when opening : %q\n", err)
		return []types.Restaurant{}, err
	}
	defer client.Close()

	return pg_util.QueryRestaurants(client, query)
}

func (db *PostgresDatabase) GetRestaurantByUrl(url string) (types.Restaurant, error) {
	query := `SELECT idrestaurant, name, url, gpscoord FROM restaurant WHERE url=$1`

	client, err := db.Connect()
	if err != nil {
		log.Fatalf("caught database err when opening : %q\n", err)
		return types.Restaurant{}, err
	}
	defer client.Close()

	row := client.QueryRow(query, url)

	return types.ScanRestaurant(row)
}

func (db *PostgresDatabase) SearchRestaurant(query string) ([]types.Restaurant, error) {
	sql_query := `SELECT idrestaurant, name, url, gpscoord FROM restaurant
WHERE idrestaurant IN (SELECT r.idrestaurant FROM restaurant r JOIN suggestions_restaurant sr ON sr.idrestaurant=r.idrestaurant WHERE UPPER(sr.keyword) LIKE $1)`

	query = "%" + strings.ToUpper(query) + "%"
	log.Println(query)

	client, err := db.Connect()
	if err != nil {
		log.Fatalf("caught database err when opening : %q\n", err)
		return []types.Restaurant{}, err
	}
	defer client.Close()

	return pg_util.QueryRestaurants(client, sql_query, query)
}

func (db *PostgresDatabase) SearchRestaurantByName(name string) ([]types.Restaurant, error) {
	query := `SELECT idrestaurant, name, url, gpscoord FROM restaurant WHERE UPPER(name) LIKE $1`
	name = "%" + strings.ToUpper(name) + "%"

	client, err := db.Connect()
	if err != nil {
		log.Fatalf("caught database err when opening : %q\n", err)
		return []types.Restaurant{}, err
	}
	defer client.Close()

	return pg_util.QueryRestaurants(client, query, name)
}

func (db *PostgresDatabase) SearchSchoolByName(name string) ([]types.School, error) {
	query := "SELECT " + types.SCHOOL_SQL_ATTR + " FROM school WHERE UPPER(name) LIKE $1"

	name = "%" + strings.ToUpper(name) + "%"

	client, err := db.Connect()
	if err != nil {
		log.Fatalf("caught database err when opening : %q\n", err)
		return []types.School{}, err
	}
	defer client.Close()

	return pg_util.QuerySchools(client, query, name)
}

func (db *PostgresDatabase) GetSchoolOfUser(id_user int) (types.School, error) {
	query := `SELECT s.idschool, s.name, s.coords FROM school s JOIN users u ON u.idschool = s.idschool WHERE u.iduser=$1`
	client, err := db.Connect()
	if err != nil {
		log.Fatalf("caught database err when opening : %q\n", err)
		return types.School{}, err
	}
	defer client.Close()

	return types.ScanSchool(client.QueryRow(query, id_user))
}

func (db *PostgresDatabase) GetMealsFromRestaurant(id_restaurant int) ([]types.Meal, error) {
	query := `SELECT idmeal, typemeal, foodies, day FROM meal WHERE idrestaurant=$1`
	client, err := db.Connect()
	if err != nil {
		log.Fatalf("caught database err when opening : %q\n", err)
		return []types.Meal{}, err
	}
	defer client.Close()
	rows, err := client.Query(query, id_restaurant)
	if err != nil {
		log.Fatalf("caught database er when querying : %q\n", err)
	}
	defer rows.Close()
	var meals []types.Meal

	for rows.Next() {
		meal, err := types.ScanMeals(rows)
		if err != nil {
			log.Fatalf("caught database err when iterating through meals : %q\n", err)
			return []types.Meal{}, err
		}
		meals = append(meals, meal)
	}
	return meals, nil
}

var ErrWrongEmailFormat = errors.New("email must finished by @etu.umontpellier.fr")
var ErrShortPassword = errors.New("password must be 6 characters long")
