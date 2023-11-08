package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

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

func (db *PostgresDatabase) UpdateMail(id int, mail string) error {
	client, err := db.Connect()
	if err != nil {
		return fmt.Errorf("err when updating mail : %v", err)
	}
	defer client.Close()

	query := `UPDATE users SET mail=$1 WHERE iduser=$2`

	_, err = client.Exec(query, mail, id)

	return err
}

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

func (db *PostgresDatabase) RestaurantHasBeenLiked(id_restaurant int, id_user int) (bool, error) {
	fmt.Println("id_user : ", id_user)
	client, err := db.Connect()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	query := `SELECT COUNT(*) FROM favoriterestaurant WHERE idrestaurant=$1 and iduser=$2`

	var exists int

	err = client.QueryRow(query, id_restaurant, id_user).Scan(&exists)
	fmt.Println(exists)

	return exists > 0, err
}

func (db *PostgresDatabase) GetUserByAuthCustomName(custom_name string) (types.User, error) {
	client, err := db.Connect()
	if err != nil {
		return types.User{}, err
	}
	defer client.Close()

	var user types.User

	query := `SELECT u.iduser, u.mail, u.password, u.name, u.idschool, u.nonce, u.name_modified, u.token, u.ical, u.salt FROM users u JOIN federal_credentials fc ON fc.user_id=u.iduser WHERE fc.custom_name=$1;`

	err = client.QueryRow(query, custom_name).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.IDSchool, &user.Nonce, &user.Name_modified, &user.Refresh_token, &user.Ical, &user.Salt)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.User{}, ErrUserNotFound
		} else {
			panic(err)
		}
	}
	return user, nil
}

