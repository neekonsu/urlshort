package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if destination, match := pathsToUrls[r.URL.Path]; match {
			http.Redirect(w, r, destination, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	// This is necessary, cannot replace with var pathsMap map[string]string
	pathsMap := makeMap(paths)
	return MapHandler(pathsMap, fallback), nil
}

func parseYAML(yml []byte) ([]path, error) {
	var paths []path
	err := yaml.Unmarshal(yml, &paths)
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func makeMap(paths []path) map[string]string {
	pathsMap := make(map[string]string)
	for _, p := range paths {
		pathsMap[p.Path] = p.URL
	}
	return pathsMap
}

type path struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
