package documents

import (
	"github.com/emicklei/go-restful"
	//"gopkg.in/mgo.v2"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"opencoredata.org/ocdServices/connectors"
)

// might need one for other metadata too...
func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/documents").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	service.Route(service.GET("/download/{filename}").To(GetFileByName))
	service.Route(service.GET("/download/{UUID}/{format}").To(GetFileByUUID))
	return service
}

func GetFileByName(request *restful.Request, response *restful.Response) {
	filename := request.PathParameter("filename")

	// call mongo and lookup the redirection to use...
	session, err := connectors.GetMongoCon()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	mongodb := session.DB("test")

	// var buf []byte
	file, _ := mongodb.GridFS("fs").Open(filename)
	buf := make([]byte, file.Size())
	file2, err := file.Read(buf)

	// err = c.Find(bson.M{"measure": vars["measurements"], "leg": vars["leg"]}).One(&results)
	if err != nil {
		log.Printf("Error calling aggregation_janusURLSet : %v  length %d", err, file2)
	}

	log.Printf("\n%s", string(buf))

	response.AddHeader("Content-Disposition", "inline; filename=\"myfile.txt\"")
	response.Write(buf)
}

func GetFileByUUID(request *restful.Request, response *restful.Response) {

	UUID := request.PathParameter("UUID")
	format := request.PathParameter("format")
	URI := fmt.Sprintf("http://opencoredata.org/id/dataset/%s", UUID)

	// call mongo and lookup the redirection to use...
	session, err := connectors.GetMongoCon()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// case switch this area  (scoping a response for each one) on CSVMETA, SCHEMAORG, JSON  (and CSV?)
	switch format {
	case "CSV":
		c := session.DB("test").C("schemaorg")
		result := SchemaOrgMetadata{} // need this struct  (it's everywhere.   what can I do about that?   move only ocdServices from ocdWeb?)
		err = c.Find(bson.M{"url": URI}).One(&result)
		if err != nil {
			log.Printf("URL lookup error: %v", err)
		}

		filename := result.Name // the file name
		//  can I just 303 now to the download?  Perhaps I shouldn't in case some client don't follow
		mongodb := session.DB("test")
		file, _ := mongodb.GridFS("fs").Open(filename)
		buf := make([]byte, file.Size())
		file2, err := file.Read(buf)
		if err != nil {
			log.Printf("Error calling aggregation_janusURLSet : %v  length %d", err, file2)
		}

		response.AddHeader("Content-Disposition", "inline; filename=\"myfile.txt\"")
		response.Write(buf)
	case "JSON":
		c := session.DB("test").C("schemaorg")
		result := SchemaOrgMetadata{} // need this struct  (it's everywhere.   what can I do about that?   move only ocdServices from ocdWeb?)
		err = c.Find(bson.M{"url": URI}).One(&result)
		if err != nil {
			log.Printf("URL lookup error: %v", err)
		}
		jsonldtext, _ := json.MarshalIndent(result, "", " ") // results as embeddale JSON-LD
		if err != nil {
			log.Printf("Error calling in GetFileBuyUUID : %v ", err)
		}
		response.Write(jsonldtext)
	case "CSVW":
		c := session.DB("test").C("csvwmeta")
		result := CSVWMeta{} // need this struct  (it's everywhere.   what can I do about that?   move only ocdServices from ocdWeb?)
		err = c.Find(bson.M{"url": URI}).One(&result)
		if err != nil {
			log.Printf("URL lookup error: %v", err)
		}
		jsonldtext, _ := json.MarshalIndent(result, "", " ") // results as embeddale JSON-LD
		if err != nil {
			log.Printf("Error calling in GetFileBuyUUID : %v ", err)
		}
		response.Write(jsonldtext)
	}

}
