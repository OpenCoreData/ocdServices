package spatial

import (
	"github.com/emicklei/go-restful"
	//"gopkg.in/mgo.v2"
	// "encoding/json"
	// "fmt"
	gj "github.com/kpawlik/geojson"
	// "github.com/mb0/wkt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"opencoredata.org/ocdServices/connectors"
	"strconv"
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
	Expedition    string
	CruiseType    string
	EndPortCall   string
	Operator      string
	Participant   string
	Program       string
	Scheduler     string
	StartPortCall string
	LegSiteHole   string
	Track         string
	Vessel        string
	Note          string
}

// might need one for other metadata too...
func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/spatial").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	service.Route(service.GET("/geojson").To(GeoJSONCall))
	service.Route(service.GET("/expeditions").To(Expeditions))
	return service
}

func Expeditions(request *restful.Request, response *restful.Response) {
	session, err := connectors.GetMongoCon()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// session.SetMode(mgo.Monotonic, true)
	c := session.DB("expedire").C("expeditions")

	var results []CruiseGL
	err = c.Find(bson.M{}).All(&results)
	if err != nil {
		log.Printf("Error calling for all expeditions: %v", err)
	}

	//geom, err := wkt.Parse([]byte(`POINT ZM(1.0 2.0 3.0 4.0)`))

	// loop on results...
	// build out the feature like below, using wkt stuff...

	// Build the geojson section
	var (
		// fc *gj.FeatureCollection
		f  *gj.Feature
		fa []*gj.Feature
	)

	// feature with propertises
	for _, item := range results {
		track := item.Track
		//newp := gj.NewLineString()

		log.Printf("%s\n\n", track)

		// geom, err := wkt.Parse([]byte(track))
		// if err != nil {
		// 	// log.Printf("%s\n\n", track)
		// 	// panic(err)
		// 	// log.Printf("%v", err)
		// }

		// log.Printf("GEOM: %v  \n\n", geom)

		// for _, point := range &geom {
		// 	// cd := gj.Coordinate{gj.Coord(geom.}
		// 	log.Printf("%v  \n\n", point)
		// }

		// temp stuff for coordinates
		c := gj.Coordinates{}
		cd := gj.Coordinate{gj.Coord(12.2), gj.Coord(45.5)}
		cd2 := gj.Coordinate{gj.Coord(22.2), gj.Coord(55.5)}

		c = append(c, cd)
		c = append(c, cd2)

		newp := gj.NewLineString(c)

		// lat, err := strconv.ParseFloat(item.Spatial.Geo.Latitude, 64)
		// long, err := strconv.ParseFloat(item.Spatial.Geo.Longitude, 64)
		// if err != nil {
		// 	panic(err)
		// }
		// p := gj.NewPoint(gj.Coordinate{gj.Coord(long), gj.Coord(lat)})

		//p := gj.NewLineString(gj.Coordinates{{1, 1}, {2.001, 3}, {4001, 1223}})

		props := map[string]interface{}{"Site": item.Note}
		f = gj.NewFeature(newp, props, nil)
		fa = append(fa, f)
	}

	fc := gj.FeatureCollection{Type: "FeatureCollection", Features: fa}
	// fc.Features = fa

	gjstr, err := gj.Marshal(fc)
	if err != nil {
		panic(err)
	}

	response.Write([]byte(gjstr))
}

// CHANGE..  nobody wants the location of datasets..
// CHANGE this to location of features to compliment the cruises above
func GeoJSONCall(request *restful.Request, response *restful.Response) {
	session, err := connectors.GetMongoCon()
	if err != nil {
		panic(err)
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
			panic(err)
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
		panic(err)
	}

	response.Write([]byte(gjstr))
}
