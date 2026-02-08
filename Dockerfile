# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /build

COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.22-alpine AS backend-builder

WORKDIR /build

COPY backend/go.mod ./
COPY backend/ ./

# Download dependencies and build
RUN go mod download && \
    go mod tidy

COPY --from=frontend-builder /build/dist ./frontend/dist

RUN CGO_ENABLED=0 GOOS=linux go build -o plocate-ui .

# Stage 3: Final Runtime Image
FROM ubuntu:22.04

# Install plocate and other dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    plocate \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Create app directory and necessary directories
WORKDIR /app

# Copy binary from builder
COPY --from=backend-builder /build/plocate-ui /app/plocate-ui

# Create directories for database and config
RUN mkdir -p /app/data /app/config

# Expose port
EXPOSE 8080

# Run as root to allow updatedb to run (plocate requires root for indexing)
# In production, you may want to configure this differently
USER root

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app/plocate-ui", "--health"] || exit 1

CMD ["/app/plocate-ui"]
