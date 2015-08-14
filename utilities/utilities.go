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

// // old 2s structs below

// type MessageTest struct {
// 	Head    *HeadType
// 	Results *ResultsType
// }

// type HeadType struct {
// 	Link []string
// 	Vars []string
// }

// type ResultsType struct {
// 	Distinct bool
// 	Ordered  bool
// 	Bindings []BindingsType
// }

// type BindingsType struct {
// 	Count CountType
// 	Label LabelType
// 	Id    IdType
// }

// type CountType struct {
// 	//   Type string
// 	//   Datatype string
// 	Value string
// }

// type LabelType struct {
// 	//   Type string
// 	Value string
// }

// type IdType struct {
// 	//   Type string
// 	Value string
// }

// type s2sCollec struct {
// 	Items []*s2sEntry
// }

// type s2sEntry struct {
// 	Label string `json:"label"`
// 	Id    string `json:"id"`
// 	Count string `json:"count"`
// }

// func (kp s2sCollec) Len() int {
// 	return len(kp.Items)
// }

// Add elements to the end of a collection
// func (kp *s2sCollec) Add(kvp *s2sEntry) {
// 	if kp.Items == nil {
// 		kp.Items = make([]*s2sEntry, 0, 4)
// 	}
// 	n := len(kp.Items)
// 	if n+1 > cap(kp.Items) {
// 		s := make([]*s2sEntry, n, 2*n+1)
// 		copy(s, kp.Items)
// 		kp.Items = s
// 	}
// 	kp.Items = kp.Items[0 : n+1]
// 	kp.Items[n] = kvp
// }

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

func SparqlCSVFunc(sparql string) string {
	//if r, err := http.Get("http://data.oceandrilling.org/sparql?query=" + url.QueryEscape(sparql)); err == nil {
	//	r.Header.Add("Accept", "text/csv")
	//	b, err := ioutil.ReadAll(r.Body)
	//	r.Body.Close()
	//	if err == nil {
	//		return string(b)
	//	}
	//}
	//return ""

	client := &http.Client{}

	//resp, err := client.Get("http://data.oceandrilling.org/sparql?query=" + url.QueryEscape(sparql)); err == nil {

	req, err := http.NewRequest("GET", "http://data.oceandrilling.org/sparql?query="+url.QueryEscape(sparql), nil)
	req.Header.Add("Accept", "text/csv")
	resp, err := client.Do(req)
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err == nil {
		return string(b)
	}
	return ""
}

func SparqlJSONFunc(sparql string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://data.oceandrilling.org/sparql?query="+url.QueryEscape(sparql), nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err == nil {
		return string(b)
	}
	return ""
}

func SparqlJSONLDFunc(sparql string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://data.oceandrilling.org/sparql?query="+url.QueryEscape(sparql), nil)
	req.Header.Add("Accept", "application/ld+json")
	resp, err := client.Do(req)
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err == nil {
		return string(b)
	}
	return ""
}
