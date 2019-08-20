package graph

import (
	"github.com/emicklei/go-restful"
)

// New function for the service calls for graph
func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/beta/graph").
		Doc("BETA: graph query").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.GET("/csdco/details/{id}").To(ProjDetails).
		Doc("BETA:  csdco project graph call").
		Param(service.PathParameter("id", "ID or something like that").DataType("string")).
		Operation("CSDCOProjs"))

	service.Route(service.GET("/csdco/search").To(CSDCOGraph).
		Doc("BETA:  csdco project graph call").
		Param(service.PathParameter("id", "ID or something like that").DataType("string")).
		Operation("CSDCOProjs"))

	return service
}

func ProjDetails(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	data := GetCSDCOProj(id)
	response.WriteEntity(data)
}

func CSDCOGraph(request *restful.Request, response *restful.Response) {
	q := request.QueryParameter("q")

	sr := CSDCOGraphCall(q)

	response.AddHeader("Content-Type", "application/json")
	response.Write(sr)
}
