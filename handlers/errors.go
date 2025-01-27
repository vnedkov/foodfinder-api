package handlers

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// handleError is a helper function that writes an error response to the client
func handleError(w http.ResponseWriter, r *http.Request, err error, responseStatus int, msg ...string) bool {
	var errMsg string
	if err == nil {
		return false
	}

	log.Err(err).Msg(errMsg)
	http.Error(w, errMsg, responseStatus)
	return true
}
