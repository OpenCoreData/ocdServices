package main

//http://ernestmicklei.com/2012/11/24/go-restful-first-working-example/

import (
	// "github.com/chris-ramon/graphql-go/types"
	// "encoding/json"
	"github.com/emicklei/go-restful"
	// "github.com/sogko/graphql-go-handler"
	// "github.com/chris-ramon/graphql"
	"log"
	"net/http"
	"opencoredata.org/ocdServices/documents"
	"opencoredata.org/ocdServices/expeditions"
	// ocdGraphql "opencoredata.org/ocdServices/graphql"
	"opencoredata.org/ocdServices/neptune"
	"opencoredata.org/ocdServices/spatial"
)

func main() {
	// Graphql section

	// http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
	// 	result := ocdGraphql.ExecuteQuery(r.URL.Query()["query"][0], ocdGraphql.Schema)
	// 	json.NewEncoder(w).Encode(result)
	// })

	// go func() {
	// 	// http.ListenAndServe("localhost:8081", serverMuxA)
	// 	http.ListenAndServe(":7890", nil)
	// }()

	// end graphql section

	// REST section
	wsContainer := restful.NewContainer()
	// u := UserResource{}
	// u.RegisterTo(wsContainer)

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type"},
		CookiesAllowed: false,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)
	// Add container filter to respond to OPTIONS
	wsContainer.Filter(wsContainer.OPTIONSFilter)

	wsContainer.Add(neptune.New())
	wsContainer.Add(expeditions.New())
	wsContainer.Add(expeditions.NewNG())
	wsContainer.Add(documents.New())
	wsContainer.Add(spatial.New())

	log.Printf("Services on localhost:6789")
	log.Printf("Serving graphql and HTML on localhost:7890/graphql")

	server := &http.Server{Addr: ":6789", Handler: wsContainer}
	server.ListenAndServe()

}
