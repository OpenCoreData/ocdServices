package documents

import (
	"github.com/emicklei/go-restful"
	//"gopkg.in/mgo.v2"
	"bytes"
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
        Doc("Access Open Core Data Documents").
    	Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
        
	service.Route(service.GET("/download/{filename}").To(GetFileByName).
        Doc("get document by file name").
        Param(service.PathParameter("filename", "Access a dataset by its name in the data store").DataType("string")).
        Operation("Get File By Name"))	
    
    service.Route(service.GET("/download/{UUID}/{format}").To(GetFileByUUID).
        Doc("get doucment by unique ID and specified format").
        Param(service.PathParameter("UUID", "UUID of the file").DataType("string")).
        Param(service.PathParameter("format", "Requested format, one of: JSON, CSV").DataType("string")).
        Operation("Get File By UUID"))	
    
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

	// Get the file
	file, _ := mongodb.GridFS("fs").Open(filename)
	buf := make([]byte, file.Size())
	fileBuf, err := file.Read(buf)

	// Build the CSV header from the CSVW metadata  (since we don't store it in the file)
	c := mongodb.C("csvwmeta")
	result := CSVWMeta{} // need this struct  (it's everywhere.   what can I do about that?   move only ocdServices from ocdWeb?)
	err = c.Find(bson.M{"dc_title": filename}).One(&result)
	if err != nil {
		log.Printf("URL lookup error: %v", err)
	}
	var buffer bytes.Buffer
	for key, column := range result.TableSchema.Columns {
		buffer.WriteString(fmt.Sprintf("%s", column.Name))
		if key+1 < len(result.TableSchema.Columns) {
			buffer.WriteString(fmt.Sprintf("\t"))
		}
	}
	buffer.WriteString("\n")

	// err = c.Find(bson.M{"measure": vars["measurements"], "leg": vars["leg"]}).One(&results)
	if err != nil {
		log.Printf("Error in GetFileByName: %v  length %d", err, fileBuf)
	}

	fullPackage := buffer.String() + string(buf)
	log.Printf("\n%s", fullPackage)

	response.AddHeader("Content-Disposition", "inline; filename=\"ocdDataFile.csv\"")
	response.Write([]byte(fullPackage))  // this seems a convoluted way to get this point..  cating a byte.buffer and []byte
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
		fileBuf, err := file.Read(buf)
		if err != nil {
			log.Printf("Error calling aggregation_janusURLSet : %v  length %d", err, fileBuf)
		}
		response.AddHeader("Content-Disposition", "inline; filename=\"ocdDataFile.csv\"")
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
    // case "Datacite":    // ref: https://golang.org/src/encoding/xml/example_test.go
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
