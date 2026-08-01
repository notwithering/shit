package handler

import "net/http"

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404: Not Found: "+r.URL.Path, http.StatusNotFound)
}

func internalServerError(w http.ResponseWriter, err error) {
	http.Error(w, "500: Internal Server Error: "+err.Error(), http.StatusInternalServerError)
}
