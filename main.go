package main

//http://ernestmicklei.com/2012/11/24/go-restful-first-working-example/

import (
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"opencoredata.org/ocdServices/expeditions"
	"opencoredata.org/ocdServices/neptune"
)

func main() {

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

	log.Printf("Listening on localhost:6789")

	server := &http.Server{Addr: ":6789", Handler: wsContainer}
	server.ListenAndServe()

}
