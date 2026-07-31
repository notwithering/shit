package handler

import "net/http"

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404: Not Found "+r.URL.Path, http.StatusNotFound)
}
