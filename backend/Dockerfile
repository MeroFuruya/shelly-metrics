ARG GO_VERSION=1.23.0
FROM golang:${GO_VERSION}-alpine AS build

COPY go.mod go.sum ./
RUN go mod download -x

COPY backend/ backend/
COPY core/ core/
RUN go build -o /shelly-analytics ./backend

FROM alpine:3.12

COPY --from=build /shelly-analytics /usr/local/bin/shelly-analytics

ENV MINICHAT_PORT=3000
EXPOSE $MINICHAT_PORT

CMD [ "shelly-analytics" ]