package janus

import (
	"database/sql"
	// "strings"
)

// To use GoCSV package I have to subclass sql.NULLString
// ref: https://github.com/gocarina/gocsv/issues/43#issuecomment-244700051

type SQLNullString struct {
	sql.NullString
}

func (s SQLNullString) MarshalCSV() (string, error) {
	return s.String, nil
}

// func (s string) MarshalCSV() (string, error) {
// 	return strings.TrimSpace(s), nil
// }

type AgeModel struct { // make an []struct?
	Leg                   int64         `json:"Leg" dc_description:"Number identifying the cruise."`
	Site                  int64         `json:"Site" dc_description:"Number identifying the site from which the core was retrieved. A site is the position of a beacon around which holes are drilled."`
	Hole                  string        `json:"Hole" dc_description:"Letter identifying the hole at a site from which a core was retrieved or data was collected."`
	Age_model_type        string        `json:"Age_model_type" dc_description:"The type of the age model - where or when it was constructed (shipboard, initial reports, post-moratorium)."`
	Depth_mbsf            string        `json:"Depth_mbsf" dc_unit:"mbsf" dc_unit_descript:"meters below sea floor (mbsf)" dc_description:"Depth in meters below sea floor (mbsf)." xs_type:"float"`
	Age_ma                string        `json:"Age_ma" dc_unit:"ma" dc_unit_descript:"millions of years (ma)" dc_description:"Age in millions of years (ma)." xs_type:"float"`
	Control_point_comment SQLNullString `json:"Control_point_comment" dc_description:"Remark on a control point for an age model."`
	Age_model_comment     SQLNullString `json:"Age_model_comment" dc_description:"Remark on an age model."`
}
