package neptune

import (
	"database/sql"
	"github.com/emicklei/go-restful"
	_ "github.com/lib/pq"
	"log"
	"opencoredata.org/ocdServices/connectors"
	"os"
)

type User struct {
	Id, Name string
}

// Need a struct for SQL responce from Neptune so I can write it back.
type NeptuneSamples struct {
	Leg               string
	Site              string
	Hole              string
	Hole_id           string
	Ocean_code        string
	Taxon_abundance   string
	Taxon             string
	Fossil_group      string
	Water_depth       string
	Sample_age_ma     string
	Sample_depth_mbsf string
	Latitude          float64
	Longitude         float64
}

func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/neptune").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	//service (restful.OPTIONSFilter)
	//service.Route(service.GET("/{delteexample}").To(FindUser))
	service.Route(service.GET("/{samples}").To(SampleCall))

	return service
}

func SampleCall(request *restful.Request, response *restful.Response) {

	// read ENV vars for password and usernams
	username := os.Getenv("CHRONOS_USERNAME")
	password := os.Getenv("CHRONOS_PASSWORD")

	// connect to Chronos DB
	db, err := connectors.ChronosConn(username, password)
	if err != nil {
		log.Printf("error connectiong: %s", err)
		return
	}
	defer db.Close()

	sqlstring := SQL_samples

	rows, err := db.Query(sqlstring)
	if err != nil {
		log.Fatal(err)
	}

	// Read out results into an array of structs
	allitems := []NeptuneSamples{}
	for rows.Next() {
		var leg, site, hole, hole_id, ocean_code, taxon_abundance, taxon, fossil_group, water_depth, sample_age_ma, sample_depth_mbsf sql.NullString
		var latitide, longitude float64
		if err := rows.Scan(&sample_age_ma, &sample_depth_mbsf, &water_depth, &leg, &site, &hole, &hole_id, &latitide, &longitude, &ocean_code, &taxon_abundance, &taxon, &fossil_group); err != nil {
			log.Fatal(err)
		}
		item := NeptuneSamples{Leg: leg.String, Site: site.String, Hole: hole.String, Hole_id: hole_id.String, Ocean_code: ocean_code.String, Taxon_abundance: taxon_abundance.String, Taxon: taxon.String, Fossil_group: fossil_group.String, Water_depth: water_depth.String, Sample_age_ma: sample_age_ma.String, Sample_depth_mbsf: sample_depth_mbsf.String, Latitude: latitide, Longitude: longitude}
		allitems = append(allitems, item)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	//  Must it be a struct?
	response.WriteEntity(allitems)
}
