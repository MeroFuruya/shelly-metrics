package main

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/MeroFuruya/shelly-metrics/core/database"
	"github.com/MeroFuruya/shelly-metrics/core/logging"
	"github.com/MeroFuruya/shelly-metrics/core/shelly"
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

	logger.Debug().Msg("Starting shelly metrics")

	err = database.InitConnection(context.TODO(), database.DatabaseConfig{
		Host:     args.PostgresHost,
		Port:     args.PostgresPort,
		Database: args.PostgresDatabase,
		User:     args.PostgresUser,
		Password: args.PostgresPassword,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to initialize database connection")
		return
	}

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
	result = conn.SendBatch(context.Background(), &batch)
	result.Close()

	data := make(chan []byte, 100)
	go shelly.Run(data)
	for v := range data {
		if batch, err := shelly.Parse(v); err != nil {
			logger.Error().Err(err).Msg("Failed to parse data")
		} else {
			result := conn.SendBatch(context.Background(), &batch.Batch)
			result.Close()
		}
	}
}
