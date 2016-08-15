package janus

import (
	// "database/sql"

	"database/sql"
	"fmt"
	"log"

	"github.com/emicklei/go-restful"
	"github.com/jmoiron/sqlx"
	"opencoredata.org/ocdServices/connectors"
)

type AgeModelx struct { // make an []struct?
	Leg                   int64          `json:"Leg" dc_description:"Number identifying the cruise."`
	Site                  int64          `json:"Site" dc_description:"Number identifying the site from which the core was retrieved. A site is the position of a beacon around which holes are drilled."`
	Hole                  string         `json:"Hole" dc_description:"Letter identifying the hole at a site from which a core was retrieved or data was collected."`
	Age_model_type        string         `json:"Age_model_type" dc_description:"The type of the age model - where or when it was constructed (shipboard, initial reports, post-moratorium)."`
	Depth_mbsf            string         `json:"Depth_mbsf" dc_unit:"mbsf" dc_unit_descript:"meters below sea floor (mbsf)" dc_description:"Depth in meters below sea floor (mbsf)." xs_type:"float"`
	Age_ma                string         `json:"Age_ma" dc_unit:"ma" dc_unit_descript:"millions of years (ma)" dc_description:"Age in millions of years (ma)." xs_type:"float"`
	Control_point_comment sql.NullString `json:"Control_point_comment" dc_description:"Remark on a control point for an age model."`
	Age_model_comment     sql.NullString `json:"Age_model_comment" dc_description:"Remark on an age model."`
}

// TestNG function for evaluate janus calls
func TestNGx(request *restful.Request, response *restful.Response) {

	// get the Oracle connection
	conn, err := connectors.GetJanusConX()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//call func test
	output2, _ := TestFuncx(conn)
	response.Write([]byte(fmt.Sprintf("%v", output2)))

	// ctx.ResponseWriter.Header().Set("Content-Type", "text/csv") // setting the content type header to text/csv
	// ctx.ResponseWriter.Header().Set("Content-Type", "text/csv")
	// ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment;filename=TheCSVFileName.csv")
	// ctx.ResponseWriter.Write(b.Bytes())

}

// TestFunc is a test function for better formatting style....
// Keep as a called function to allow use in other code (like ocdBulk) and to future use like gRPC
// func TestFuncx(db *sqlx.DB) (*[]AgeModelx, error) {
func TestFuncx(db *sqlx.DB) (string, error) {

	results, err := db.Queryx(`SELECT * FROM ocd_age_model WHERE leg = 138 AND site = 844 AND hole = 'B'`)
	if err != nil {
		log.Printf(`Error with: %s`, err)
	}
	defer results.Close()

	csvdata, _ := ResultsToCSV(results)

	// var teststruct AgeModelx  // does this need to be an array?
	// test, _ := ResultsToCSVViaStruct(results, teststruct)

	return csvdata, nil

}
