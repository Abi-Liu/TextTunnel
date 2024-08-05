package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/Abi-Liu/TextTunnel/internal/models"
)

func (c *HttpClient) GetUserByAuthToken() (models.User, error) {
	res, err := c.Get(BASE_URL + "/users")
	if err != nil {
		return models.User{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{}
	err = json.Unmarshal(data, &user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (c *HttpClient) PostSignUp(username, password string) (models.User, error) {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	req := Req{
		Username: username,
		Password: password,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return models.User{}, err
	}

	reqBody := bytes.NewReader(data)
	res, err := c.Post(BASE_URL+"/users", reqBody)
	if err != nil {
		return models.User{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	var error models.Error

	if res.StatusCode >= 400 {
		err = json.Unmarshal(body, &error)
		if err != nil {
			return models.User{}, err
		}
		return models.User{}, fmt.Errorf(error.Error)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (c *HttpClient) Login(username, password string) (models.User, error) {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	req := Req{
		Username: username,
		Password: password,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return models.User{}, err
	}

	reqBody := bytes.NewReader(data)
	res, err := c.Post(BASE_URL+"/login", reqBody)
	if err != nil {
		return models.User{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	var error models.Error
	if res.StatusCode >= 400 {
		err = json.Unmarshal(body, &error)
		if err != nil {
			return models.User{}, err
		}
		return models.User{}, fmt.Errorf(error.Error)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
