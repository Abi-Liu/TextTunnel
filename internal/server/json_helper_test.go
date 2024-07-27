package server

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespondWithJson(t *testing.T) {
	tests := map[string]struct {
		statusCode         int
		expectedStatusCode int
		payload            interface{}
		expectedBody       string
	}{
		"valid payload": {
			expectedStatusCode: 200,
			statusCode:         200,
			payload:            map[string]string{"message": "Hello World"},
			expectedBody:       `{"message":"Hello World"}`,
		},
		"empty payload": {
			expectedStatusCode: 201,
			statusCode:         201,
			payload:            map[string]string{},
			expectedBody:       `{}`,
		},
		"null payload": {
			expectedStatusCode: 204,
			statusCode:         204,
			payload:            nil,
			expectedBody:       "null",
		},
	}

	for name, test := range tests {
		t.Logf("Running test: %s", name)
		rr := httptest.NewRecorder()
		RespondWithJson(rr, test.statusCode, test.payload)

		assert.Equal(t, test.expectedStatusCode, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, test.expectedBody, rr.Body.String())
	}
}

func TestRespondWithError(t *testing.T) {
	tests := map[string]struct {
		statusCode         int
		expectedStatusCode int
		msg                string
		expectedBody       string
	}{
		"internal error": {
			statusCode:         500,
			expectedStatusCode: 500,
			msg:                "Oops something went wrong!",
			expectedBody:       `{"error":"Oops something went wrong!"}`,
		},
		"user error": {
			statusCode:         400,
			expectedStatusCode: 400,
			msg:                "bad request!!!",
			expectedBody:       `{"error":"bad request!!!"}`,
		},
	}

	for name, test := range tests {
		t.Logf("Running tests: %s", name)
		rr := httptest.NewRecorder()
		RespondWithError(rr, test.statusCode, test.msg)

		assert.Equal(t, test.expectedStatusCode, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, test.expectedBody, rr.Body.String())
	}
}
