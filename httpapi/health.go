package httpapi

import (
	"log"
	"net/http"
)

func Health(w http.ResponseWriter, _ *http.Request) {
	log.Println("health check - OK")
	w.WriteHeader(http.StatusOK)
}
