package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	sparql "github.com/knakk/sparql"
)

const gensearch = `
# Comments are ignored, except those tagging a query.


# tag: csdcograph
PREFIX text: <http://jena.apache.org/text#>
PREFIX schema: <http://schema.org/>

select ?s ?score ?literal ?g ?type
where { 
  
 {
  (?s ?score ?literal ?g ) text:query (schema:text "{{.}}" ) .
  }
  UNION
  {
    (?s ?score ?literal ?g) text:query (schema:description "{{.}}") .
  }
  OPTIONAL {
    ?s <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> ?type .
  }
}

# tag: csdcographOLD
SELECT ?s ?p ?o ?type  ?g 
WHERE 
{  GRAPH ?g
    {
      ?s ?p ?o .
       ?s <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> ?type .
     
    }
    FILTER regex(?o, "{{.}}", "i")
}
LIMIT 100


`

// CSDCOGraphCall is from the new JS based search UI
func CSDCOGraphCall(qs string) []byte {
	repo, err := getJena()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(gensearch)
	bank := sparql.LoadBank(f)

	// Need to convert the qs string into a format
	// Jena Lucene likes.   So for now just split the string
	// and join with OR.

	qsa := strings.Split(qs, " ")
	j := strings.Join(qsa, " OR ")

	qsj := fmt.Sprintf("( %s )", j)
	log.Println(qsj)

	q, err := bank.Prepare("csdcograph", qsj)
	if err != nil {
		log.Printf("%s\n", err)
	}

	// log.WithFields(log.Fields{  // don't do log with fields in this code base yet
	// 	"SPARQL": q,
	// }).Info("A SPARQL call in CSDCO namespace")

	log.Println(q)

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	// rr := &ResourceResults
	b, err := json.MarshalIndent(res, " ", "")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// fmt.Println(string(b))

	// return fmt.Sprintf("%v", res.Bindings())
	// for this one don't return the map..  return JSON of the results
	return b
}
