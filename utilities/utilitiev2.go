package utilities

import (
	"bytes"
	"fmt"
	sparql "github.com/knakk/sparql"
	"log"
	"time"
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

`

func LegSite() *sparql.Results {
	repo, err := sparql.NewRepo("http://data.oceandrilling.org/sparql",
		//sparql.DigestAuth("dba", "dba"),
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Fatal(err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("legcall2")
	if err != nil {
		log.Fatal(err)
	}

	res, err := repo.Query(q)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	// q, err := bank.Prepare("my-query", struct{ Limit, Offset int }{10, 100})
	q, err := bank.Prepare("legcall", struct{ Leg string }{"113"})
	if err != nil {
		log.Fatal(err)
	}

	res, err := repo.Query(q)
	if err != nil {
		log.Fatal(err)
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
