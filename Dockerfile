# Dockerfile to containerize a go application, the first stage is a build stage and second is the runtime stage that will contains a healthcheck

# Build stage
FROM golang:1.20.2-alpine3.17 as builder

# Maintainer
LABEL maintainer="Mattéo Chrétien <contact@matteochretien.com>"

# Install git
RUN apk update && apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Use cache mount to speed up install of existing dependencies
RUN --mount=type=cache,target=/root/.cache/go-build go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o restapi .

# Runtime stage, the binary will run in a user mode
FROM alpine:3.17

# Maintainer
LABEL maintainer="Mattéo Chrétien <contact@matteochretien.com>"

# Install curl
RUN apk update && apk add --no-cache curl

# Set the Current Working Directory inside the container
WORKDIR /app

# Create user appuser to run the application
RUN adduser -D -g '' appuser

# Switch to non-root user
USER appuser

# Copy the Pre-built binary file from the previous stage
COPY --from=builder --chown=appuser:appuser /app/restapi /app/restapi

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["/app/restapi"]

# Healthcheck
HEALTHCHECK --interval=5s --timeout=3s --start-period=5s --retries=3 CMD curl -f http://localhost:3000/health || exit 1