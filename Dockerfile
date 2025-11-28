# --- Build Stage ---
FROM golang:1.25.4-alpine AS builder

# Set working directory
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download modules
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 is required for a static build on Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# --- Production Stage ---
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install basic dependencies (optional, but useful for debugging)
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/main .

# Copy entrypoint script
COPY entrypoint.sh .

# Give execution permission to entrypoint
RUN chmod +x entrypoint.sh

# Create the media directory for the volume
RUN mkdir -p /app/media

# Expose the port (must match PORT in .env)
EXPOSE 5005

# Set Entrypoint
ENTRYPOINT ["./entrypoint.sh"]

# Default Command
CMD ["./main"]