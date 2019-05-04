package expeditions

import (
	"bytes"
	"encoding/json"
	"log"
	"text/template"

	"github.com/emicklei/go-restful"
	utilities "opencoredata.org/ocdServices/internal/utilities"
)

type TaxaReply struct {
	Head    *TaxaHead
	Results *TaxaResults
}

type TaxaHead struct {
	Link []string
	Vars []string
}

type TaxaResults struct {
	Distinct bool
	Ordered  bool
	Bindings []TaxaBindings
}

type TaxaBindings struct {
	Pro TaxaEntryType
	Vol TaxaEntryType
}

type TaxaEntryType struct {
	Type     string
	Datatype string
	Value    string
}

// TODO...
// need CSDCO project  /csdco/{project}
// return data in schema.org or other metadata format in JSON

func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/expeditions").
		Doc("Return Metadata on an expedition or project").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.GET("/{leg}").To(ExpCall).
		Doc("Get expedition or project information").
		Param(service.PathParameter("leg", "Leg in format like 123 or 312U").DataType("string")).
		Operation("ExpCall"))

	return service
}

func ExpCall(request *restful.Request, response *restful.Response) {

	const SPARQL = `
PREFIX iodp: <http://data.oceandrilling.org/core/1/>
SELECT DISTINCT  ?pro ?vol
FROM <http://data.oceandrilling.org/codices#>
WHERE {
   ?uri iodp:expedition "113" .
   ?uri iodp:scientificprospectus ?pro .
   ?uri iodp:scientificreportvolume ?vol
}
`

	// create the SPARQL call from a template
	var buff = bytes.NewBufferString("")
	t, err := template.New("lith template").Parse(SPARQL)
	if err != nil {
		log.Printf("lith template creation failed: %s", err)
	}
	//err = t.Execute(buff, id) //  instead of os.Stdout create a function to call.
	err = t.Execute(buff, "id")
	if err != nil {
		log.Printf("lith template execution failed: %s", err)
	}

	log.Printf("lith sparql call here:\n %s", string(buff.Bytes()))

	reply := []byte(utilities.SparqlCallFunc(string(buff.Bytes())))
	var m TaxaReply
	data := json.Unmarshal(reply, &m)
	if err != nil {
		log.Printf("lith json unmarshall execution failed: %s", data)
	}
	log.Printf("lith JSON here:\n %s", reply)

	//return m
	//  Must it be a struct?
	log.Println("Going to try and send expedition info now")
	response.WriteEntity(m)

}
