package documents

// DataCiteExp for the values needed for the DataCite entry for an expeditions
type DataCiteExp struct {
	ExpDOI          string   // Is this the ID of the expedition or something else
	ExpURI          string   // something like http://data.rvdata.us/id/cruise/TN272 for R2R
	ResourceType    string   // Field_expedition
	CreatorName     string   // Open Core Data
	CreatorDOI      string   // re3data DOI
	Title           string   // Expedition XXX on Joides Resoultion
	Abstract        string   // * Deep-AUV Magnetic and Seismic Study of the Hawaiian Jurassic Crust, Leg 1
	DateCollected   string   // ** Really a data of a specific format 2011-11-05/2011-12-17
	ContributorName string   // Joides Resolution Science Office || Continental Scientific Drilling Corrdinating Office
	RelatedDOIs     []string // 1 or more related DOI's
	Long            string   // longitude
	Lat             string   // latitude
	Publisher       string   // Rolling Deck to Repository (R2R) Program
	Version         string   // 1, 2, 3, etc
	PubYear         string   // 2016
}

// * need a short abstract or is there a length?
// ** do I really have this..  I might in some of the "trivia" files

type RDFDistinctParams struct {
	CoreCount              string
	Coredata               string
	CoreInterval           string
	CoreRecovery           string
	Drilled                string
	Hole                   string
	Initialreportvolume    string
	Leg                    string
	Legsite                string
	Logdata                string
	Penetration            string
	PercentRecovery        string
	Preliminaryreport      string
	Program                string
	Scientificprospectus   string
	Scientificreportvolume string
	Site                   string
	Waterdepth             string
	Vcddata                string
	Proceedingreport       string
	Expeditionsite         string
	Geom                   string
	Note                   string
	Prcoeedingreport       string
	Uri                    string
	Vcdata                 string
	Cruisetype             string
	Endportcall            string
	Expedition             string
	Legsitehole            string
	Operator               string
	Participant            string
	Scheduler              string
	Startportcall          string
	Track                  string
	Vessel                 string
	Lat                    string
	Long                   string
}

// XMLtemplate is a text/template source for the DataCiteExp struct
const XMLtemplate = `
<?xml version="1.0"?>
<resource xmlns="http://datacite.org/schema/kernel-3" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://datacite.org/schema/kernel-3 http://schema.datacite.org/meta/kernel-3/metadata.xsd">
  <identifier identifierType="DOI">{{.ExpDOI}}</identifier>
  <alternateIdentifiers>
     <alternateIdentifier alternateIdentifierType="URL">{{.ExpURI}}</alternateIdentifier>
  </alternateIdentifiers>
  <resourceType resourceTypeGeneral="Event">{{.ResourceType}}</resourceType>
  <creators>
    <creator>
      <creatorName>{{.CreatorName}}</creatorName>
      <nameIdentifier nameIdentifierScheme="DOI">{{.CreatorDOI}}</nameIdentifier>
    </creator>
  </creators>
  <titles>
    <title>{{.Title}}</title>
  </titles>
  <descriptions>
    <description descriptionType="Abstract">{{.Abstract}}</description>
  </descriptions>
  <dates>
    <date dateType="Collected">{{.DateCollected}}</date>
  </dates>
  <language>en</language>
  <contributors>
    <contributor contributorType="Producer">
      <contributorName>{{.ContributorName}}</contributorName>
    </contributor>
    <contributor contributorType="Funder">
      <contributorName>National Science Foundation</contributorName>
      <nameIdentifier nameIdentifierScheme="DOI">10.13039/100000001</nameIdentifier>
    </contributor>
  </contributors>
  <relatedIdentifiers>
                {{range $ITEMS := .RelatedDOIs}}
    <relatedIdentifier relatedIdentifierType="DOI" relationType="IsReferencedBy">{{.}}</relatedIdentifier>   // TODO this needs to be a loop
                {{end}}
  </relatedIdentifiers>
  <geoLocations>
   <geoLocation>
      <geoLocationPoint>{{.Long}} {{.Lat}}</geoLocationPoint>
    </geoLocation>
  </geoLocations>
  <publisher>{{.Publisher}}</publisher>
  <version>{{.Version}}</version>
  <publicationYear>{{.Year}}</publicationYear>
</resource>
`
