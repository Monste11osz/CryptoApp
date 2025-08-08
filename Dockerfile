FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
ENV GOPROXY=direct
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main ./cmd/main.go

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/main ./

CMD ["./main"]