FROM golang:1.22.5

WORKDIR /app

ENV CGO_ENABLED 0

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./model ./model

RUN go build -o main ./cmd/api/main.go

CMD ["./main"]