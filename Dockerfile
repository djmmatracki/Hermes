FROM golang:1.19.4-alpine3.17

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
ENTRYPOINT go run ./cmd/main.go
EXPOSE 8000