package janus

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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
		log.Fatal(err)
	}
	defer conn.Close()

	sqlstring := `SELECT * FROM ocd_age_model WHERE leg = 138 AND site = 844 AND hole = 'B'`

	// Test 1  rows to CSV
	output2, _ := TestFuncx(conn, sqlstring)

	// TEST 2  rows to struct to X
	data := []AgeModelx{}
	GetData(conn, sqlstring, &data)
	log.Print("testng data value")
	log.Println(data)

	response.Header().Set("Content-Type", "text/csv") // setting the content type header to text/csv
	response.Header().Set("Content-Disposition", "attachment;filename=ocdDataDownload.csv")
	response.Write([]byte(fmt.Sprintf("%v", output2)))
}

// TestFuncx is a test function for better formatting style....
// Keep as a called function to allow use in other code (like ocdBulk) and to future use like gRPC
// func TestFuncx(db *sqlx.DB) (*[]AgeModelx, error) {
func TestFuncx(db *sqlx.DB, sqlstring string) (string, error) {
	results, err := db.Queryx(sqlstring)
	if err != nil {
		log.Printf(`Error with: %s`, err)
	}
	defer results.Close()

	csvdata, _ := ResultsToCSV(results)
	return csvdata, nil
}

func GetData(db *sqlx.DB, sqlstring string, destold interface{}) {
	// arr := reflect.ValueOf(dest).Elem()
	// v := reflect.New(reflect.TypeOf(dest).Elem().Elem())

	db.MapperFunc(strings.ToUpper)

	// rows, err := db.Queryx(sqlstring)
	// if err != nil {
	// 	log.Printf(`Error 1 with: %s`, err)
	// }

	places := []AgeModelx{}
	err := db.Select(&places, sqlstring)
	if err != nil {
		log.Printf(`Error 2 with: %s`, err)
		return
	}
	log.Print("getdata places value")
	log.Print(places)

	// var dest []interface{}
	// for rows.Next() {
	// 	// cols is an []interface{} of all of the column results
	// 	dest, err = rows.SliceScan()
	// 	if err != nil {
	// 		log.Printf(`Error in row Next: %s`, err)
	// 	}
	// }

	// log.Print("getdata dest value")
	// log.Print(dest)

	// if err == nil {
	// 	if err = rows.StructScan(v.Interface()); err == nil {
	// 		arr.Set(reflect.Append(arr, v.Elem()))
	// 	}
	// }
}
