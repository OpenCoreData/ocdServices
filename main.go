package main

//http://ernestmicklei.com/2012/11/24/go-restful-first-working-example/

import (
	// "github.com/chris-ramon/graphql-go/types"
	"github.com/emicklei/go-restful"
	"github.com/sogko/graphql-go-handler"
	"log"
	"net/http"
	"opencoredata.org/ocdServices/documents"
	"opencoredata.org/ocdServices/expeditions"
	"opencoredata.org/ocdServices/graphql"
	"opencoredata.org/ocdServices/neptune"
)

func main() {
	// Graphql section
	h := gqlhandler.New(&gqlhandler.Config{
		Schema: &graphql.Schema,
		Pretty: true,
	})

	// serve a GraphQL endpoint at `/graphql`
	http.Handle("/graphql", h)

	go func() {
		// http.ListenAndServe("localhost:8081", serverMuxA)
		http.ListenAndServe(":7890", nil)
	}()

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

	log.Printf("Listening on localhost:6789")
	log.Printf("Serving graphql on localhost:7890/graphql")

	server := &http.Server{Addr: ":6789", Handler: wsContainer}
	server.ListenAndServe()

}
