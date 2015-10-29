package graphql

import (
	"github.com/chris-ramon/graphql-go/types"
)

var queryType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name:        "Query",
	Description: "Simple query example",
	Fields: types.GraphQLFieldConfigMap{
		"latestPost": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
			Resolve: func(p types.GQLFRParams) interface{} {
				return "Hello World!"
			},
		},
		"name": &types.GraphQLFieldConfig{
			Type:        types.GraphQLString,
			Description: "The name of the human.",
			Resolve: func(p types.GQLFRParams) interface{} {
				return "first test items"
			},
		},
	},
})

var Schema, _ = types.NewGraphQLSchema(types.GraphQLSchemaConfig{
	Query: queryType,
})

func test(p types.GQLFRParams) interface{} {
	return "test return could be anything"
}
