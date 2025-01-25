package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMultiUserRequirements(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)

	// FIXME Wrap these handle functions in a domain-level package to prvent direct use of the handler package
	HealthHandler(responseRecorder, request)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, `{"status": "OK"}`, responseRecorder.Body.String())
}
