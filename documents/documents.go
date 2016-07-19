package documents

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"

	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"opencoredata.org/ocdServices/connectors"
	"opencoredata.org/ocdServices/structs"
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

	service.Route(service.GET("/keyword/{keyword}").To(GetFilesByKeyword).
		Doc("get array of documents based on a keyword search").
		Param(service.PathParameter("keyword", "keyword to search on").DataType("string")).
		Operation("Get File ID set by Keyword"))

	//  service.Route(service.GET("/collection/measurement/{measurement}/jrso/lsh/{leg}/{site}/{hole}").To(GetByMesJRSO).
	// Doc("get array of documents based on a measurement in JRSO data filtered by LSH (optional)").
	// Param(service.PathParameter("keyword", "keyword to search on").DataType("string")).
	// Operation("Get File ID set by measurment and expedition information"))

	//  service.Route(service.GET("/collection/measurement/{measurement}/csdco/project/{project}").To(GetByMesCSDCO).
	// Doc("get array of documents based on a measurement in JRSO data filtered by LSH (optional)").
	// Param(service.PathParameter("keyword", "keyword to search on").DataType("string")).
	// Operation("Get File ID set by measurment and expedition information"))

	return service
}

// GetFilesByKeyword returns a set of schema.org JSON elements for datasets
// that match the keyword.  This likely should be converted to a ?q= format.
func GetFilesByKeyword(request *restful.Request, response *restful.Response) {

	keyword := request.PathParameter("keyword")

	session, err := connectors.GetMongoCon()
	if err != nil {
		log.Printf("ERROR:  %v", err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("schemaorg")

	var results []structs.SchemaOrgMetadata

	err = c.Find(bson.M{"$text": bson.M{"$search": keyword}}).All(&results)

	if err != nil {
		log.Printf("Error calling for GetFilesByKeyword: %v", err)
		results = nil
	}

	resultSet, err := json.Marshal(results)
	if err != nil {
		log.Printf("ERROR:  %v", err)
	}

	response.Write([]byte(resultSet))
}

func GetFileByName(request *restful.Request, response *restful.Response) {
	filename := request.PathParameter("filename")

	// call mongo and lookup the redirection to use...
	session, err := connectors.GetMongoCon()
	if err != nil {
		log.Print(err)
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

	// put the column heading in..   // todo..  this could be a switch in the URL that puts headers on or not.
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
	response.Write([]byte(fullPackage)) // this seems a convoluted way to get this point..  cating a byte.buffer and []byte
}

func GetFileByUUID(request *restful.Request, response *restful.Response) {

	UUID := request.PathParameter("UUID")
	format := request.PathParameter("format")
	URI := fmt.Sprintf("http://opencoredata.org/id/dataset/%s", UUID)

	// call mongo and lookup the redirection to use...
	session, err := connectors.GetMongoCon()
	if err != nil {
		log.Print(err)
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
	case "DATACITE": // case "Datacite":    // ref: https://golang.org/src/encoding/xml/example_test.go
		c := session.DB("test").C("schemaorg")
		result := SchemaOrgMetadata{} // need this struct  (it's everywhere.   what can I do about that?   move only ocdServices from ocdWeb?)
		err = c.Find(bson.M{"url": URI}).One(&result)
		if err != nil {
			log.Printf("URL lookup error: %v", err)
		}
		xmltext, _ := xml.MarshalIndent(result, "", " ") // results as XML
		if err != nil {
			log.Printf("Error calling in GetFileBuyUUID : %v ", err)
		}
		response.Write(xmltext)
	case "CSVW":
		c := session.DB("test").C("csvwmeta")
		result := CSVWMeta{} // need this struct  (it's everywhere.   what can I do about that?   move only ocdServices from ocdWeb?)
		err = c.Find(bson.M{"url": URI}).One(&result)
		if err != nil {
			log.Printf("URL lookup error: %v", err)
		}
		jsonldtext, _ := json.MarshalIndent(result, "", " ") // results as embeddale CSVW JSON-LD
		if err != nil {
			log.Printf("Error calling in GetFileBuyUUID : %v ", err)
		}
		response.Write(jsonldtext)
	}

}
