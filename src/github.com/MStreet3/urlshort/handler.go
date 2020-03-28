package urlshort

import (
	"encoding/json"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

type pathUrl struct {
	Path string
	URL  string
}

type pathUrlParser func([]byte) ([]pathUrl, error)
type urlHandler func([]byte, http.Handler) (http.HandlerFunc, error)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func createHandler(parser pathUrlParser) urlHandler {
	return func(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
		pathUrls, err := parser(data)
		if err != nil {
			return nil, err
		}
		pathToURLMap := buildMap(pathUrls)
		return MapHandler(pathToURLMap, fallback), nil
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return createHandler(parseYaml)(yml, fallback)
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return createHandler(parseJson)(jsn, fallback)
}

func parseYaml(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	if err := yaml.Unmarshal(data, &pathUrls); err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func parseJson(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	if err := json.Unmarshal(data, &pathUrls); err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathToURLMap := make(map[string]string)
	for _, pu := range pathUrls {
		pathToURLMap[pu.Path] = pu.URL
	}
	return pathToURLMap
}
