package spatial

import (
	"encoding/json"

	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	// "encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	gj "github.com/kpawlik/geojson"
	"github.com/mb0/wkt"
	"gopkg.in/mgo.v2/bson"
	"opencoredata.org/ocdServices/connectors"
	"opencoredata.org/ocdServices/structs"
	"opencoredata.org/ocdServices/utilities"
)

type GeoLatLong struct {
	id           string `bson:"_id,omitempty"` // I don't really want the ID, so leave it lower case
	Opencoresite string
	Spatial      Spatial
}

type Spatial struct {
	Geo Geo
}

type Geo struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type CruiseGL struct {
	Expedition    string `modeDescription:"Cruise GL Struct". description:"Exepedition docs"`
	CruiseType    string `description:"Cruise Type "`
	EndPortCall   string `description:"Ending port call"`
	Operator      string `description:"Operator"`
	Participant   string `description:"Participant"`
	Program       string `description:"Program"`
	Scheduler     string `description:"Scheduler"`
	StartPortCall string `description:"Start Port Call"`
	LegSiteHole   string `description:"Leg Site Hole"`
	Track         string `description:"Track"`
	Vessel        string `description:"Vessel"`
	Note          string `description:"Note"`
	URI           string `description:"URI"`
}

// CSDCO struct for informationrelated to CSDCO projects
// TODO  remove the _ from the structs..  should not be used...
type CSDCO struct {
	LocationName           string
	LocationType           string
	Project                string
	LocationID             string
	Site                   string
	Hole                   string
	SiteHole               string
	OriginalID             string
	HoleID                 string
	Platform               string
	Date                   string
	WaterDepthM            string
	Country                string
	State_Province         string
	County_Region          string
	PI                     string
	Lat                    string
	Long                   string
	Elevation              string
	Position               string
	StorageLocationWorking string
	StorageLocationArchive string
	SampleType             string
	Comment                string
	mblfT                  string
	mblfB                  string
	MetadataSource         string
}

// might need one for other metadata too...
func New() *restful.WebService {
	service := new(restful.WebService)

	service.
		Path("/api/v1/spatial").
		Doc("Spatial servives to Open Core Data holdings").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.GET("/datasets").
		To(DatasetCall).
		Doc("get datasets").
		Operation("DatasetCall"))

	service.Route(service.GET("/expeditions").To(AllExpeditions).
		Doc("get all expeditions").
		Operation("AllExpeditions"))

	service.Route(service.GET("/continental").To(CSDCOFeatures).
		Doc("get all continental drill sites").
		Operation("CSDCOFeatures"))

	service.Route(service.GET("/expedition/{leg}").To(Expeditions).
		Doc("get expedition by leg").
		Param(service.PathParameter("leg", "Leg in format like 123 or 312U").DataType("string")).
		Operation("Expeditions"))

	service.Route(service.GET("/expedition/{leg}/{site}").To(LegSite). // TODO  expand to work with site as well...  Also Hole?
										Doc("get expedition by leg and site").
										Param(service.PathParameter("leg", "Leg in format like 123 or 312U").DataType("string")).
										Param(service.PathParameter("site", "Site in format like 1234").DataType("string")).
										Operation("LegSite"))

	service.Route(service.GET("/search/jrso").To(WKTFeaturesJRSO). // TODO make work with WKT or GeoJSON
									Doc("get expeditions from a spatial polygon defined by wkt").
									Param(service.QueryParameter("geowithin", "Polygon in WKT format within which to look for features.  Try `POLYGON((-72.2021484375 38.58896696823242,-59.1943359375 38.58896696823242,-59.1943359375 28.11801628757283,-72.2021484375 28.11801628757283,-72.2021484375 38.58896696823242))`").DataType("string")).
									Param(service.QueryParameter("abstracts", "If set `true` then abstracts are sent, otherwise abstracts are not sent.  Default is not to send").DataType("string")).
									ReturnsError(400, "Unable to handle request", nil).
									Operation("WKTFeaturesJRSO"))

	service.Route(service.GET("/search/jrso/abstract").To(GetAbstract). // Get abstract based on ID/URI  like  "URI": "http://opencoredata.org/id/expedition/167/1018/B"
										Doc("Get abstract based on ID/URI").
										Param(service.QueryParameter("uri", "Get abstracts associated with the URI such as: http://opencoredata.org/id/expedition/167/1018/B").DataType("string")).
										ReturnsError(400, "Unable to handle request", nil).
										Operation("TBD"))

	return service
}

