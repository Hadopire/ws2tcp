package utils

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

func JSON(c *websocket.Conn, i interface{}) error {
	w, err := c.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return err
	}

	err = w.Close()
	return err
}
