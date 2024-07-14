FROM golang:1.22 as go-build

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./game-server

ENTRYPOINT /app/game-server