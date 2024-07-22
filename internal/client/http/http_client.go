package http

import (
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const BASE_URL = "http://localhost:8080"

type HttpClient struct {
	Client    *http.Client
	AuthToken string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ApiKey    string    `json:"api_key"`
}

type Error struct {
	Error string `json:"error"`
}

func CreateHttpClient() *HttpClient {
	client := &HttpClient{
		Client: &http.Client{
			Timeout: 20 * time.Second,
		},
		AuthToken: "",
	}
	return client
}

func (c *HttpClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *HttpClient) Post(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.Do(req)
}

func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+c.AuthToken)
	return c.Client.Do(req)
}

func (c *HttpClient) SetAuthToken(token string) {
	c.AuthToken = token
}
