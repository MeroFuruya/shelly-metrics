package main

import (
	"context"
	"io"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/MeroFuruya/shelly-analytics/core/database"
	"github.com/MeroFuruya/shelly-analytics/core/logging"
	"github.com/MeroFuruya/shelly-analytics/core/shelly"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load .env file")
	}

	args := GetArgs()

	logging.SetupLogger(args.Format.Format)

	logger := logging.GetLogger("main")

	logger.Debug().Msg("Starting Shelly Analytics")

	database.InitConnection(context.TODO(), database.DatabaseConfig{
		Host:     args.PostgresHost,
		Port:     args.PostgresPort,
		Database: args.PostgresDatabase,
		User:     args.PostgresUser,
		Password: args.PostgresPassword,
	})

	conn := database.GetConnection()

	batch := pgx.Batch{}
	shelly.CreateMetricTables(&batch)
	result := conn.SendBatch(context.TODO(), &batch)
	result.Close()
	batch = pgx.Batch{}
	shelly.CreateMetricHyperTablesBatch(&batch)
	result = conn.SendBatch(context.TODO(), &batch)
	result.Close()

	batch = pgx.Batch{}
	shelly.InsertMetric(shelly.NewDeviceMetric{
		Type:      "SHSW-1",
		Latitude:  0.0,
		Longitude: 0.0,
		Display:   "Shelly 1",
		Friendly:  "NW, Middle of the earth",
		Family:    "Shelly",
		Timestamp: time.Now(),
	}, &batch)
	result = conn.SendBatch(context.Background(), &batch)
	result.Close()

	c := shelly.NewShellyClient(shelly.ShellyOptions{
		ReconnectEnabled: true,
		ReconnectTryMax:  5,
		OnDataReceived: func(r io.Reader) {
			if batch, err := shelly.Parse(r); err != nil {
				logger.Error().Err(err).Msg("Failed to parse data")
			} else {
				result := conn.SendBatch(context.Background(), &batch.Batch)
				result.Close()
			}
		},
	})
	c.Open()
	c.Listen()
	c.Close()
}
