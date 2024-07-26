package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func (c HttpClient) FetchRooms() ([]Room, error) {
	res, err := c.Get(BASE_URL + "/rooms")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	rooms := []Room{}
	err = json.Unmarshal(data, &rooms)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (c HttpClient) CreateRoom(name string) (Room, error) {
	type Req struct {
		Name string `json:"name"`
	}
	params := Req{Name: name}
	data, err := json.Marshal(params)
	if err != nil {
		return Room{}, fmt.Errorf("Failed to Marshal to json: %s", err)
	}
	reader := bytes.NewReader(data)
	res, err := c.Post(BASE_URL+"/rooms", reader)
	if err != nil {
		return Room{}, fmt.Errorf("Failed to create room: %s", err)
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return Room{}, fmt.Errorf("Failed to read response body: %s", err)
	}
	room := Room{}
	err = json.Unmarshal(dat, &room)
	if err != nil {
		return Room{}, fmt.Errorf("Failed to Unmarshal data: %s", err)
	}

	return room, nil
}
