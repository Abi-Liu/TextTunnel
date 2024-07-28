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

func TestHashPassword(t *testing.T) {
	tests := map[string]struct {
		password    string
		expectedErr bool
	}{
		"valid password": {
			password:    "testpw123",
			expectedErr: false,
		},
		"empty password": {
			password:    "",
			expectedErr: false,
		},
	}

	for name, test := range tests {
		t.Logf("Running test: %s", name)
		pw, err := HashPassword(test.password)
		if test.expectedErr {
			assert.Error(t, err)
			assert.Equal(t, "", pw)
		} else {
			assert.Nil(t, err)
			assert.NotEqual(t, "", pw)
		}
	}
}

func TestCompareHashAndPassword(t *testing.T) {
	password := "Password123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password")
	}

	tests := map[string]struct {
		pw            string
		expectedMatch bool
	}{
		"passwords match": {
			pw:            password,
			expectedMatch: true,
		},
		"passwords do not match": {
			pw:            "random",
			expectedMatch: false,
		},
		"empty password": {
			pw:            "",
			expectedMatch: false,
		},
	}

	for name, test := range tests {
		t.Logf("Running test: %s", name)
		match := CompareHashAndPassword(hashedPassword, test.pw)
		assert.Equal(t, test.expectedMatch, match)
	}
}
