package janus

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/gocarina/gocsv"
	"github.com/jmoiron/sqlx"
	"opencoredata.org/ocdServices/connectors"
)

type AgeModelxx struct { // make an []struct?
	Leg                   int64          `json:"Leg" dc_description:"Number identifying the cruise."`
	Site                  int64          `json:"Site" dc_description:"Number identifying the site from which the core was retrieved. A site is the position of a beacon around which holes are drilled."`
	Hole                  string         `json:"Hole" dc_description:"Letter identifying the hole at a site from which a core was retrieved or data was collected."`
	Age_model_type        string         `json:"Age_model_type" dc_description:"The type of the age model - where or when it was constructed (shipboard, initial reports, post-moratorium)."`
	Depth_mbsf            string         `json:"Depth_mbsf" dc_unit:"mbsf" dc_unit_descript:"meters below sea floor (mbsf)" dc_description:"Depth in meters below sea floor (mbsf)." xs_type:"float"`
	Age_ma                string         `json:"Age_ma" dc_unit:"ma" dc_unit_descript:"millions of years (ma)" dc_description:"Age in millions of years (ma)." xs_type:"float"`
	Control_point_comment sql.NullString `json:"Control_point_comment" dc_description:"Remark on a control point for an age model."`
	Age_model_comment     sql.NullString `json:"Age_model_comment" dc_description:"Remark on an age model."`
}

// GetAgeModles function for evaluate janus calls
func GetAgeModles(request *restful.Request, response *restful.Response) {
	conn, err := connectors.GetJanusConX() // get the Oracle connection
	if err != nil {
		log.Print(err)
	}
	defer conn.Close()

	sqlstring := `SELECT * FROM ocd_age_model WHERE leg = 138 AND site = 844 AND hole = 'B'`

	data := &[]AgeModelxx{}
	GetDatav3(conn, sqlstring, &data)
	log.Print(data)

	csvContent, err := gocsv.MarshalString(data) // Get all clients as CSV string
	if err != nil {
		log.Print(err)
	}

	response.Header().Set("Content-Type", "text/csv") // setting the content type header to text/csv
	response.Header().Set("Content-Disposition", "attachment;filename=ocdDataDownload.csv")
	response.Write([]byte(fmt.Sprintf("%v", csvContent)))
}

// GetDatav2 does the search and modifies the results struct pointer
func GetDatav3(db *sqlx.DB, sqlstring string, results **[]AgeModelxx) {
	db.MapperFunc(strings.ToUpper)
	err := db.Select(*results, sqlstring)
	if err != nil {
		log.Printf(`Error v2 with: %s`, err)
	}
}
