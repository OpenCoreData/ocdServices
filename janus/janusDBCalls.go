package janus

import (
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

// GetAgeModelData does the search and modifies the results struct pointer
// Not placed in the above function to allow this to be called by stand along apps like ocdBulk
func GetAgeModelData(db *sqlx.DB, sqlstring string, results **[]AgeModel) {
	db.MapperFunc(strings.ToUpper)
	err := db.Select(*results, sqlstring)
	if err != nil {
		log.Printf(`Error v2 with: %s`, err)
	}
}
