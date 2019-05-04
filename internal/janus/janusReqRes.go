package janus

import (
	// "database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/gocarina/gocsv"
	"opencoredata.org/ocdServices/internal/connectors"
)

// GetAgeModle function for evaluate janus calls
func GetAgeModel(request *restful.Request, response *restful.Response) {
	conn, err := connectors.GetJanusConX() // get the Oracle connection
	if err != nil {
		log.Print(err)
		http.Error(response, err.Error(), 500)
		return
	}
	defer conn.Close()

	sqlstring := `SELECT * FROM ocd_age_model WHERE leg = 138 AND site = 844 AND hole = 'B'`

	// Call to Struct based version
	data := &[]AgeModel{}
	GetAgeModelData(conn, sqlstring, &data)
	csvContent, err := gocsv.MarshalString(data) // Get all clients as CSV string
	if err != nil {
		log.Print(err)
	}

	response.Header().Set("Content-Type", "text/csv") // setting the content type header to text/csv
	response.Header().Set("Content-Disposition", "attachment;filename=ocdDataDownload.csv")
	response.Write([]byte(fmt.Sprintf("%v", csvContent)))
}
