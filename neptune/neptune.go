package neptune

import (
	// "bytes"
	"database/sql"
	"github.com/emicklei/go-restful"
	_ "github.com/lib/pq"
	"log"
	//"fmt"
	_ "github.com/lib/pq"
)

type User struct {
	Id, Name string
}

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

// need a structu above for my SQL responce from Neptune.  I can then
//write tihs back to the user.

func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/neptune").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	//service (restful.OPTIONSFilter)

	//service.Route(service.GET("/{delteexample}").To(FindUser))
	service.Route(service.GET("/{samples}").To(SampleCall))

	return service
}

func SampleCall(request *restful.Request, response *restful.Response) {
	db, err := sql.Open("postgres", "user=USERNAME password=PASSWORD host=localhost port=9000 dbname=neptune sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	log.Println("Now in the service call")

	// ask Pat about this querry..  and the DISTING aspect I added (why would it need to be added?)
	sqlstring := `SELECT 
                a.sample_age_ma, a.sample_depth_mbsf, c.water_depth, c.leg, c.site, c.hole, a.hole_id, 
                c.latitude, c.longitude, c.ocean_code, d.taxon_abundance, 
                TRIM(INITCAP(e.genus) || ' ' || LOWER(e.species) || ' ' || LOWER(COALESCE(e.subspecies,' '))) AS taxon, e.fossil_group 
            FROM 
                neptune_sample a, neptune_core b, neptune_hole_summary c, neptune_sample_taxa d, 
                neptune_taxonomy e  
            WHERE
                a.hole_id = b.hole_id 
                AND b.hole_id = c.hole_id 
                AND a.core = b.core 
                AND a.sample_id = d.sample_id 
                AND e.taxon_id = d.taxon_id            
            ORDER BY 
                a.sample_age_ma
	  LIMIT 1000`

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
	log.Println("Going to try and send allitems now")
	response.WriteEntity(allitems)
}

// func FindUser(request *restful.Request, response *restful.Response) {
// 	id := request.PathParameter("user-id")
// 	// here you would fetch user from some persistence system
// 	usr := &User{Id: id, Name: "John Doe"}
// 	response.WriteEntity(usr)
// }
