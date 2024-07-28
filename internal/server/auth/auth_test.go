package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAuthorizationKey(t *testing.T) {
	tests := map[string]struct {
		req           *http.Request
		token         string
		expectedToken string
		expectedErr   error
	}{
		"valid token": {
			req:           httptest.NewRequest("GET", "https://example.com", nil),
			token:         "randomtesttoken123",
			expectedToken: "randomtesttoken123",
			expectedErr:   nil,
		},
		"no token": {
			req:           httptest.NewRequest("GET", "https://example.com", nil),
			token:         "",
			expectedToken: "",
			expectedErr:   errors.New("invalid authorization key"),
		},
	}

	for name, test := range tests {
		t.Logf("Running test: %s", name)
		if name == "valid token" {
			test.req.Header.Set("Authorization", "Bearer "+test.token)
		}
		token, err := GetAuthorizationKey(test.req)
		if name == "valid token" {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedToken, token)
		} else {
			assert.Error(t, err)
			assert.Equal(t, token, "")
		}
	}
}
