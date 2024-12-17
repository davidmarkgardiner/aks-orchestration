# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o script-orchestrator

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/script-orchestrator .
COPY templates/ templates/
COPY scripts/ scripts/

# Make scripts executable
RUN chmod +x scripts/*.sh

EXPOSE 8080
CMD ["./script-orchestrator"] 