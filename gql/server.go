package gql

import (
	"fmt"

	"github.com/graph-gophers/graphql-go"
	"github.com/nedrocks/delphisbe/gql/resolver"
	"github.com/nedrocks/delphisbe/gql/schema"
)

func NewServer() (*graphql.Schema, error) {
	schema, err := initializeSchema()
	if err != nil {
		fmt.Printf("Error creating schema: %+v\n", err)
		return nil, err
	}
	fmt.Printf("%+v", schema)
	return schema, nil
}

func initializeSchema() (*graphql.Schema, error) {
	return graphql.MustParseSchema(schema.String(), &resolver.Resolver{}), nil
}
