package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/MStreet3/urlshort"
)

func main() {
	yamlFileName := flag.String("yaml", "urls.yaml", "a yaml file describing shortened urls")
	jsonFileName := flag.String("json", "urls.json", "a json file describing shortened urls")
	flag.Parse()

	yamlFile, err := ioutil.ReadFile(*yamlFileName)

	if err != nil {
		exit(fmt.Sprintf("Failed to open the YAML file: %s\n", *yamlFileName))
	}

	jsonFile, err := ioutil.ReadFile(*jsonFileName)

	if err != nil {
		exit(fmt.Sprintf("Failed to open the YAML file: %s\n", *jsonFileName))
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler(yamlFile, mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler
	jsonHandler, err := urlshort.JSONHandler(jsonFile, yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
