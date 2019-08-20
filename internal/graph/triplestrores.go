package graph

import (
	"log"
	"time"

	sparql "github.com/knakk/sparql"
)

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
	// repo, err := sparql.NewRepo("http://graph.opencoredata.org/opencore/sparql",

	repo, err := sparql.NewRepo("http://localhost:3030/doa/query",
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Printf("%s\n", err)
	}
	return repo, err
}
