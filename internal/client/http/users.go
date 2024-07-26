package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func (c *HttpClient) GetUserByAuthToken() (User, error) {
	res, err := c.Get(BASE_URL + "/users")
	if err != nil {
		return User{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return User{}, err
	}

	user := User{}
	err = json.Unmarshal(data, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (c *HttpClient) PostSignUp(username, password string) (User, error) {
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
		return User{}, err
	}

	reqBody := bytes.NewReader(data)
	res, err := c.Post(BASE_URL+"/users", reqBody)
	if err != nil {
		return User{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return User{}, err
	}

	var user User
	var error Error

	if res.StatusCode >= 400 {
		err = json.Unmarshal(body, &error)
		if err != nil {
			return User{}, err
		}
		return User{}, fmt.Errorf(error.Error)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (c *HttpClient) Login(username, password string) (User, error) {
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
		return User{}, err
	}

	reqBody := bytes.NewReader(data)
	res, err := c.Post(BASE_URL+"/login", reqBody)
	if err != nil {
		return User{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return User{}, err
	}

	var user User
	var error Error
	if res.StatusCode >= 400 {
		err = json.Unmarshal(body, &error)
		if err != nil {
			return User{}, err
		}
		return User{}, fmt.Errorf(error.Error)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
