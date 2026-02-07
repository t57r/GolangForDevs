# --- build stage ---
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build server binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /bin/documentstore-server ./cmd/server

# --- runtime stage ---
FROM alpine:3.20

# Non-root user
RUN adduser -D -H -u 10001 appuser
USER appuser

WORKDIR /app

COPY --from=builder /bin/documentstore-server /app/server

EXPOSE 4000

# Server listens on TCP :4000
ENTRYPOINT ["/app/server"]