// GetAbstract fuunction to pull abstract based on URI request
func GetAbstract(request *restful.Request, response *restful.Response) {
	uri := request.QueryParameter("uri")

	session, err := connectors.GetMongoCon()
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	// session.SetMode(mgo.Monotonic, true)
	c := session.DB("expedire").C("featuresAbsGeoJSON") // featuresGeoJSON  and featuresAbsGeoJSON
	var result structs.ExpeditionGeoJSON

	err = c.Find(bson.M{"uri": uri}).Select(bson.M{"uri": 1, "abstract": 1}).One(&result) // pull only ONE not ALL .All(&results)
	if err != nil {
		log.Printf("Error calling for abstract by URI %v", err)
	}

	jsontext, _ := json.MarshalIndent(result, "", " ")
	response.Write([]byte(jsontext))
}

// WKTFeaturesJRSO get features for JRSO data using WKT Polygon string
func WKTFeaturesJRSO(request *restful.Request, response *restful.Response) {
	wktstring := request.QueryParameter("geowithin")
	abstracts := request.QueryParameter("abstracts")

	session, err := connectors.GetMongoCon()
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	// session.SetMode(mgo.Monotonic, true)
	c := session.DB("expedire").C("featuresAbsGeoJSON") // featuresGeoJSON  and featuresAbsGeoJSON
	var results []structs.ExpeditionGeoJSON

	parsedwkt, err := WKTPolygonToFloatArray(wktstring)
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	err = c.Find(bson.M{
		"coordinates": bson.M{
			"$geoWithin": bson.M{
				"$geometry": bson.M{"type": "Polygon", "coordinates": parsedwkt},
			},
		},
	}).All(&results)

	if err != nil {
		log.Printf("Error making spatial call to MongoDB: %v", err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	// check to see if we got anything, if not return 204, success, no content
	if len(results) == 0 {
		log.Printf("Everything OK, no conent\n")
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNoContent, "No features were found in the specified POLYGON") // put in some JSON "empty" here?
	}

	// Build the geojson section
	var (
		// fc *gj.FeatureCollection
		f  *gj.Feature
		fa []*gj.Feature
	)

	// feature with propertises
	for _, item := range results {

		// c := gj.Coordinates{}
		cd := gj.Coordinate{gj.Coord(item.Coordinates[0]), gj.Coord(item.Coordinates[1])}
		// c = append(c, cd)

		// Set prop entries
		// TODO..  swith on if item.Hole exist.....
		props := map[string]interface{}{"description": item.Expedition} //  "popupContent": item.Expedition,
		// "URL": fmt.Sprintf("<a target='_blank' href='http://opencoredata.org/id/expedition/%s/%s/%s'>%s_%s%s</a>",
		// 	item.Expedition, item.Site, item.Hole, item.Expedition, item.Site, item.Hole)}
		if item.Uri != "" {
			props["URI"] = item.Uri
		}
		if item.Hole != "" {
			props["Hole"] = item.Hole
		}
		if item.Expedition != "" {
			props["Expedition"] = item.Expedition
		}
		if item.Site != "" {
			props["Site"] = item.Site
		}
		if item.Program != "" {
			props["Program"] = item.Program
		}
		if item.Waterdepth != "" {
			props["Water Depth"] = item.Waterdepth
		}
		if item.CoreCount != "" {
			props["Core Count"] = item.CoreCount
		}
		if item.Initialreportvolume != "" {
			props["Initial report volume"] = item.Initialreportvolume
		}
		if item.Coredata != "" {
			props["Coredata"] = item.Coredata
		}
		if item.Logdata != "" {
			props["Logdata"] = item.Logdata
		}
		if item.Geom != "" {
			props["Geom"] = item.Geom
		}
		if item.Scientificprospectus != "" {
			props["Scientific prospectus"] = item.Scientificprospectus
		}
		if item.CoreRecovery != "" {
			props["Core Recovery"] = item.CoreRecovery
		}
		if item.Penetration != "" {
			props["Penetration"] = item.Penetration
		}
		if item.Scientificreportvolume != "" {
			props["Scientific report volume"] = item.Scientificreportvolume
		}
		if item.Expeditionsite != "" {
			props["Expedition site"] = item.Expeditionsite
		}
		if item.Preliminaryreport != "" {
			props["Preliminary report"] = item.Preliminaryreport
		}
		if item.CoreInterval != "" {
			props["Core Interval"] = item.CoreInterval
		}
		if item.PercentRecovery != "" {
			props["Percent Recovery"] = item.PercentRecovery
		}
		if item.Drilled != "" {
			props["Drilled"] = item.Drilled
		}
		if item.Vcdata != "" {
			props["VC data"] = item.Vcdata
		}
		if item.Note != "" {
			props["Note"] = item.Note
		}
		if item.Prcoeedingreport != "" {
			props["Prcoeeding report"] = item.Prcoeedingreport
		}
		if abstracts == "true" {
			if item.Abstract != "" {
				props["Abstracts"] = item.Abstract
			}
		}

		// newp := gj.NewMultiPoint(c)
		newp := gj.NewPoint(cd)
		f = gj.NewFeature(newp, props, nil)
		fa = append(fa, f)
	}

	fc := gj.FeatureCollection{Type: "FeatureCollection", Features: fa}
	gjstr, err := gj.Marshal(fc)
	if err != nil {
		log.Println(err)
	}

	response.Write([]byte(gjstr))
}

// WKTPolygonToFloatArray return float64 array for WKT Poly string
func WKTPolygonToFloatArray(wkt string) ([][][]float64, error) {
	// Take WKT string and parse down
	parsedString := strings.TrimSuffix(strings.TrimPrefix(wkt, "POLYGON(("), "))")
	wktarray := strings.Split(parsedString, ",")

	f := [][][]float64{}
	c := [][]float64{}

	for _, item := range wktarray {
		coordSet := strings.Split(item, " ")
		// TODO..  catch these errors..  this is bad form!  The whole function needs an error
		x, err := strconv.ParseFloat(coordSet[0], 64)
		y, err := strconv.ParseFloat(coordSet[1], 64)
		cd := []float64{x, y}
		c = append(c, cd)
		if err != nil {
			log.Println(err)
			return f, errors.New("Error in the conversion of WKT to GeoJSON Polygon to support spatial search")
		}
	}

	f = append(f, c)

	fmt.Println(f)
	return f, nil
}

// WKTPolygontoGeoJSON convert WKT to GeoJSON for Polygons.
// Ended up not using this in ocdService since mgo driver needed a bson
// strcuture for query
func WKTPolygontoGeoJSON(wkt string) string {
	var (
		// fc *gj.FeatureCollection
		//f  *gj.Feature
		//fa []*gj.Feature
		newp *gj.Polygon
	)

	// Take WKT string and parse down
	parsedString := strings.TrimSuffix(strings.TrimPrefix(wkt, "POLYGON(("), "))")
	wktarray := strings.Split(parsedString, ",")

	c := gj.Coordinates{}
	for _, item := range wktarray {
		coordSet := strings.Split(item, " ")
		// TODO..  catch these errors..  this is bad form!  The whole function needs an error
		x, _ := strconv.ParseFloat(coordSet[0], 64)
		y, _ := strconv.ParseFloat(coordSet[1], 64)
		cd := gj.Coordinate{gj.Coord(x), gj.Coord(y)}
		c = append(c, cd)

	}
	newml := gj.MultiLine{c}
	newp = gj.NewPolygon(newml)
	//f = gj.NewFeature(newp, nil, nil)
	//fa = append(fa, f)

	//fc := gj.FeatureCollection{Type: "FeatureCollection", Features: fa}
	gjstr, err := gj.Marshal(newp)
	if err != nil {
		//panic(err)
		log.Printf("Error event: %v \n", err)
	}

	fmt.Println(gjstr)

	return gjstr
}

func CSDCOFeatures(request *restful.Request, response *restful.Response) {
	session, err := connectors.GetMongoCon()
	if err != nil {
		//panic(err)
		log.Printf("Error event: %v \n", err)
	}
	defer session.Close()

	// session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("csdco")
	var results []CSDCO
	// c.Find(bson.M{"_id": bson.ObjectIdHex(p.ByName("id"))}).Select(bson.M{"races": bson.M{"$elemMatch": bson.M{"_id": bson.ObjectIdHex(p.ByName("raceId"))}}}).One(&result)
	err = c.Find(bson.M{}).Select(bson.M{"lat": 1, "long": 1, "holeid": 1, "project": 1}).All(&results) // pull only Lat Long Project and HoleID
	if err != nil {
		log.Printf("Error calling for all expeditions: %v", err)
	}

	// Build the geojson section
	var (
		// fc *gj.FeatureCollection
		f  *gj.Feature
		fa []*gj.Feature
	)

	// feature with propertises
	for _, item := range results {

		// c := gj.Coordinates{}
		// TODO..  catch these errors..  this is bad form!
		x, _ := strconv.ParseFloat(item.Long, 64)
		y, _ := strconv.ParseFloat(item.Lat, 64)
		cd := gj.Coordinate{gj.Coord(x), gj.Coord(y)}
		// c = append(c, cd)

		// Turned this off and the associated for loop below..   really this could be a MongoDB aggregation on the server
		// to enable this to be far faster.
		// Grab the schema from an expedition here
		// schemameta := GetFeatures(item.Expedition, "")

		// Set prop entries
		props := map[string]interface{}{"project": item.Project, "URI": fmt.Sprintf("http://opencoredata.org/collections/csdco/%s", item.HoleID)}
		//for key, ds := range schemameta {
		//	props[fmt.Sprintf("HREF_%d", key)] = ds.Uri
		//}

		// newp := gj.NewMultiPoint(c)
		newp := gj.NewPoint(cd)
		f = gj.NewFeature(newp, props, nil)
		fa = append(fa, f)
	}

	fc := gj.FeatureCollection{Type: "FeatureCollection", Features: fa}
	gjstr, err := gj.Marshal(fc)
	if err != nil {
		//panic(err)
		log.Printf("Error event: %v \n", err)
	}
	response.Write([]byte(gjstr))
}

func AllExpeditions(request *restful.Request, response *restful.Response) {
	session, err := connectors.GetMongoCon()
	if err != nil {
		//panic(err)
		log.Printf("Error event: %v \n", err)
	}
	defer session.Close()

	// session.SetMode(mgo.Monotonic, true)
	c := session.DB("expedire").C("expeditions")
	var results []CruiseGL
	err = c.Find(bson.M{}).All(&results) // change to only those with WKT entry?  There are some without
	if err != nil {
		log.Printf("Error calling for all expeditions: %v", err)
	}

	// Build the geojson section
	var (
		// fc *gj.FeatureCollection
		f  *gj.Feature
		fa []*gj.Feature
	)

	// feature with propertises
	for _, item := range results {
		track := item.Track
		if track != "" {
			geom, err := wkt.Parse([]byte(track))
			if err != nil {
				log.Printf("ERROR:  %v", err)
			}
			c := gj.Coordinates{}
			mp := geom.(*wkt.MultiPoint)
			for _, coord := range mp.Coords {
				cd := gj.Coordinate{gj.Coord(coord.X), gj.Coord(coord.Y)}
				c = append(c, cd)
			}

			// Turned this off and the associated for loop below..   really this could be a MongoDB aggregation on the server
			// to enable this to be far faster.
			// Grab the schema from an expedition here
			// schemameta := GetFeatures(item.Expedition, "")

			// Set prop entries
			props := map[string]interface{}{"popupContent": item.URI, "URI": item.URI}
			// "end_age":"0.0", "begin_age":fmt.Sprintf("%.2f", begin_age), "feature_type": "gpml:UnclassifiedFeature",
			//for key, ds := range schemameta {
			//	props[fmt.Sprintf("HREF_%d", key)] = ds.Uri
			//}

			newp := gj.NewMultiPoint(c)
			f = gj.NewFeature(newp, props, nil)
			// The following switch statement is nice.  It allows point and linestring
			// however, the order in which the points go in is vital and gets messed up
			// in this current implemntation.  Do either to arrays not having defined
			// order or either the items being entered poorly in the original WKT string
			// case switch on c length (1 = Point, >1 = LineString)
			// i := len(c)
			// switch {
			// case i == 1:
			// 	newp := gj.NewPoint(c[0])
			// 	f = gj.NewFeature(newp, props, nil)
			// case i >= 2:
			// 	newp := gj.NewLineString(c)
			// 	f = gj.NewFeature(newp, props, nil)
			// }

			fa = append(fa, f)
		}
	}

	fc := gj.FeatureCollection{Type: "FeatureCollection", Features: fa}
	// fc.Features = fa

	gjstr, err := gj.Marshal(fc)
	if err != nil {
		//panic(err)
		log.Printf("Error event: %v \n", err)
	}

	response.Write([]byte(gjstr))
}

func Expeditions(request *restful.Request, response *restful.Response) {
	session, err := connectors.GetMongoCon()
	if err != nil {
		//panic(err)
		log.Printf("Error event: %v \n", err)
	}
	defer session.Close()

	// session.SetMode(mgo.Monotonic, true)
	c := session.DB("expedire").C("features")

	// TODO  case switch here what we have LEG and or SITE
	var results []structs.Expedition
	legRequest := request.PathParameter("leg")
	err = c.Find(bson.M{"expedition": legRequest}).All(&results) // change to only those with WKT entry?  There are some without
	if err != nil {
		log.Printf("Error calling for all expeditions: %v", err)
	}

	// Grab the schema from an expedition here
	datasets := GetSchema(request.PathParameter("leg"), "")

	// Build the geojson section
	var (
		// fc *gj.FeatureCollection
		f  *gj.Feature
		fa []*gj.Feature
	)

	// feature with propertises
	for _, item := range results {
		track := item.Geom
		if track != "" {
			geom, err := wkt.Parse([]byte(track))
			if err != nil {
				log.Printf("ERROR:  %v", err)
			}
			c := gj.Coordinate{}
			mp := geom.(*wkt.Point)
			// for _, coord := range mp.Coords {
			// c = gj.Coordinate{mp.Coord(coord.X), mp.Coord(coord.Y)}
			// 				cd := gj.Coordinate{gj.Coord(coord.X), gj.Coord(coord.Y)}
			c = gj.Coordinate{gj.Coord(mp.X), gj.Coord(mp.Y)} //  mp.Coord(coord.X), mp.Coord(coord.Y)}
			// c = append(c, cd)
			// }

			// Set prop entries
			props := map[string]interface{}{"expedition": item.Expedition, "URL": fmt.Sprintf("http://opencoredata.org/id/expedition/%s/%s", item.Expedition, item.Site)}
			for key, ds := range datasets {
				props[fmt.Sprintf("dataset%d", key)] = ds.Name
			}

			newp := gj.NewPoint(c)
			f = gj.NewFeature(newp, props, nil)

			fa = append(fa, f)
		}
	}

	fc := gj.FeatureCollection{Type: "FeatureCollection", Features: fa}

	gjstr, err := gj.Marshal(fc)
	if err != nil {
		//panic(err)
		log.Printf("Error event: %v \n", err)
	}

	response.Write([]byte(gjstr))
}

func LegSite(request *restful.Request, response *restful.Response) {
	session, err := connectors.GetMongoCon()
	if err != nil {
		//panic(err)
		log.Printf("Error event: %v \n", err)
	}
	defer session.Close()

	// session.SetMode(mgo.Monotonic, true)
	c := session.DB("expedire").C("features")

	// TODO  case switch here what we have LEG and or SITE
	legRequest := request.PathParameter("leg")
	siteRequest := request.PathParameter("site")
	var results []structs.Expedition

	err = c.Find(bson.M{"expedition": legRequest, "site": siteRequest}).All(&results) // change to only those with WKT entry?  There are some without
	if err != nil {
		log.Printf("Error calling fo expedition:%v, error is: %v", siteRequest, err)
	}

	// Grab the datasets from an expedition here
	datasets := GetSchema(request.PathParameter("leg"), siteRequest)

	// Build the geojson section
	var (
		// fc *gj.FeatureCollection
		f  *gj.Feature
		fa []*gj.Feature
	)

	// feature with propertises
	for _, item := range results {
		track := item.Geom
		if track != "" {
			geom, err := wkt.Parse([]byte(track))
			if err != nil {
				log.Printf("ERROR:  %v", err)
			}
			c := gj.Coordinate{}
			mp := geom.(*wkt.Point)
			// for _, coord := range mp.Coords {
			// c = gj.Coordinate{mp.Coord(coord.X), mp.Coord(coord.Y)}
			// 				cd := gj.Coordinate{gj.Coord(coord.X), gj.Coord(coord.Y)}
			c = gj.Coordinate{gj.Coord(mp.X), gj.Coord(mp.Y)} //  mp.Coord(coord.X), mp.Coord(coord.Y)}
			// c = append(c, cd)
			// }

			// Set prop entries
			// TODO  use function from agemodel package
			begin_age := utilities.MaxAge(item.Expedition, item.Site, item.Hole)
			props := map[string]interface{}{"end_age": "0.0", "begin_age": fmt.Sprintf("%.2f", begin_age),
				"feature_type": "gpml:UnclassifiedFeature", "name": item.Uri, "popupContent": item.Uri,
				"Site": item.Site, "Hole": item.Hole, "URL": item.Uri}
			//props := map[string]interface{}{"popupContent": item.Uri, "Site": item.Site, "Hole": item.Hole, "URL": item.Uri}
			for key, ds := range datasets {
				props[fmt.Sprintf("dataset%d", key)] = ds.Name
			}

			newp := gj.NewPoint(c)
			f = gj.NewFeature(newp, props, nil)

			fa = append(fa, f)
		}
	}

	fc := gj.FeatureCollection{Type: "FeatureCollection", Features: fa}
	// fc.Features = fa

	gjstr, err := gj.Marshal(fc)
	if err != nil {
		//panic(err)
		log.Printf("Error event: %v \n", err)
	}

	response.Write([]byte(gjstr))
}

// GetFeatures fucntion is also used in ocdWeb (really need that ocdCommons)
// need to return an error too
// TODO  this distinct is NOT working
func GetFeatures(Leg string, Site string) []structs.Expedition {
	session, err := connectors.GetMongoCon()
	if err != nil {
		//panic(err)
		log.Printf("Error event: %v \n", err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("expedire").C("features")

	var results []structs.Expedition

	// NOTE:  switched to distinct return on this
	switch Site {
	case "":
		err = c.Find(bson.M{"expedition": Leg}).All(&results)
		// err = c.Find(bson.M{"expedition": Leg}).Distinct("Uri", &results)
	default:
		err = c.Find(bson.M{"expedition": Leg, "site": Site}).All(&results)
		// err = c.Find(bson.M{"expedition": Leg, "site": Site}).Distinct("Uri", &results)
	}

	if err != nil {
		log.Printf("Error calling for ShowExpeditions: %v", err)
		results = nil
	}

	return results
}

// GetSchema pulls back an array of the schema.org struct for
// all the leg and site supplied
func GetSchema(Leg string, Site string) []structs.SchemaOrgMetadata {
	session, err := connectors.GetMongoCon()
	if err != nil {
		//panic(err)
		log.Printf("Error event: %v \n", err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("schemaorg")

	var results []structs.SchemaOrgMetadata

	switch Site {
	case "":
		err = c.Find(bson.M{"opencoreleg": Leg}).All(&results)
	default:
		err = c.Find(bson.M{"opencoreleg": Leg, "opencoresite": Site}).All(&results)
	}

	if err != nil {
		log.Printf("Error calling for ShowExpeditions: %v", err)
		results = nil
	}

	return results
}

//  WORTHLESS, NOT USED..  DELETE THIS FUNCTION
// really used by a "feature call", do I have that?
// No but for a set of Leg "sites" I can still return the max age as a public interest item
// func GetMaxAge(Leg string, Site string) float64 {
// 	// see if an age model file exist (based on keywords search for that measurement)
// 	// if it does, pull it, get for max value in the age model column.
// 	// return that
// 	session, err := connectors.GetMongoCon()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()

// 	// classic problem with CSVW...   need to call schema first to get the files I can call on and index
// 	//  I guess just call the GetSchema above and then iterrate on the result names and hold the
// 	// oldest age..

// 	// Optional. Switch the session to a monotonic behavior.
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB("test").C("csvxmeta")

// 	var results structs.CSVWMeta

// 	switch Site {
// 	case "":
// 		err = c.Find(bson.M{"opencoreleg": Leg}).All(&results)
// 	default:
// 		err = c.Find(bson.M{"opencoreleg": Leg, "opencoresite": Site}).All(&results)
// 	}

// 	if err != nil {
// 		log.Printf("Error calling for ShowExpeditions: %v", err)
// 	}

// 	return 64.5
// }

// CHANGE..  nobody wants the location of datasets..
// CHANGE this to location of features to compliment the cruises above
// Calls into schema.org for all entries (these would be type dataset)
func DatasetCall(request *restful.Request, response *restful.Response) {
	session, err := connectors.GetMongoCon()
	if err != nil {
		//panic(err)
		log.Printf("Error event: %v \n", err)
	}
	defer session.Close()

	c := session.DB("test").C("schemaorg")
	latlong := []GeoLatLong{}
	//  {"spatial.geo.latitude": 1, "spatial.geo.longitude": 1}
	err = c.Find(nil).Select(bson.M{"opencoresite": 1, "spatial.geo.latitude": 1, "spatial.geo.longitude": 1}).All(&latlong) //.Distinct("Opencoresite", &latlong)
	// err = c.Find(nil).Distinct("opencoresite", &latlong)

	// err = c.Find(nil).All(&latlong)
	if err != nil {
		log.Printf("URL lookup error: %v", err)
	}
	// jsonldtext, _ := json.MarshalIndent(latlong, "", " ") // results as embeddale JSON-LD

	// Build the geojson section
	var (
		// fc *gj.FeatureCollection
		f  *gj.Feature
		fa []*gj.Feature
	)

	// feature with propertises
	for _, item := range latlong {
		lat, err := strconv.ParseFloat(item.Spatial.Geo.Latitude, 64)
		long, err := strconv.ParseFloat(item.Spatial.Geo.Longitude, 64)
		if err != nil {
			//panic(err)
			log.Printf("Error event: %v \n", err)
		}
		p := gj.NewPoint(gj.Coordinate{gj.Coord(long), gj.Coord(lat)})
		props := map[string]interface{}{"Site": item.Opencoresite}
		f = gj.NewFeature(p, props, nil)
		fa = append(fa, f)
	}

	fc := gj.FeatureCollection{Type: "FeatureCollection", Features: fa}
	// fc.Features = fa

	gjstr, err := gj.Marshal(fc)
	if err != nil {
		//panic(err)
		log.Printf("Error event: %v \n", err)
	}

	response.Write([]byte(gjstr))
}
