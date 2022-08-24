package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/graphql-go/graphql"
	gqlhandler "github.com/graphql-go/graphql-go-handler"
)

func getGraphQlType(dataType string) *graphql.Scalar {
	switch dataType {
	case "String":
		return graphql.String
	case "Int":
		return graphql.Int
	case "Float":
		return graphql.Float
	case "Boolean":
		return graphql.Boolean
	default:
		return graphql.String

	}
}

func generateQueries(dataMap DataMap) graphql.ObjectConfig {
	queryFields := graphql.Fields{}
	var wg sync.WaitGroup
	wg.Add(len(dataMap))
	// generate a query for each data in dataMap
	for collectionName, collectionData := range dataMap {
		go func(collectionName string, collectionData DataTypeMapArray) {
			defer wg.Done()
			// create a fields array for each key in collectionData
			fields := graphql.Fields{}
			// for each entry in collectionData create a field
			collectionArgs := graphql.FieldConfigArgument{}

			for _, entry := range collectionData {
				// create a field for each key in entry
				for key, value := range entry {
					field := graphql.Field{}
					field.Name = key
					field.Type = getGraphQlType(value.typeName)
					field.Resolve = func(p graphql.ResolveParams) (interface{}, error) {
						raw := p.Source.(DataTypeMap)[p.Info.FieldName].rawData
						dataType := p.Source.(DataTypeMap)[p.Info.FieldName].typeName
						switch dataType {
						case "String":
							return raw, nil
						case "Int":
							return strconv.ParseInt(raw, 10, 64)
						case "Float":
							return strconv.ParseFloat(raw, 64)
						case "Boolean":
							return strconv.ParseBool(raw)
						default:
							return raw, nil
						}
					}
					fields[key] = &field
					arg := graphql.ArgumentConfig{Type: getGraphQlType(value.typeName)}
					collectionArgs[key] = &arg
				}
			}
			// create a query for each data
			queryFields[collectionName] = &graphql.Field{
				Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
					Name:   collectionName,
					Fields: fields,
				},
				)),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					collection := dataMap[p.Info.FieldName]
					// if there is some arguments, filter the data
					if len(p.Args) > 0 {
						for key, value := range p.Args {
							collection = filterData(collection, key, value)
						}
						return collection, nil
					}
					// if there is no arguments, return the collection
					return collection, nil

				}}
			queryFields[collectionName].Args = collectionArgs
		}(collectionName, collectionData)
	}
	wg.Wait()
	return graphql.ObjectConfig{Name: "Query", Fields: queryFields}
}

func filterData(collection []DataTypeMap, key string, value interface{}) []DataTypeMap {
	var filteredCollection []DataTypeMap
	for _, entry := range collection {
		if convertedValue, ok := value.(int); ok {
			fieldValue, _ := strconv.ParseInt(entry[key].rawData, 10, 64)
			if fieldValue == int64(convertedValue) {
				filteredCollection = append(filteredCollection, entry)
			}
		}

		if convertedValue, ok := value.(float64); ok {
			fieldValue, _ := strconv.ParseFloat(entry[key].rawData, 64)
			if fieldValue == convertedValue {
				filteredCollection = append(filteredCollection, entry)
			}
		}

		if convertedValue, ok := value.(bool); ok {
			fieldValue, _ := strconv.ParseBool(entry[key].rawData)
			if fieldValue == convertedValue {
				filteredCollection = append(filteredCollection, entry)
			}
		}

		if convertedValue, ok := value.(string); ok {
			if entry[key].rawData == convertedValue {
				filteredCollection = append(filteredCollection, entry)
			}
		}

	}
	return filteredCollection
}

func generateGraphqlObjects(dataMap DataMap) graphql.SchemaConfig {
	query := generateQueries(dataMap)
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(query)}
	return schemaConfig
}

func createSchema(schemaConfig graphql.SchemaConfig) graphql.Schema {
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
	return schema
}

func startGraphQLServer(schema graphql.Schema) {
	h := gqlhandler.New(&gqlhandler.Config{
		Schema: &schema,
		Pretty: true,
	})
	http.Handle("/graphql", CorsMiddleware(AuthMiddleware(h)))
}

func getBodyString(r *http.Request) string {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return string(body)
}

func identifyIntrospectionQuery(body string) bool {
	if strings.Contains(body, "query IntrospectionQuery") {
		return true
	}
	return false
}

func extractQueriedDataset(body string) []string {
	writing := true
	stringo := ""
	for _, char := range body {
		if char == '{' || char == '(' || char == '}' || char == ')' {
			writing = !writing
			continue
		}
		if writing {
			if char != '"' && char != '\'' && char != ' ' {
				stringo += string(char)
			}
		}
	}
	result := strings.Split(stringo, ",")
	return result
}
