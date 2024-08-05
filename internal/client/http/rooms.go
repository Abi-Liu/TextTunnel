package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Abi-Liu/TextTunnel/internal/models"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

func (c *HttpClient) ConnectToSocket(id uuid.UUID) (*websocket.Conn, context.Context, context.CancelFunc, error) {
	header := http.Header{}
	header.Set("Authorization", "Bearer "+c.AuthToken)
	dialOpts := websocket.DialOptions{HTTPHeader: header}

	ctx, cancel := context.WithCancel(context.Background())
	conn, _, err := websocket.Dial(ctx, BASE_URL+"/ws/"+id.String(), &dialOpts)
	if err != nil {
		cancel()
		return nil, nil, nil, err
	}
	return conn, ctx, cancel, nil
}

func (c *HttpClient) FetchRooms() ([]models.Room, error) {
	res, err := c.Get(BASE_URL + "/rooms")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	rooms := []models.Room{}
	err = json.Unmarshal(data, &rooms)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (c *HttpClient) CreateRoom(name string) (models.Room, error) {
	type Req struct {
		Name string `json:"name"`
	}
	params := Req{Name: name}
	data, err := json.Marshal(params)
	if err != nil {
		return models.Room{}, fmt.Errorf("failed to Marshal to json: %s", err)
	}
	reader := bytes.NewReader(data)
	res, err := c.Post(BASE_URL+"/rooms", reader)
	if err != nil {
		return models.Room{}, fmt.Errorf("failed to create room: %s", err)
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return models.Room{}, fmt.Errorf("failed to read response body: %s", err)
	}
	room := models.Room{}
	err = json.Unmarshal(dat, &room)
	if err != nil {
		return models.Room{}, fmt.Errorf("failed to Unmarshal data: %s", err)
	}

	return room, nil
}
