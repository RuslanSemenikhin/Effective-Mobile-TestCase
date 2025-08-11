FROM golang:1.24.3 AS builder
WORKDIR /src

COPY controller/go.mod controller/go.sum ./controller/
COPY rpc ./rpc

WORKDIR /src/controller

RUN go mod download

COPY controller/. .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/controller ./cmd/controller

FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /bin/controller /app/controller

CMD ["/app/controller"]