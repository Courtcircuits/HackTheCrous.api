package pg_util

import (
	"database/sql"
	"log"

	"github.com/Courtcircuits/HackTheCrous.api/types"
)

func QuerySchools(client *sql.DB, query string, args ...any) ([]types.School, error) {
	var rows *sql.Rows
	var err error
	if len(args) == 0 {
		rows, err = client.Query(query)
	} else {
		rows, err = client.Query(query, args...)
	}

	if err != nil {
		log.Fatalf("caught database err when querying : %q\n", err)
		return []types.School{}, err
	}
	defer rows.Close()
	var schools []types.School

	for rows.Next() {
		school, err := types.ScanSchools(rows)
		if err != nil {
			log.Fatalf("caught database err when iterating through schools : %q\n", err)
			return []types.School{}, err
		}
		schools = append(schools, school)
		log.Println(school)
	}

	return schools, nil
}
