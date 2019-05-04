package connectors

import (
	"fmt"
	// _ "gopkg.in/goracle.v1"
	"os"

	"github.com/jmoiron/sqlx"
	// _ "gopkg.in/rana/ora.v4"
)

// Note to future me..   GetJanusCon was used in one location
// but if needs the Oracle client installed to work..   need to break this dep.

// GetJanusCon returns a Oracle DB connection
// func GetJanusCon() (*sql.DB, error) {

// 	username := os.Getenv("JANUS_USERNAME")
// 	password := os.Getenv("JANUS_PASSWORD")
// 	host := os.Getenv("JANUS_HOST")
// 	servicename := os.Getenv("JANUS_SERVICENAME")

// 	connectionString := fmt.Sprintf("%s/%s@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=TCP)(HOST=%s)(PORT=1521)))(CONNECT_DATA=(SERVICE_NAME=%s)))", username, password, host, servicename)
// 	// return sql.Open("goracle", connectionString)
// 	return sql.Open("ora", connectionString)
// }

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
