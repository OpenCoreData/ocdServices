package janus

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/emicklei/go-restful"
	"github.com/jmoiron/sqlx"
	"opencoredata.org/ocdServices/internal/connectors"
)

// OCDAgeModelSQL proxy function for DataByInterface
func OCDAgeModelSQL(request *restful.Request, response *restful.Response) {
	queryToCall := "ocd_age_model_sql"
	DataByInterface(request, response, queryToCall)
}

// DataByInterface function for evaluate janus calls
func DataByInterface(request *restful.Request, response *restful.Response, queryToCall string) {
	conn, err := connectors.GetJanusConX() // get the Oracle connection
	if err != nil {
		log.Print(err)
		http.Error(response, err.Error(), 500)
		return
	}

	defer conn.Close()

	// Need a struct for the input elements
	type lshStruct struct {
		Query       string
		Leg         string
		Site        string
		Hole        string
		Core        string
		Section     string
		DepthTop    string
		DepthBottom string
	}

	// need to find the name of the last string element
	fmt.Printf("URL: %s \n", request.Request.URL.String())

	lshParams := lshStruct{Query: queryToCall, //request.PathParameter("query"),
		Leg:         request.QueryParameter("leg"),
		Site:        request.QueryParameter("site"),
		Hole:        request.QueryParameter("hole"),
		Core:        request.QueryParameter("core"),
		Section:     request.QueryParameter("section"),
		DepthTop:    request.QueryParameter("depthtop"),
		DepthBottom: request.QueryParameter("depthbottom")}

	// get the template and populate
	lshqry, err := GetSQLString(lshParams.Query)
	if err != nil {
		log.Printf("Can not find query in map, need to return error code to use for this: %s", err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	ct, err := template.New("RDF template").Parse(lshqry)
	if err != nil {
		log.Printf("Template creation failed for query: %s", err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	var buff = bytes.NewBufferString("")
	err = ct.Execute(buff, lshParams)
	if err != nil {
		log.Printf("Template execution failed: %s", err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	// Need to log this to a "query" log so I can know what is called and how often
	//log.Print(string(buff.Bytes()))

	// Call to interface based version
	csvContent, err := GetJanusCSV(conn, string(buff.Bytes()))
	if err != nil {
		log.Print(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	if csvContent == "" {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNoContent, "No data found for request")
		return
	}

	response.Header().Set("Content-Type", "text/csv") // setting the content type header to text/csv
	response.Header().Set("Content-Disposition", "attachment;filename=ocdDataDownload.csv")
	response.Write([]byte(fmt.Sprintf("%v", csvContent)))
}

// GetJanusCSV is a test function for better formatting style....
// Keep as a called function to allow use in other code (like ocdBulk) and to future use like gRPC
// func TestFuncx(db *sqlx.DB) (*[]AgeModelx, error) {
func GetJanusCSV(db *sqlx.DB, sqlstring string) (string, error) {
	results, err := db.Queryx(sqlstring)
	if err != nil {
		log.Print(err)
		return "", err
	}
	defer results.Close()

	csvdata, _ := ResultsToCSV(results)
	return csvdata, nil
}

func GetSQLString(request string) (string, error) {

	m := make(map[string]string)

	m["ocd_age_model_sql"] = `SELECT * FROM ocd_age_model WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_age_profile_sql"] = `SELECT * FROM ocd_age_profile WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_chem_carb_sql"] = `SELECT * FROM ocd_chem_carb WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_chem_gas_sql"] = `SELECT * FROM ocd_chem_gas WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_chem_iw_sql"] = `SELECT * FROM ocd_chem_iw WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_chem_rock_eval_sql"] = `SELECT * FROM ocd_chem_rock_eval WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_core_closeup_image_sql"] = `SELECT * FROM ocd_core_closeup_image WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_core_image_sql"] = `SELECT * FROM ocd_core_image WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_core_sample_sql"] = `SELECT * FROM ocd_core_sample WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_core_sample_total_sql"] = `SELECT * FROM ocd_core_sample_total WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_core_section_summary_sql"] = `SELECT * FROM ocd_core_section_summary WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_core_summary_sql"] = `SELECT * FROM ocd_core_summary WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_cryomag_sql"] = `SELECT * FROM ocd_cryomag WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_dht_sql"] = `SELECT * FROM ocd_dht WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_digital_rgb_sql"] = `SELECT * FROM ocd_digital_rgb WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_gra_sql"] = `SELECT * FROM ocd_gra WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_icp_aes_sql"] = `SELECT * FROM ocd_icp_aes WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_lab_sampling_ref_sql"] = `SELECT * FROM ocd_lab_sampling_ref WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_leg_sampling_sql"] = `SELECT * FROM ocd_leg_sampling WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_mad_sql"] = `SELECT * FROM ocd_mad WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_misc_metadata_sql"] = `SELECT * FROM ocd_misc_metadata WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_misc_special_hole_sql"] = `SELECT * FROM ocd_misc_special_hole WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_ms2f_sql"] = `SELECT * FROM ocd_ms2f WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_mscl_sql"] = `SELECT * FROM ocd_mscl WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_msl_sql"] = `SELECT * FROM ocd_msl WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_ncr_sql"] = `SELECT * FROM ocd_ncr WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_ngr_sql"] = `SELECT * FROM ocd_ngr WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_paleo_datum_ref_sql"] = `SELECT * FROM ocd_paleo_datum_ref WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_paleo_image_sql"] = `SELECT * FROM ocd_paleo_image WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_paleo_occurrence_sql"] = `SELECT * FROM ocd_paleo_occurrence WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_paleo_sample_interp_sql"] = `SELECT * FROM ocd_paleo_sample_interp WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_paleo_taxonomy_ref_sql"] = `SELECT * FROM ocd_paleo_taxonomy_ref WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_paleo_zone_ref_sql"] = `SELECT * FROM ocd_paleo_zone_ref WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_prime_data_image_sql"] = `SELECT * FROM ocd_prime_data_image WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_pwl_sql"] = `SELECT * FROM ocd_pwl WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_pws1_sql"] = `SELECT * FROM ocd_pws1 WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_pws2_sql"] = `SELECT * FROM ocd_pws2 WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_pws3_sql"] = `SELECT * FROM ocd_pws3 WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_rsc_sql"] = `SELECT * FROM ocd_rsc WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_shear_strength_avs_sql"] = `SELECT * FROM ocd_shear_strength_avs WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_shear_strength_pen_sql"] = `SELECT * FROM ocd_shear_strength_pen WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_shear_strength_tor_sql"] = `SELECT * FROM ocd_shear_strength_tor WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_smear_slide_sql"] = `SELECT * FROM ocd_smear_slide WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_smear_slide_comp_ref_sql"] = `SELECT * FROM ocd_smear_slide_comp_ref WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_splice_interval_sql"] = `SELECT * FROM ocd_splice_interval WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_splice_tie_sql"] = `SELECT * FROM ocd_splice_tie WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_tensor_core_sql"] = `SELECT * FROM ocd_tensor_core WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_tensor_tool_sql"] = `SELECT * FROM ocd_tensor_tool WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_thermal_conductivity_sql"] = `SELECT * FROM ocd_thermal_conductivity WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_thin_section_sql"] = `SELECT * FROM ocd_thin_section WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_thin_section_comp_ref_sql"] = `SELECT * FROM ocd_thin_section_comp_ref WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_thin_section_image_sql"] = `SELECT * FROM ocd_thin_section_image WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_vcd_hard_rock_image_sql"] = `SELECT * FROM ocd_vcd_hard_rock_image WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_vcd_image_sql"] = `SELECT * FROM ocd_vcd_image WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_vcd_structure_image_sql"] = `SELECT * FROM ocd_vcd_structure_image WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_xrd_sql"] = `SELECT * FROM ocd_xrd WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_xrd_image_sql"] = `SELECT * FROM ocd_xrd_image WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`
	m["ocd_xrf_sql"] = `SELECT * FROM ocd_xrf WHERE  {{if .Leg }} leg = {{.Leg}} {{end}} {{if .Site }} AND site = {{.Site}} {{end}} {{if .Hole }} AND hole = upper('{{.Hole}}') {{end}}`

	if val, ok := m[request]; ok {
		return val, nil
	}

	return "", errors.New("Error, map key now found.  Unable to locate query for this call")

}
