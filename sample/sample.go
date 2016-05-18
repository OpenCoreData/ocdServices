package sample

import (
    "github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	// "encoding/json"
	"fmt"
	gj "github.com/kpawlik/geojson"
	"github.com/mb0/wkt"
	"gopkg.in/mgo.v2/bson"
	"log"
	//"opencoredata.org/ocdCommons/structs"
	"opencoredata.org/ocdServices/connectors"
	"strconv"
)


type Sample struct {
	ID  string `json:"id"`
	info1 string `json:"info1"`
}



func New() *restful.WebService {
	service := new(restful.WebService)

   	service.
		Path("/api/v1/sample").
        Doc("BETA: Services related to samples (ID's etc)").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
        
	service.Route(service.GET("/ocean/{sampleid}").
        To(SampleInfo).
        Doc("get sample information by sample ID").
        Operation("SampleCall"))


return service
    
}


func SampleInfo(request *restful.Request, response *restful.Response) {
	session, err := connectors.GetMongoCon()
	if err != nil {
		log.Print(err)
	}
	defer session.Close()
 
 // session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("csdco")
	var results []Sample
	err = c.Find(bson.M{}).All(&results) // change to only those with WKT entry?  There are some without
	if err != nil {
		log.Printf("Error calling for all expeditions: %v", err)
	}
    
}