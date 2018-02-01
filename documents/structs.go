package documents

// W3c csvw metadata structs
type CSVWMeta struct {
	Context      string       `json:"@context"`
	Dc_license   Dc_license   `json:"dc:license"`
	Dc_modified  Dc_modified  `json:"dc:modified"`
	Dc_publisher Dc_publisher `json:"dc:publisher"`
	Dc_title     string       `json:"dc:title"`
	Dcat_keyword []string     `json:"dcat:keyword"`
	TableSchema  TableSchema  `json:"tableSchema"`
	URL          string       `json:"url"`
}

type Dc_license struct {
	Id string `json:"@id"`
}

type Dc_modified struct {
	Type  string `json:"@type"`
	Value string `json:"@value"`
}

type Dc_publisher struct {
	Schema_name string     `json:"schema:name"`
	Schema_url  Schema_url `json:"schema:url"`
}

type Schema_url struct {
	Id string `json:"@id"`
}

type TableSchema struct {
	AboutURL   string    `json:"aboutUrl"`
	Columns    []Columns `json:"columns"`
	PrimaryKey string    `json:"primaryKey"`
}

type Columns struct {
	Datatype       string   `json:"datatype"`
	Dc_description string   `json:"dc:description"`
	Name           string   `json:"name"`
	Required       bool     `json:"required"`
	Titles         []string `json:"titles"`
}

// schema.org Dataset metadata structs
type SchemaOrgMetadata struct {
	Context             Context      `json:"@context"` // was []interface{}  should be Context struct (which has 3 items in it for each voc)
	ID                  string       `json:"@id"`
	Type                string       `json:"@type"`
	Author              Author       `json:"author"`
	Description         string       `json:"description"`
	Distribution        Distribution `json:"distribution"`
	GlviewDataset       string       `json:"geolink:dataset"`
	GlviewKeywords      string       `json:"geolink:keywords"`
	OpenCoreLeg         string       `json:"opencore:leg"`
	OpenCoreSite        string       `json:"opencore:site"`
	OpenCoreHole        string       `json:"opencore:hole"`
	OpenCoreMeasurement string       `json:"opencore:measurement"`
	Keywords            string       `json:"keywords"`
	Name                string       `json:"name"`
	Spatial             Spatial      `json:"spatial"`
	URL                 string       `json:"url"`
}

type Context struct {
	Schema   string `json:"@vocab"`
	GeoLink  string `json:"geolink"`  // namespace prefix in the rest of the struct
	OpenCore string `json:"opencore"` // namespace prefix in the rest of the struct
}

type Author struct {
	Type        string `json:"@type"`
	ID          string `json:"@id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	URL         string `json:"url"`
}

type Distribution struct {
	Type           string `json:"@type"`
	ID             string `json:"@id"`
	ContentURL     string `json:"contentUrl"`
	DatePublished  string `json:"datePublished"`
	EncodingFormat string `json:"encodingFormat"`
	InLanguage     string `json:"inLanguage"`
}

type Spatial struct {
	Type string `json:"@type"`
	ID   string `json:"@id"`
	Geo  Geo    `json:"geo"`
}

type Geo struct {
	Type      string `json:"@type"`
	ID        string `json:"@id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
