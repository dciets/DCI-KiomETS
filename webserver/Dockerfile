# node (angular) build
FROM node:21.7 as node-build

WORKDIR /frontend

## install dependencies
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install

## build angular app
COPY frontend ./
RUN npm run build

# go build
FROM golang:1.22 as go-build

WORKDIR /backend

## install dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

## run backend webserver
COPY backend ./
RUN go build -o /app/webserver-binary ./src

COPY --from=node-build /frontend/dist /app/dist

ENTRYPOINT /app/webserver-binary
