FROM golang:1.19

WORKDIR /usr/src/app

ENV GOOS=linux
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags '-w -s' -a -installsuffix cgo -o webserver cmd/main.go

CMD ["/usr/src/app/webserver"]