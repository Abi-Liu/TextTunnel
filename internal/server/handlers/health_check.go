package handlers

import (
	"net/http"
)

func GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	RespondWithJson(w, 200, "OK")
}
