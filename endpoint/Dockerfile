FROM golang:1.24.3 AS builder
WORKDIR /src

COPY endpoint/go.mod endpoint/go.sum ./endpoint/
COPY rpc ./rpc

WORKDIR /src/endpoint

RUN go mod download

COPY endpoint/. .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/endpoint ./cmd/endpoint

FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /bin/endpoint /app/endpoint

CMD ["/app/endpoint"]