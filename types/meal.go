package types

import (
	"database/sql"
	"time"
)

type Period string

const (
	DEJEUNER      Period = "Déjeuner"
	PETITDEJEUNER Period = "Petit déjeuner"
	DINER         Period = "Dîner"
)

type Meal struct {
	ID      int
	Type    Period
	Foodies any
	Day     time.Time
}

func ScanMeals(rows *sql.Rows) {
	panic("to implement")
}
