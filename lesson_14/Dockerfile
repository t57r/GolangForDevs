FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -trimpath -ldflags="-s -w" -o /bin/fiber-server ./cmd/server

FROM alpine:3.20

RUN adduser -D -H -u 10001 appuser
USER appuser

WORKDIR /app
COPY --from=builder /bin/fiber-server /app/server

EXPOSE 8080
ENTRYPOINT ["/app/server"]
