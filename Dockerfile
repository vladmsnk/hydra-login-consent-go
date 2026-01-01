FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /login_consent ./app

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary and UI templates
COPY --from=builder /login_consent .
COPY --from=builder /app/ui ./ui

EXPOSE 3000

CMD ["./login_consent"]

