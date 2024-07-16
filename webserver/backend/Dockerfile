FROM golang:1.22

WORKDIR /app/

## install dependencies
COPY go.mod go.sum  ./
RUN go mod download

## run backend webserver
COPY ./ /app/
RUN go build -o /app/backend-binary ./src

ENTRYPOINT /app/backend-binary
