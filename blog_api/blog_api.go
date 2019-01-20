package main

import (
	"blog_api/types"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"net/http"
	"os"
	"time"
)

func main() {

	http.HandleFunc("/blog", func(res http.ResponseWriter, req *http.Request) {
		if _, err := res.Write(serve(req.FormValue("query"))); err != nil {
			stderr(err)
			return
		}
	})

	stderr(http.ListenAndServe(":9959", nil))
}

func serve(query string) []byte {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"post": types.Post,
			"tags": types.Tags,
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"create": types.PostCreate,
			"update": types.PostUpdate,
			"delete": types.PostDelete,
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		stderr(err)
		return nil
	}
	result := graphql.Do(graphql.Params{Schema: schema, RequestString: query})
	if result.Errors != nil {
		stderr(fmt.Errorf("wrong query"))
		return nil
	}

	JSON, err := json.Marshal(result.Data)
	if err != nil {
		stderr(err)
		return nil
	}
	return JSON
}

func stderr(err error) {
	if _, err := fmt.Fprintln(os.Stderr, time.Now().String()[:19], err); err != nil {
		return
	}
}