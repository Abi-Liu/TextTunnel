package http

import (
	"encoding/json"
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
