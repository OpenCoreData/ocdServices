package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	sparql "github.com/knakk/sparql"
)

const queries = `
# Comments are ignored, except those tagging a query.


# tag: csdcograph
PREFIX   ex: <http://www.example.org/resources#>
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX text: <http://jena.apache.org/text#>
prefix schema: <http://schema.org/>
prefix bds: <http://www.bigdata.com/rdf/search#>
prefix csdco: <http://opencoredata.org/id/voc/csdco/v1/>
select DISTINCT ?s ?proj ?locationname ?pi ?country ?state_province ?lat ?long ?score ?rank
{
  (?s ?score) text:query (rdfs:comment '{{.}}') .
  ?s csdco:pi ?pi .
  ?s csdco:project ?proj .
  ?s csdco:locationname ?locationname .
  ?s csdco:country ?country .
  ?s csdco:state_province ?state_province .

  ?s <http://www.w3.org/2003/01/geo/wgs84_pos#lat> ?lat .
  ?s <http://www.w3.org/2003/01/geo/wgs84_pos#long> ?long .   
}
ORDER BY DESC(?score)



# tag: csdcographDEPRECATED
prefix schema: <http://schema.org/>
prefix bds: <http://www.bigdata.com/rdf/search#>
prefix csdco: <http://opencoredata.org/id/voc/csdco/v1/>
select DISTINCT ?s ?proj ?locationname ?pi ?country ?state_province ?lat ?long ?score ?rank
{
  BIND("{{.}}" AS ?q)
  ?s csdco:project ?proj .
  ?s csdco:locationname ?locationname .
  ?s csdco:pi ?pi .
  ?s csdco:country ?country .
  ?s csdco:state_province ?state_province .
  ?s <http://www.w3.org/2003/01/geo/wgs84_pos#lat> ?lat .
  ?s <http://www.w3.org/2003/01/geo/wgs84_pos#long> ?long .
  {
  ?s rdf:type csdco:CSDCOProject .
  ?s csdco:project ?sproj .
  ?sproj bds:search ?q .
  ?sproj bds:minRelevance "0.25" .
  ?sproj bds:relevance ?score .
  ?sproj bds:maxRank "1000" .
  ?sproj bds:rank ?rank .
  }
  UNION {
  ?s rdf:type csdco:CSDCOProject .
  ?s csdco:pi ?spi .
  ?spi bds:search ?q .
  ?spi bds:minRelevance "0.25" .
  ?spi bds:relevance ?score .
  ?spi bds:maxRank "1000" .
  ?spi bds:rank ?rank .
  }
  UNION {
  ?s rdf:type csdco:CSDCOProject .
  ?s csdco:country ?scountry .
  ?scountry bds:search ?q .
  ?scountry bds:minRelevance "0.25" .
  ?scountry bds:relevance ?score .
  ?scountry bds:maxRank "1000" .
  ?scountry bds:rank ?rank .
  }
  UNION {
  ?s rdf:type csdco:CSDCOProject .
  ?s csdco:state_province ?ssp .
  ?ssp bds:search ?q .
  ?ssp bds:minRelevance "0.25" .
  ?ssp bds:relevance ?score .
  ?ssp bds:maxRank "1000" .
  ?ssp bds:rank ?rank .
  }
  UNION {
  ?s rdf:type csdco:CSDCOProject .
  ?s csdco:locationname ?slocationname .
  ?slocationname bds:search ?q .
  ?slocationname bds:minRelevance "0.25" .
  ?slocationname bds:relevance ?score .
  ?slocationname bds:maxRank "1000" .
  ?slocationname bds:rank ?rank .
  }
}
ORDER BY DESC(?score)

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

#tag: csdcoproj
SELECT *
WHERE  
{    
	?uri rdf:type <http://opencoredata.org/id/voc/csdco/v1/CSDCOProject> .    
	?uri <http://opencoredata.org/id/voc/csdco/v1/holeid> "{{.ID}}" .  
	?uri ?p ?o .
}

`

func getCSDCO() (*sparql.Repo, error) {
	repo, err := sparql.NewRepo("http://opencoredata.org/blazegraph/namespace/csdco/sparql",
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Printf("%s\n", err)
	}
	return repo, err
}

func getJena() (*sparql.Repo, error) {
	repo, err := sparql.NewRepo("http://graph.opencoredata.org/opencore/sparql",
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Printf("%s\n", err)
	}
	return repo, err
}

// CSDCOGraphCall is from the new JS based search UI
func CSDCOGraphCall(qs string) []byte {
	repo, err := getJena()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("csdcograph", qs)
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

// TODO ..  error nice if I don't get what I expect in these params
func GetCSDCOProj(identity string) *sparql.Results {
	// repo, err := sparql.NewRepo("http://data.oceandrilling.org/sparql",
	repo, err := sparql.NewRepo("http://localhost:9999/blazegraph/namespace/opencore/sparql",
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Print(err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	// q, err := bank.Prepare("my-query", struct{ Limit, Offset int }{10, 100})
	q, err := bank.Prepare("csdcoproj", struct{ ID string }{identity})
	if err != nil {
		log.Print(err)
	}

	log.Println(q)

	res, err := repo.Query(q)
	if err != nil {
		log.Print(err)
	}

	return res
}

func MaxAge(leg string, site string, hole string) float64 {
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

	q, err := bank.Prepare("agemodel", struct{ URI string }{fmt.Sprintf("<http://data.oceandrilling.org/codices/lsh/%s/%s/%s>", leg, site, hole)})
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

func AgeModel(leg string, site string, hole string) *sparql.Results {
	repo, err := sparql.NewRepo("http://data.oceandrilling.org/sparql",
		//sparql.DigestAuth("dba", "dba"),
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Print(err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("agemodel", struct{ URI string }{fmt.Sprintf("<http://data.oceandrilling.org/codices/lsh/%s/%s/%s>", leg, site, hole)})
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
