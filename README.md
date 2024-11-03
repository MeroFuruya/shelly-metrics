# shelly-metrics

I found a public api endpoint that returns some interesting data: `wss://info-board.shelly.cloud/`.

This app scrapes that endpoint and saves it to a postgres database.

Later we might be able to build a frontend to give a nice view of the data.

## Setup

- Make sure you have Timescaledb installed: [docs.timescale.com](https://docs.timescale.com/self-hosted/latest/install/)
- Create a `.env` file in the root of the project according to the `.env.example` file
- run `go run backend/*.go --format pretty` to start the backend