func (db *PostgresDatabase) UpdateUser(id_user int, name string, ical string, school_id int) (*types.User, error) {
	query_update_user := `UPDATE users SET name=$1, idschool=$2, ical=$3 WHERE iduser=$4 RETURNING iduser, mail, password, name, idschool, ical, nonce, name_modified, token, salt`

	client, err := db.Connect()
	if err != nil {
		return nil, errors.New("error while connection to database")
	}
	defer client.Close()

	user, err := types.ScanUser(client.QueryRow(query_update_user, name, school_id, ical, id_user))

	return &user, err
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

type UserCreationFunc func(*UserCreationOpts) error

type UserCreationOpts struct {
	password        *string
	email           string
	activation_code string
}

func defaultUserOpts() *UserCreationOpts {
	return &UserCreationOpts{
		password:        nil,
		email:           "",
		activation_code: util.GenActivationCode(),
	}
}

func withPassword(password *string) UserCreationFunc {
	return func(uco *UserCreationOpts) error {
		if len(*password) < 6 {
			return ErrShortPassword
		}
		uco.password = password
		return nil
	}
}

func withEmail(email string) UserCreationFunc {
	return func(uco *UserCreationOpts) error {
		if match, _ := regexp.MatchString("^.*@etu\\.umontpellier\\.fr$", email); !match {
			return ErrWrongEmailFormat
		}
		uco.email = email
		return nil
	}
}

func withGmail(email string) UserCreationFunc {
	return func(uco *UserCreationOpts) error {
		uco.email = email
		return nil
	}
}

func (db *PostgresDatabase) CreateUser(opts ...UserCreationFunc) (types.User, error) {
	user_opts := defaultUserOpts()
	for _, fn := range opts {
		err := fn(user_opts)
		if err != nil {
			return types.User{}, err
		}
	}
	return db.createUser(*user_opts)
}

func (db *PostgresDatabase) CreateLocalUser(mail string, password string) (types.User, error) {
	user, err := db.CreateUser(withPassword(&password), withEmail(mail))

	client, err := db.Connect()
	if err != nil {
		return types.User{}, err
	}
	query := `INSERT INTO federal_credentials(user_id, provider, created_at, custom_name) VALUES($1, $2, $3, $4);`

	_, err = client.Exec(query, user.ID, "local", time.Now(), mail)
	defer client.Close()
	err = util.SendConfirmationMail(user.Email.String, user.Nonce.String)
	return user, err
}

func (db *PostgresDatabase) createUser(opts UserCreationOpts) (types.User, error) {
	client, err := db.Connect()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer client.Close()

	var user types.User
	if opts.password != nil { //creating local user because password is defined
		hashed_password, salt := util.HashAndSalted(*opts.password)
		query := `INSERT INTO users(mail, password, nonce, salt) VALUES ($1, $2, $3, $4) RETURNING iduser, mail, password, name, idschool, nonce, name_modified, token, ical, salt;`

		user, err = types.ScanUser(client.QueryRow(query, opts.email, hashed_password, opts.activation_code, salt))
		if err != nil {
			return types.User{}, err
		}

	} else {
		query := `INSERT INTO users(mail, nonce) VALUES ($1, $2) RETURNING iduser, mail, password, name, idschool, nonce, name_modified, token, ical, salt;`

		user, err = types.ScanUser(client.QueryRow(query, opts.email, opts.activation_code))
		if err != nil {
			return types.User{}, err
		}
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

func (db *PostgresDatabase) GetRestaurant(id_restaurant int) (types.Restaurant, error) {
	query := `SELECT idrestaurant, name, url, gpscoord FROM restaurant WHERE idrestaurant=$1`

	client, err := db.Connect()
	if err != nil {
		log.Fatalf("caught database err when opening : %q\n", err)
		return types.Restaurant{}, err
	}
	defer client.Close()

	row := client.QueryRow(query, id_restaurant)

	return types.ScanRestaurant(row)
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

func (db *PostgresDatabase) GetFavoriteRestaurants(id_user int) ([]types.Restaurant, error) {
	query := `SELECT r.idrestaurant, r.name, r.url, r.gpscoord FROM restaurant r JOIN favoriterestaurant fr ON fr.idrestaurant = r.idrestaurant WHERE fr.iduser=$1`

	client, err := db.Connect()
	if err != nil {
		log.Fatalf("caught database err when opening : %q\n", err)
		return []types.Restaurant{}, err
	}
	defer client.Close()
	return pg_util.QueryRestaurants(client, query, id_user)
}

func (db *PostgresDatabase) GetCalendarOfUser(id_user int) (string, error) {
	query := `SELECT ical FROM users WHERE iduser=$1`
	client, err := db.Connect()
	if err != nil {
		log.Fatalf("caught database err when opening : %v\n", err)
		return "", err
	}
	defer client.Close()
	var ical string
	row := client.QueryRow(query, id_user)
	err = row.Scan(&ical)
	if err != nil {
		log.Fatalf("caught error while scanning ical : %v\n", err)
	}
	return ical, err
}

func (db *PostgresDatabase) AddRestaurantAsFavorite(id_user int, id_restaurant int) error {
	query := `INSERT INTO favoriterestaurant(idrestaurant, iduser) VALUES($1, $2)`
	client, err := db.Connect()
	if err != nil {
		log.Fatalf("caught database err when opening : %v\n", err)
		return err
	}
	defer client.Close()
	_, err = client.Exec(query, id_restaurant, id_user)
	return err
}
func (db *PostgresDatabase) DeleteRestaurantFromFavorite(id_user int, id_restaurant int) error {
	query := `DELETE FROM favoriterestaurant WHERE idrestaurant=$1 AND iduser=$2`
	client, err := db.Connect()
	if err != nil {
		log.Fatalf("caught database err when opening : %v\n", err)
		return err
	}
	defer client.Close()
	_, err = client.Exec(query, id_restaurant, id_user)
	return err
}

func (db *PostgresDatabase) CreateGoogleUser(email string) (types.User, error) {
	user, err := db.CreateUser(withGmail(email))
	if err != nil {
		return types.User{}, err
	}
	client, err := db.Connect()

	query := `INSERT INTO federal_credentials(user_id, provider, created_at, custom_name) VALUES($1, $2, $3, $4);`
	if err != nil {
		log.Fatalf("caught database err when opening : %v\n", err)
		return types.User{}, err
	}
	defer client.Close()

	_, err = client.Exec(query, user.ID, "google", time.Now(), email)
	return user, err
}

func (db *PostgresDatabase) ConfirmMail(id_user int, nonce string) error {
	user, err := db.GetUserByID(id_user)
	if err != nil {
		return err
	}

	if nonce != user.Nonce.String {
		return errors.New("err : Bad nonce")
	}

	query := `UPDATE users SET nonce=NULL WHERE iduser=$1`
	client, err := db.Connect()
	if err != nil {
		return err
	}

	defer client.Close()

	_, err = client.Exec(query, id_user)
	return err
}

var ErrWrongEmailFormat = errors.New("email must finished by @etu.umontpellier.fr")
var ErrShortPassword = errors.New("password must be 6 characters long")
