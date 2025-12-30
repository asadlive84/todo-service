FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
# COPY config ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todo-app cmd/main.go cmd/cmd.go

RUN ls -la migrations/ || echo "Migrations folder found"
RUN ls -la .env || echo ".env file found"

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates curl

COPY --from=builder /app/todo-app .
COPY --from=builder /app/config ./config 

COPY migrations ./migrations

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

CMD ["./todo-app"]