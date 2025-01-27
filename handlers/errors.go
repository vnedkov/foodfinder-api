package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// errorResponse is a struct that represents an error response
type errorResponse struct {
	Error string `json:"error"`
}

// handleError is a helper function that writes an error response to the client
func handleError(w http.ResponseWriter, r *http.Request, err error, responseStatus int, msg ...string) bool {
	errMsg := strings.Join(msg, "")
	if err == nil {
		return false
	}

	log.Err(err).Msg(errMsg)

	resp, _ := json.Marshal(errorResponse{Error: errMsg})
	http.Error(w, string(resp), responseStatus)
	return true
}
