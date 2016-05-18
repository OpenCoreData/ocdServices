package janus

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
    "strings"
	"github.com/emicklei/go-restful"
	"opencoredata.org/ocdServices/connectors"
)

func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/janus").
		Doc("Access Open Core Data Documents").
		Consumes("application/x-www-form-urlencoded").
		Produces("text/plain")

	service.Route(service.GET("/rockeval").To(RockEval).
		Doc("Rock Eval").
		Param(service.QueryParameter("leg", "Leg of expedition")).
		Param(service.QueryParameter("site", "Site of expedition")).
		Param(service.QueryParameter("hole", "Hole of expedition")).
		Param(service.QueryParameter("core", "Core")).
		Param(service.QueryParameter("section", "Core section")).
		Param(service.QueryParameter("depthtop", "Depth top")).
		Param(service.QueryParameter("depthbottom", "Depth bottom")).
		Operation("Rock evaluation query"))

	return service
}

// RockEval function for evaluate janus calls
func RockEval(request *restful.Request, response *restful.Response) {

	// get the Oracle connection
	conn, err := connectors.GetJanusCon()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//get the parameters into a struct for the template
	type lshStruct struct {
		Leg         string
		Site        string
		Hole        string
		Core        string
		Section     string
		DepthTop    string
		DepthBottom string
	}

	lshParams := lshStruct{Leg: request.QueryParameter("leg"),
		Site:        request.QueryParameter("site"),
		Hole:        request.QueryParameter("hole"),
		Section:     request.QueryParameter("section"),
		DepthTop:    request.QueryParameter("depthtop"),
		DepthBottom: request.QueryParameter("depthbottom")}

	// get the template and populate
	lshqry := SQL_RockEval
	ct, err := template.New("RDF template").Parse(lshqry)
	if err != nil {
		log.Printf("Template creation failed for query: %s", err)
	}

	var buff = bytes.NewBufferString("")
	err = ct.Execute(buff, lshParams)
	if err != nil {
		log.Printf("Template execution failed: %s", err)
	}

	// do the query
	rows, err := conn.Query(string(buff.Bytes()))
	if err != nil {
		log.Printf(`Error with "%s": %s :%v`, string(buff.Bytes()), err, rows)
		// return
	}

	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("Failed to get columns", err)
		return
	}

	// Result is your slice string.
	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	var buffer bytes.Buffer
	
	log.Printf("%v", cols)

	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			fmt.Println("Failed to scan row", err)
			return
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}

		buffer.WriteString(strings.Join(result, ","))
		buffer.WriteString("\n")

	}

	response.Write([]byte(buffer.String()))

}
