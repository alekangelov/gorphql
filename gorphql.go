package gorphql

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

type Options struct {
	Host string `default:"localhost" description:"The host to listen on"`
	Port int    `default:"8080" description:"The port to listen on"`
	Path string `default:"/graphql" description:"The path to listen on"`
}

var types = ""
var inputs = ""
var queries = ""
var mutations = ""
var subscriptions = ""
var results = map[string]interface{}{
	"query":        map[string]interface{}{},
	"mutation":     map[string]interface{}{},
	"subscription": map[string]interface{}{},
}

type ExtensionType string

type Context interface{}

const (
	QUERY        = "query"
	MUTATION     = "mutation"
	SUBSCRIPTION = "subscription"
)

type ExtendType struct {
	Extension ExtensionType
	Field     string
	Args      map[string]string
	Result    func(args map[string]string, context Context) interface{}
}

func ExtendSchema(e *ExtendType) {
	fmt.Printf("Extending schema with %s\n", e.Extension)
	inpus := ""
	for k, v := range e.Args {
		inpus += fmt.Sprintf("%s: %s", k, v)
	}
	fmt.Printf("Arguments: %s\n", inpus)
	fmt.Printf("Field: %s\n", e.Field)

	results[string(e.Extension)].(map[string]interface{})[e.Field] = e.Result

	fmt.Printf("Results: %v\n", results)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", r.URL.Path)
}

func Init(o *Options) {
	fmt.Println("Initializing GorphQL...")
	http.HandleFunc(o.Path, handler)
	fmt.Printf("Listening on http://%s:%d%s\n", o.Host, o.Port, o.Path)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", o.Host, o.Port), nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
