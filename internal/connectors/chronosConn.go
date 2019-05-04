package connectors

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

// var (   //  need flags for the items in the open string
// 	fUsername    = flag.String("db.username", "username", "username to connect as (if you don't provide the dsn")
// 	fPassword    = flag.String("db.password", "password", "password to connect with (if you don't provide the dsn")
// )

// GetConnection returns a connection - using GetDSN if dsn is empty
func ChronosConn(username string, password string) (*sql.DB, error) {

	paramString := fmt.Sprintf("user=%s password=%s host=localhost port=9000 dbname=neptune sslmode=disable", username, password)
	db, err := sql.Open("postgres", paramString)
	if err != nil {
		log.Fatal(err)
	}

	return db, err
}
