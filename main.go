package main

import (
	"log"
	"time"
)

func setupServer() {
	start := time.Now()
	setupConfigurationGlobals()
	dataMap := getDataMap()

	schemaConfig := generateGraphqlObjects(dataMap)

	graphqlSchema := createSchema(schemaConfig)

	startGraphQLServer(graphqlSchema)

	setupAuthRoute(dataMap)

	elapsed := time.Since(start)
	log.Printf("Server started in %s", elapsed)
}

func main() {
	setupServer()
	httpListen()
}
