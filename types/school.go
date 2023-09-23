package types

import (
	"database/sql"

	"github.com/Courtcircuits/HackTheCrous.api/graph/model"
	"github.com/Courtcircuits/HackTheCrous.api/util"
)

type School struct {
	ID     int              `json:"id,omitempty"`
	Name   string           `json:"name,omitempty"`
	Coords util.Coordinates `json:"coords,omitempty"`
}

const SCHOOL_SQL_ATTR string = "idschool, name, coords"

func ScanSchools(rows *sql.Rows) (School, error) {
	type SQL_School struct {
		Idschool sql.NullInt32  `json:"idschool,omitempty"`
		Name     sql.NullString `json:"name,omitempty"`
		Coords   sql.NullString `json:"coords,omitempty"`
	}

	var sql_school SQL_School
	err := rows.Scan(&sql_school.Idschool, &sql_school.Name, &sql_school.Coords)

	if err != nil {
		return School{}, err
	}

	coords, err := util.Parse_coordinates(sql_school.Coords.String)

	return School{
		ID:     int(sql_school.Idschool.Int32),
		Name:   sql_school.Name.String,
		Coords: coords,
	}, err
}

func ScanSchool(row *sql.Row) (School, error) {
	type SQL_School struct {
		Idschool sql.NullInt32  `json:"idschool,omitempty"`
		Name     sql.NullString `json:"name,omitempty"`
		Coords   sql.NullString `json:"coords,omitempty"`
	}

	var sql_school SQL_School
	err := row.Scan(&sql_school.Idschool, &sql_school.Name, &sql_school.Coords)

	if err != nil {
		return School{}, err
	}

	coords, err := util.Parse_coordinates(sql_school.Coords.String)

	return School{
		ID:     int(sql_school.Idschool.Int32),
		Name:   sql_school.Name.String,
		Coords: coords,
	}, err
}

func (s School) ToGraphQL() *model.School {
	return &model.School{
		Idschool: &s.ID,
		Name:     &s.Name,
		Coords:   s.Coords.ToGraphQL(),
	}
}
