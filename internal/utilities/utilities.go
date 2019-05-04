package utilities

import (
	"io/ioutil"
	//"log"
	"net/http"
	"net/url"
)

type MessageTestDef struct {
	Head    *HeadTypeDef
	Results *ResultsTypeDef
}

type HeadTypeDef struct {
	Link []string
	Vars []string
}

type ResultsTypeDef struct {
	Distinct bool
	Ordered  bool
	Bindings []BindingsTypeDef
}

type BindingsTypeDef struct {
	Key   KeyType
	Value ValueType
}

type KeyType struct {
	Type     string
	Datatype string
	Value    string
}

type ValueType struct {
	Type  string
	Value string
}

// Util function for SPARQL call to Virtuoso server (JSON format request) and URL escaped SPARQL query
func SparqlCallFunc(sparql string) string {
	if r, err := http.Get("http://data.oceandrilling.org/sparql?default-graph-uri=&should-sponge=&query=" + url.QueryEscape(sparql) + "&debug=on&timeout=&format=application%2Fsparql-results%2Bjson&CXML_redir_for_subjs=&CXML_redir_for_hrefs=&save=display&fname="); err == nil {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err == nil {
			return string(b)
		}
	}
	return ""
}
