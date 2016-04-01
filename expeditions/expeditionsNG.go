package expeditions

import (
	"github.com/emicklei/go-restful"
	utilitiesv2 "opencoredata.org/ocdServices/utilities"
)

func NewNG() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/beta/expeditionsng").
        Doc("BETA: Not versioned, not for general use.").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
        
	service.Route(service.GET("/{leg}").To(LithCallNG).
		Doc("BETA:  expedition information by leg").
		Param(service.PathParameter("leg", "Leg in format like 123 or 312U").DataType("string")).
		Operation("LithCallNG"))
	
    
    service.Route(service.GET("/test/{leg}").To(LegSiteNG).
		Doc("BETA:  Version 2:  expedition information by leg").
		Param(service.PathParameter("leg", "Leg in format like 123 or 312U").DataType("string")).
		Operation("LegSiteNG"))

	return service
}

func LithCallNG(request *restful.Request, response *restful.Response) {

	data := utilitiesv2.ExpPublications("113")
	response.WriteEntity(data)

}

func LegSiteNG(request *restful.Request, response *restful.Response) {

	data := utilitiesv2.LegSite()
	response.WriteEntity(data)

}
