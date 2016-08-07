package janus

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/emicklei/go-restful"
	"github.com/kisielk/sqlstruct"
	"opencoredata.org/ocdServices/connectors"
)

type AgeModel struct {
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
func TestNG(request *restful.Request, response *restful.Response) {

	// get the Oracle connection
	conn, err := connectors.GetJanusCon()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//call func test
	output2, _ := TestFunc(conn)
	response.Write([]byte(fmt.Sprintf("%v", output2)))

}

// TestFunc is a test function for better formatting style....
func TestFunc(db *sql.DB) (*AgeModel, error) {
	var r AgeModel

	rows, err := db.Query(`SELECT * FROM ocd_age_model WHERE leg = 138 AND site = 844 AND hole = 'B'`)
	if err != nil {
		log.Printf(`Error with: %s`, err)
	}
	defer rows.Close()

	for rows.Next() {
		err := sqlstruct.Scan(&r, rows) //.Scan(&r.Leg, &r.Site, &r.Hole, &r.Age_model_type, &r.Depth_mbsf, &r.Age_ma, &r.Control_point_comment, &r.Age_model_comment)
		if err != nil {
			log.Print(err)
		}
		log.Print(r)
	}

	// if err != nil {
	// 	log.Print(err)
	// 	return nil, nil //errors.Wrap(err, "failed to select t1")  // firgure out this error issue
	// }

	return &r, nil
}
