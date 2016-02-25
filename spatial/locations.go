package spatial

import (
	"github.com/emicklei/go-restful"
	//"gopkg.in/mgo.v2"
	// "encoding/json"
	// "fmt"
	gj "github.com/kpawlik/geojson"
	"github.com/mb0/wkt"
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
	URI           string
}

// might need one for other metadata too...
func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/spatial").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	service.Route(service.GET("/datasets").To(DatasetCall))
	service.Route(service.GET("/expeditions").To(Expeditions))
	service.Route(service.GET("/expedition/{leg}").To(Expeditions))
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

	// TODO  case switch here on LEG being sent and rebuild the
	// query to adjust.
	legRequest := request.PathParameter("leg")
	var results []CruiseGL
	if legRequest == "" {
		err = c.Find(bson.M{}).All(&results) // change to only those with WKT entry?  There are some without
		if err != nil {
			log.Printf("Error calling for all expeditions: %v", err)
		}
	}
	if legRequest != "" {
		err = c.Find(bson.M{"expedition": legRequest}).All(&results) // change to only those with WKT entry?  There are some without
		if err != nil {
			log.Printf("Error calling fo expedition:%v, error is: %v", legRequest, err)
		}
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
				log.Printf("items: %v  %v \n\n", coord.X, coord.Y)
				// temp stuff for coordinates
				cd := gj.Coordinate{gj.Coord(coord.X), gj.Coord(coord.Y)}
				c = append(c, cd)
			}

			// Set prop entries
			props := map[string]interface{}{"Site": item.Note, "URL": item.URI}

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
		panic(err)
	}

	response.Write([]byte(gjstr))
}

// CHANGE..  nobody wants the location of datasets..
// CHANGE this to location of features to compliment the cruises above
// Calls into schema.org for all entries (these would be type dataset)
func DatasetCall(request *restful.Request, response *restful.Response) {
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
