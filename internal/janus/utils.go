package janus

import (
	"bytes"
	"encoding/csv"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// ResultsToCSVViaStruct takes a sqlx.result and struct and
// returns CSV (or also JSON?)
// question:  do I need an ageMOdel stub and then a handler for struct scanning?
//http://stackoverflow.com/questions/30855499/using-sqlx-rows-structscan-for-interface-args
func ResultsToCSVViaStruct(results *sqlx.Rows, currentStruct AgeModel) (string, error) {

	return "test", nil
}

// ResultsToCSV takes a sqlx.result without a struct
// and returns CSV string
func ResultsToCSV(results *sqlx.Rows) (string, error) {

	b := &bytes.Buffer{} // creates IO Writer
	csvWriter := csv.NewWriter(b)
	//csvWriter.Comma = 0x0009  // tab sep

	firstLine := true
	for results.Next() {
		row, err := results.SliceScan() // or do a results.StructScan(struct)
		if err != nil {
			log.Fatal(err)
		}
		if firstLine {
			firstLine = false
			cols, err := results.Columns()
			if err != nil {
				log.Fatal(err)
			}
			csvWriter.Write(cols)
		}

		rowStrings := make([]string, len(row))
		for i, col := range row {
			//log.Print(reflect.TypeOf(col))
			switch col.(type) {
			case float64:
				rowStrings[i] = strconv.FormatFloat(col.(float64), 'f', 6, 64)
			case int64:
				rowStrings[i] = strconv.FormatInt(col.(int64), 10)
			case bool:
				rowStrings[i] = strconv.FormatBool(col.(bool))
			case []byte:
				rowStrings[i] = string(col.([]byte))
			case string:
				rowStrings[i] = strings.Trim(col.(string), " ")
			// case sql.NullString:  // cannot use col.(sql.NullString) (type sql.NullString) as type string in argument to strings.Trim
			// 	rowStrings[i] = strings.Trim(col.(sql.NullString), " ")
			case time.Time:
				rowStrings[i] = col.(time.Time).String()
			case nil:
				rowStrings[i] = "NULL"
			default:
				log.Print(col)
			}
		}
		csvWriter.Write(rowStrings)
	}

	csvWriter.Flush()
	return b.String(), nil

}
