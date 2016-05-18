package graphql

import (
	// "encoding/json"
	// "fmt"
	// "github.com/chris-ramon/graphql"
	"gopkg.in/mgo.v2"
	"log"
	"opencoredata.org/ocdServices/connectors"
)

type user struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Measure string `json:"measure"`
	Leg     string `json:"leg"`
	Count   int    `json:"count"`
}

type MLCount struct {
	id      string `bson:"_id,omitempty"` // I don't really want the ID, so leave it lower case
	Measure string `json:"measure"`
	Leg     string `json:"leg"`
	Count   int    `json:"count"`
}

var data map[string]user

/*
   Create User object type with fields "id" and "name" by using GraphQLObjectTypeConfig:
       - Name: name of object type
       - Fields: a map of fields by using GraphQLFieldConfigMap
   Setup type of field use GraphQLFieldConfig
*/
// var UserType = graphql.NewObject(
// 	graphql.ObjectConfig{
// 		Name: "User",
// 		Fields: graphql.FieldConfigMap{
// 			"id": &graphql.FieldConfig{
// 				Type: graphql.String,
// 			},
// 			"name": &graphql.FieldConfig{
// 				Type: graphql.String,
// 			},
// 			"measure": &graphql.FieldConfig{
// 				Type: graphql.String,
// 			},
// 			"leg": &graphql.FieldConfig{
// 				Type: graphql.String,
// 			},
// 			"count": &graphql.FieldConfig{
// 				Type: graphql.String,
// 			},
// 		},
// 	},
// )

/*
   Create Query object type with fields "user" has type [userType] by using GraphQLObjectTypeConfig:
       - Name: name of object type
       - Fields: a map of fields by using GraphQLFieldConfigMap
   Setup type of field use GraphQLFieldConfig to define:
       - Type: type of field
       - Args: arguments to query with current field
       - Resolve: function to query data using params from [Args] and return value with current type
*/
// var QueryType = graphql.NewObject(
// 	graphql.ObjectConfig{
// 		Name: "Query",
// 		Fields: graphql.FieldConfigMap{
// 			"user": &graphql.FieldConfig{
// 				Type: UserType,
// 				Args: graphql.FieldConfigArgument{
// 					"id": &graphql.ArgumentConfig{
// 						Type: graphql.String,
// 					},
// 				},
// 				Resolve: func(p graphql.GQLFRParams) interface{} {
// 					_, isOK := p.Args["id"].(string)
// 					if isOK {
// 						// this := user{ID: "23", Name: "this name"}
// 						// return this
// 						jsonstring, _ := json.MarshalIndent(gridCall(), " ", "")
// 						log.Printf("%s\n", jsonstring)
// 						return jsonstring
// 					}
// 					return nil
// 				},
// 			},
// 		},
// 	})

func gridCall() []MLCount {
	session, err := connectors.GetMongoCon()
	if err != nil {
		log.Print(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("aggregation_janusMLCountv2")

	var results []MLCount
	err = c.Find(nil).All(&results)
	if err != nil {
		log.Printf("Error calling aggregation_janusMLCount : %v", err)
	}

	return results

}

// func ExecuteQuery(query string, schema graphql.Schema) *graphql.Result {
// 	params := graphql.Params{
// 		Schema:        schema,
// 		RequestString: query,
// 	}
// 	resultChannel := make(chan *graphql.Result)
// 	go graphql.Graphql(params, resultChannel)
// 	result := <-resultChannel
// 	if len(result.Errors) > 0 {
// 		fmt.Println("wrong result, unexpected errors: %v", result.Errors)
// 	}
// 	return result
// }

// var Schema, _ = graphql.NewSchema(
// 	graphql.SchemaConfig{
// 		Query: QueryType,
// 	},
// )
