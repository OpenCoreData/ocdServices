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
		Doc("Rock evaluation").
		Param(service.QueryParameter("leg", "Leg of expedition")).
		Param(service.QueryParameter("site", "Site of expedition")).
		Param(service.QueryParameter("hole", "Hole of expedition")).
		Param(service.QueryParameter("core", "Core")).
		Param(service.QueryParameter("section", "Core section")).
		Param(service.QueryParameter("depthtop", "Depth top")).
		Param(service.QueryParameter("depthbottom", "Depth bottom")).
		Operation("Rock evaluation query"))

	service.Route(service.GET("/agemodel").To(GetAgeModel).
		Doc("Test service GetAgeModles").
		Param(service.QueryParameter("leg", "Leg of expedition")).
		Param(service.QueryParameter("site", "Site of expedition")).
		Param(service.QueryParameter("hole", "Hole of expedition")).
		Operation("GetAgeModles evaluation query"))

	service.Route(service.GET("/ocd_age_model_sql").To(OCDAgeModelSQL).
		Doc("Janus: Get age model data").
		Param(service.QueryParameter("leg", "Leg of expedition")).
		Param(service.QueryParameter("site", "Site of expedition")).
		Param(service.QueryParameter("hole", "Hole of expedition")).
		Param(service.QueryParameter("core", "Core value from hole")).
		Param(service.QueryParameter("section", "Section of core")).
		Param(service.QueryParameter("depthtop", "Top Depth (mbsf)")).
		Param(service.QueryParameter("depthbottom", "Bottom Depth (mbsf)")).
		Operation("OCD_age_model_sql"))

	return service
}
