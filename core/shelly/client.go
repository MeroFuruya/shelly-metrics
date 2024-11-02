package shelly

import (
	"net/url"

	"github.com/MeroFuruya/shelly-analytics/core/logging"
	"github.com/gorilla/websocket"
)

var SHELLY_URL = url.URL{
	Scheme: "wss",
	Host:   "info-board.shelly.cloud",
}

var logger = logging.GetLogger("shelly.client")

func channelIsClosed(ch <-chan []byte) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}

func newConn() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(SHELLY_URL.String(), nil)
	return conn, err
}

func Run(data chan []byte) {
	for {
		if channelIsClosed(data) {
			logger.Debug().Msg("Data channel is closed, exiting")
			return
		}

		conn, err := newConn()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to connect to Shelly")
			return
		}

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				logger.Error().Err(err).Msg("Failed to read message")
				break
			}

			data <- message
		}
	}
}
