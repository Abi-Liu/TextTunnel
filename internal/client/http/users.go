package http

import (
	"encoding/json"
	"io"
)

func (h *HttpClient) GetUserByAuthToken() (User, error) {
	res, err := h.Get(BASE_URL + "/users")
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
