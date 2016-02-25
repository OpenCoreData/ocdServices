package expeditions

import (
	"github.com/emicklei/go-restful"
	utilitiesv2 "opencoredata.org/ocdServices/utilities"
)

func NewNG() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/expeditionsng").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	service.Route(service.GET("/{leg}").To(LithCallNG))
	service.Route(service.GET("/test/{leg}").To(LegSiteNG))

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
