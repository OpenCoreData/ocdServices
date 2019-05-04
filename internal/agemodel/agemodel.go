package agemodel

import (
	"github.com/emicklei/go-restful"
	utilitiesv2 "opencoredata.org/ocdServices/internal/utilities"
)

func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/beta/ages").
		Doc("BETA: Not versioned, not for general use.").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.GET("/agemodel/{leg}/{site}/{hole}").To(AgeModel).
		Doc("BETA:  get age model information based on Leg Site Hole").
		Param(service.PathParameter("leg", "Leg in format like 207 ").DataType("string")).
		Param(service.PathParameter("site", "Site in format like 1260 ").DataType("string")).
		Param(service.PathParameter("hole", "Hole in format like B ").DataType("string")).
		Operation("AgeModel"))

	service.Route(service.GET("/maxage/{leg}/{site}/{hole}").To(MaxAge).
		Doc("BETA:  get max age based on Leg Site Hole").
		Param(service.PathParameter("leg", "Leg in format like 207 ").DataType("string")).
		Param(service.PathParameter("site", "Site in format like 1260 ").DataType("string")).
		Param(service.PathParameter("hole", "Hole in format like B ").DataType("string")).
		Operation("MaxAge"))

	return service
}

func MaxAge(request *restful.Request, response *restful.Response) {

	legRequest := request.PathParameter("leg")
	siteRequest := request.PathParameter("site")
	holeRequest := request.PathParameter("hole")

	data := utilitiesv2.MaxAge(legRequest, siteRequest, holeRequest)
	response.WriteEntity(data)

}

func AgeModel(request *restful.Request, response *restful.Response) {

	legRequest := request.PathParameter("leg")
	siteRequest := request.PathParameter("site")
	holeRequest := request.PathParameter("hole")

	data := utilitiesv2.AgeModel(legRequest, siteRequest, holeRequest)
	response.WriteEntity(data)

}

// func MinAge

// func MaxAge

// func Locs
