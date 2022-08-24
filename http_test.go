package main

import (
	"github.com/graphql-go/graphql"
	gqlhandler "github.com/graphql-go/graphql-go-handler"
	"testing"
)

func TestListen(t *testing.T) {
	go httpListen()

}

func TestCorsMiddleware(t *testing.T) {
	schemaConfig := graphql.SchemaConfig{}
	schema, _ := graphql.NewSchema(schemaConfig)
	h := gqlhandler.New(&gqlhandler.Config{Schema: &schema})
	CorsMiddleware(h)
}
