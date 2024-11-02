package shelly

import (
	"errors"
	"io"
	"net/url"

	"github.com/MeroFuruya/shelly-analytics/core/logging"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

var SHELLY_URL = url.URL{
	Scheme: "wss",
	Host:   "info-board.shelly.cloud",
}

type ShellyOptions struct {
	ReconnectEnabled bool
	ReconnectTryMax  int
	OnDataReceived   func(r io.Reader)
}

type ShellyClient struct {
	Conn             *websocket.Conn
	OnDataReceived   func(r io.Reader)
	Logger           zerolog.Logger
	Done             chan struct{}
	ReconnectEnabled bool
	ReconnectTry     int
	ReconnectTryMax  int
}

func NewShellyClient(options ShellyOptions) *ShellyClient {
	return &ShellyClient{
		ReconnectEnabled: options.ReconnectEnabled,
		ReconnectTryMax:  options.ReconnectTryMax,
		Logger:           logging.GetLogger("shelly"),
		OnDataReceived:   options.OnDataReceived,
	}
}

func (c *ShellyClient) Open() error {
	c.Logger.Info().
		Str("url", SHELLY_URL.String()).
		Msg("Connecting to Shelly")

	conn, _, err := websocket.DefaultDialer.Dial(SHELLY_URL.String(), nil)
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to connect to Shelly")
		return err
	}
	c.Conn = conn

	c.Done = make(chan struct{})

	return nil
}

func (c *ShellyClient) Listen() error {
	for {
		if err := c.receive(); err != nil {
			c.Logger.Error().Err(err).Msg("Failed to receive data from Shelly")
			if c.ReconnectEnabled {
				if err := c.reconnect(); err != nil {
					return err
				}
			}
		}
	}
}

func (c *ShellyClient) OpenAndListen() error {
	if err := c.Open(); err != nil {
		return err
	}
	return c.Listen()
}

func (c *ShellyClient) receive() error {
	c.Logger.Info().Msg("Listening for messages from Shelly")
	defer close(c.Done)
	for {
		select {
		case <-c.Done:
			return nil
		default:
			_, r, err := c.Conn.NextReader()
			if err != nil {
				c.Logger.Error().Err(err).Msg("Failed to read message from Shelly")
				return err
			}
			if c.OnDataReceived != nil {
				c.OnDataReceived(r)
			}
		}
	}
}

func (c *ShellyClient) reconnect() error {
	c.Logger.Info().Msg("Reconnecting to Shelly")
	c.ReconnectTry++
	if c.ReconnectTry > c.ReconnectTryMax {
		c.Logger.Error().Msg("Failed to reconnect to Shelly: max tries reached")
		return errors.New("max tries reached")
	}
	c.Close()
	c.Open()
	return nil
}

func (c *ShellyClient) Close() {
	close(c.Done)
	if c.Conn != nil {
		c.Conn.Close()
	}
}
