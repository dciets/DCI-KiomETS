# go build
FROM golang:1.22 as go-build

WORKDIR /backend

## install dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

## run backend server
COPY backend/**/*.go ./
ENV PORT=:8080
RUN CGO_ENABLED=0 GOOS=linux go build -o ./binary


# node (angular) build
FROM node:21.7 as node-build

WORKDIR /frontend

## install dependencies
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install

## build angular app
COPY frontend ./
RUN npm run build

# final stage
FROM scratch

COPY --from=go-build /backend/binary /app/backend-binary
COPY --from=node-build /frontend/dist /app/dist

CMD ["/app/backend-binary"]
