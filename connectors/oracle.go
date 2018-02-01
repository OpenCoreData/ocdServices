package connectors

import (
	"database/sql"
	"fmt"
	// _ "gopkg.in/goracle.v1"
	"os"

	"github.com/jmoiron/sqlx"
	_ "gopkg.in/rana/ora.v4"
)

// GetJanusCon returns a Oracle DB connection
func GetJanusCon() (*sql.DB, error) {

	username := os.Getenv("JANUS_USERNAME")
	password := os.Getenv("JANUS_PASSWORD")
	host := os.Getenv("JANUS_HOST")
	servicename := os.Getenv("JANUS_SERVICENAME")

	connectionString := fmt.Sprintf("%s/%s@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=TCP)(HOST=%s)(PORT=1521)))(CONNECT_DATA=(SERVICE_NAME=%s)))", username, password, host, servicename)
	// return sql.Open("goracle", connectionString)
	return sql.Open("ora", connectionString)
}

// GetJanusConX opens with the SQLX set for testing....
func GetJanusConX() (*sqlx.DB, error) {

	username := os.Getenv("JANUS_USERNAME")
	password := os.Getenv("JANUS_PASSWORD")
	host := os.Getenv("JANUS_HOST")
	servicename := os.Getenv("JANUS_SERVICENAME")

	connectionString := fmt.Sprintf("%s/%s@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=TCP)(HOST=%s)(PORT=1521)))(CONNECT_DATA=(SERVICE_NAME=%s)))", username, password, host, servicename)
	// return sql.Open("goracle", connectionString)
	return sqlx.Open("ora", connectionString)
}
