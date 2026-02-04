# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o api ./cmd/api
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o worker ./cmd/worker

# Runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/api /app/api
COPY --from=builder /app/worker /app/worker

EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/api"]
