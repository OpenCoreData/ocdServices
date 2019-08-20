package graph

import (
	"bytes"
	"log"

	sparql "github.com/knakk/sparql"
)

const projdetails = `
# Comments are ignored, except those tagging a query.

#tag: csdcoproj
SELECT *
WHERE  
{    
	?uri rdf:type <http://opencoredata.org/id/voc/csdco/v1/CSDCOProject> .    
	?uri <http://opencoredata.org/id/voc/csdco/v1/holeid> "{{.ID}}" .  
	?uri ?p ?o .
}

`

// GetCSDCOProj  TODO ..  error nice if I don't get what I expect in these params
func GetCSDCOProj(identity string) *sparql.Results {
	// repo, err := sparql.NewRepo("http://data.oceandrilling.org/sparql",
	repo, err := getJena()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(projdetails)
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
