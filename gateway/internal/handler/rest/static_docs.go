package rest

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/apidocs
var apiDocsFS embed.FS

// APIDocsPage serves the static API documentation HTML site.
func APIDocsPage() http.Handler {
	sub, err := fs.Sub(apiDocsFS, "static/apidocs")
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "documentation not available", http.StatusInternalServerError)
		})
	}
	return http.FileServer(http.FS(sub))
}
