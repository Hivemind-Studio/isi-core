FROM golang:1.22-alpine AS builder
WORKDIR /src

ENV CGO_ENABLED=0 GOOS=linux

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o /bin/server ./cmd/api

FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

COPY --from=builder /bin/server ./server
COPY template ./template
COPY db/migrations ./db/migrations

ENV APP_PORT=:3000

EXPOSE 3000

# Start the API
CMD ["./server"]