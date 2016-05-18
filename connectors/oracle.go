package connectors

import (
	"database/sql"
	"fmt"
	// _ "gopkg.in/goracle.v1"
   	_ "gopkg.in/rana/ora.v3"  
	"os"

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
