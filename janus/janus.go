package janus

import (
	"github.com/emicklei/go-restful"
)

func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/janus").
		Doc("Access Open Core Data Documents").
		Consumes("application/x-www-form-urlencoded").
		Produces("text/plain")

	service.Route(service.GET("/rockeval").To(RockEval).
		Doc("Rock Eval").
		Param(service.QueryParameter("leg", "Leg of expedition")).
		Param(service.QueryParameter("site", "Site of expedition")).
		Param(service.QueryParameter("hole", "Hole of expedition")).
		Param(service.QueryParameter("core", "Core")).
		Param(service.QueryParameter("section", "Core section")).
		Param(service.QueryParameter("depthtop", "Depth top")).
		Param(service.QueryParameter("depthbottom", "Depth bottom")).
		Operation("Rock evaluation query"))

	return service
}
