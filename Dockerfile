FROM golang:1.16-buster as builder

WORKDIR /app


COPY go.* ./
RUN go mod download


COPY . ./


RUN go build cmd/service/main.go


FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/main /app/server
CMD ["/app/server"]