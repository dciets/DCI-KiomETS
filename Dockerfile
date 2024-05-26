# go build
FROM golang:1.22 as go-build

WORKDIR /webserver/backend

## install dependencies
COPY webserver/backend/go.mod webserver/backend/go.sum ./
RUN go mod download

## run backend webserver
COPY webserver/backend/**/*.go ./
ENV PORT=:8080
RUN CGO_ENABLED=0 GOOS=linux go build -o ./binary


# node (angular) build
FROM node:21.7 as node-build

WORKDIR /webserver/frontend

## install dependencies
COPY webserver/frontend/package.json webserver/frontend/package-lock.json ./
RUN npm install

## build angular app
COPY webserver/frontend ./
RUN npm run build

# final stage
FROM scratch

COPY --from=go-build /webserver/backend/binary /app/webserver-binary
COPY --from=node-build /webserver/frontend/dist /app/dist

CMD ["/app/webserver-binary"]
