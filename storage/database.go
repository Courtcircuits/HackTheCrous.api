package storage

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/Courtcircuits/HackTheCrous.api/util"
	"github.com/gofiber/fiber/v2/log"
)

type DatabaseV1 struct {
	user     string
	password string
	host     string
	port     string
	database string
}

type IDatabase interface {
	Connect() (*sql.DB, error)
	Select(query string, args ...any) (*sql.Rows, error)
}

func newDatabase() *DatabaseV1 {
	return &DatabaseV1{
		util.Get("PG_USER"),
		util.Get("PG_PASSWORD"),
		util.Get("PG_HOST"),
		util.Get("PG_PORT"),
		util.Get("PG_DATABASE"),
	}
}

func (db *DatabaseV1) Connect() (*sql.DB, error) {
	connStr := "user=" + db.user + " password=" + db.password + " host=" + db.host + " port=" + db.port + " dbname=" + db.database + " sslmode=disable"
	return sql.Open("postgres", connStr)
}

func (db *DatabaseV1) Select(query string, args ...any) (*sql.Rows, error) {
	conn, err := db.Connect()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer conn.Close()

	if len(args) == 0 {
		return conn.Query(query)
	}
	return conn.Query(query, args...)
}

var Database IDatabase = newDatabase()

// func (db *Database) Insert(query string, args ...string) (sql.Result, error) {
// 	conn, err := db.Connect()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer conn.Close()

// 	return conn.Exec(query, args)
// }
