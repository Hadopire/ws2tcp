package proxy

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/Hadopire/ws2tcp/errors"
	"github.com/Hadopire/ws2tcp/utils"
	"github.com/gorilla/websocket"
)

type Target struct {
	Host string `json:"host"`
}

const (
	wsWriteWait    = 10 * time.Second
	wsPongWait     = 60 * time.Second
	wsPingPeriod   = (wsPongWait * 9) / 10
	maxMessageSize = 512
)

func ReadTarget(target *net.TCPConn, buffer chan []byte) {
	defer close(buffer)

	for {
		data := make([]byte, maxMessageSize)
		i, err := target.Read(data)
		if err != nil {
			log.Println(err)
			return
		}

		buffer <- data[:i]
	}
}

func TargetToSource(source *websocket.Conn, target *net.TCPConn) {
	ticker := time.NewTicker(wsPingPeriod)
	defer func() {
		ticker.Stop()
		source.Close()
		target.Close()
	}()

	buffer := make(chan []byte, 64)
	go ReadTarget(target, buffer)

	for {
		select {
		case data, ok := <-buffer:
			source.SetWriteDeadline(time.Now().Add(wsWriteWait))
			if !ok {
				source.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := source.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			source.SetWriteDeadline(time.Now().Add(wsWriteWait))
			err := source.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func SourceToTarget(source *websocket.Conn, target *net.TCPConn) {
	defer func() {
		source.Close()
		target.Close()
	}()

	source.SetReadLimit(maxMessageSize)
	source.SetReadDeadline(time.Now().Add(wsPongWait))
	source.SetPongHandler(func(string) error { source.SetReadDeadline(time.Now().Add(wsPongWait)); return nil })
	for {
		_, message, err := source.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		_, err = target.Write(message)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func Init(source *websocket.Conn) (*net.TCPConn, error) {
	source.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, message, err := source.ReadMessage()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var t Target
	err = json.Unmarshal(message, &t)
	if err != nil {
		log.Println(err)
		utils.JSON(source, errors.InvalidJSON)
		return nil, err
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", t.Host)
	if err != nil {
		log.Println(err)
		utils.JSON(source, errors.Hostname)
		return nil, err
	}
	target, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(err)
		utils.JSON(source, errors.ConnectionFailed)
		return nil, err
	}

	return target, nil
}
