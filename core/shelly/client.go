package shelly

import (
	"net/url"

	"github.com/MeroFuruya/shelly-metrics/core/logging"
	"github.com/gorilla/websocket"
)

var SHELLY_URL = url.URL{
	Scheme: "wss",
	Host:   "info-board.shelly.cloud",
}

func Run(data chan<- []byte) {
	var logger = logging.GetLogger("shelly.client")
	for {
		conn, _, err := websocket.DefaultDialer.Dial(SHELLY_URL.String(), nil)
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
