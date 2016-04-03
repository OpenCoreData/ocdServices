package utilities

import (
	"bytes"
	"fmt"
	sparql "github.com/knakk/sparql"
	"log"
	"time"
    "strconv"
)

const queries = `
# Comments are ignored, except those tagging a query.

# tag: my-query
SELECT *
WHERE {
  ?s ?p ?o
} LIMIT {{.Limit}} OFFSET {{.Offset}}

#tag: legcall
PREFIX iodp: <http://data.oceandrilling.org/core/1/>
SELECT DISTINCT  ?pro ?vol
FROM <http://data.oceandrilling.org/codices#>
WHERE {
   ?uri iodp:expedition "{{.Leg}}" .
   ?uri iodp:scientificprospectus ?pro .
   ?uri iodp:scientificreportvolume ?vol
}

#tag: legcall2
PREFIX iodp: <http://data.oceandrilling.org/core/1/>
SELECT DISTINCT  ?p ?o
FROM <http://data.oceandrilling.org/codices#>
WHERE {
   <http://data.oceandrilling.org/codices/lsh/113/689> ?p ?o .
}

#tag: agemodel
PREFIX chronos: <http://www.chronos.org/loc-schema#>
PREFIX geo: <http://www.w3.org/2003/01/geo/wgs84_pos#>
SELECT  ?age ?depth  
FROM <http://chronos.org/janusamp#>
WHERE {
  ?ob chronos:age ?age . 
  ?ob chronos:depth ?depth . 
  ?ob <http://purl.org/linked-data/cube#dataSet>  ?dataset .
  ?dataset geo:long ?long .
  ?dataset geo:lat ?lat .
  ?dataset <http://www.w3.org/2000/01/rdf-schema#seeAlso> {{.URI}} .
}
ORDER BY ?dataset ASC(?depth)


`

// TODO ..  error nice if I don't get what I expect in these params
func MaxAge(leg string, site string, hole string)  float64 {
	repo, err := sparql.NewRepo("http://data.oceandrilling.org/sparql",
		//sparql.DigestAuth("dba", "dba"),
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Print(err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)
    
    fmt.Printf("Making call for %s %s %s\n", leg, site, hole)

	q, err := bank.Prepare("agemodel",struct{ URI string }{fmt.Sprintf("<http://data.oceandrilling.org/codices/lsh/%s/%s/%s>", leg, site, hole)})
	if err != nil {
		log.Print(err)
	}

	res, err := repo.Query(q)
	if err != nil {
		log.Print(err)
	}

    maxage := 0.0
    solutionsTest := res.Solutions() // map[string][]rdf.Term
	  for _, i := range solutionsTest {
  		age, err := strconv.ParseFloat(fmt.Sprint(i["age"]), 64)
          if err != nil {
		log.Print(err)
	}
        if age > maxage {
            maxage = age
        }
	}
    
	fmt.Println("res.Solutions():")
	for k, i := range solutionsTest {
		fmt.Printf("At postion %v with %v \n", k, i)
	}
  
      
	return maxage

}

func AgeModel(leg string, site string, hole string)  *sparql.Results {
	repo, err := sparql.NewRepo("http://data.oceandrilling.org/sparql",
		//sparql.DigestAuth("dba", "dba"),
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Print(err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("agemodel",struct{ URI string }{fmt.Sprintf("<http://data.oceandrilling.org/codices/lsh/%s/%s/%s>", leg, site, hole)})
	if err != nil {
		log.Print(err)
	}

	res, err := repo.Query(q)
	if err != nil {
		log.Print(err)
	}

	// Print loop testing
	bindingsTest := res.Results.Bindings // map[string][]rdf.Term
	fmt.Println("res.Resuolts.Bindings:")
	for k, i := range bindingsTest {
		fmt.Printf("At postion %v with %v and %v\n", k, i["p"], i["o"])
	}

	bindingsTest2 := res.Bindings() // map[string][]rdf.Term
	fmt.Println("res.Bindings():")
	for k, i := range bindingsTest2 {
		fmt.Printf("At postion %v with %v \n", k, i)
	}

	solutionsTest := res.Solutions() // map[string][]rdf.Term
	fmt.Println("res.Solutions():")
	for k, i := range solutionsTest {
		fmt.Printf("At postion %v with %v \n", k, i)
	}

	return res

}


func LegSite() *sparql.Results {
	repo, err := sparql.NewRepo("http://data.oceandrilling.org/sparql",
		//sparql.DigestAuth("dba", "dba"),
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Print(err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("legcall2")
	if err != nil {
		log.Print(err)
	}

	res, err := repo.Query(q)
	if err != nil {
		log.Print(err)
	}

	// Print loop testing
	bindingsTest := res.Results.Bindings // map[string][]rdf.Term
	fmt.Println("res.Resuolts.Bindings:")
	for k, i := range bindingsTest {
		fmt.Printf("At postion %v with %v and %v\n", k, i["p"], i["o"])
	}

	bindingsTest2 := res.Bindings() // map[string][]rdf.Term
	fmt.Println("res.Bindings():")
	for k, i := range bindingsTest2 {
		fmt.Printf("At postion %v with %v \n", k, i)
	}

	solutionsTest := res.Solutions() // map[string][]rdf.Term
	fmt.Println("res.Solutions():")
	for k, i := range solutionsTest {
		fmt.Printf("At postion %v with %v \n", k, i)
	}

	return res

}

func ExpPublications(queryterm string) string {
	repo, err := sparql.NewRepo("http://data.oceandrilling.org/sparql",
		//sparql.DigestAuth("dba", "dba"),
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Print(err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	// q, err := bank.Prepare("my-query", struct{ Limit, Offset int }{10, 100})
	q, err := bank.Prepare("legcall", struct{ Leg string }{"113"})
	if err != nil {
		log.Print(err)
	}

	res, err := repo.Query(q)
	if err != nil {
		log.Print(err)
	}

	// Print loop testing
	bindingsTest := res.Results.Bindings // map[string][]rdf.Term
	fmt.Println("res.Resuolts.Bindings:")
	for k, i := range bindingsTest {
		fmt.Printf("At postion %v with %v and %v\n", k, i["pro"], i["vol"])
	}

	bindingsTest2 := res.Bindings() // map[string][]rdf.Term
	fmt.Println("res.Bindings():")
	for k, i := range bindingsTest2 {
		fmt.Printf("At postion %v with %v \n", k, i)
	}

	solutionsTest := res.Solutions() // map[string][]rdf.Term
	fmt.Println("res.Solutions():")
	for k, i := range solutionsTest {
		fmt.Printf("At postion %v with %v \n", k, i)
	}

	return q
}
