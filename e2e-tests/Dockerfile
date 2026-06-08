FROM golang:1.21-alpine

WORKDIR /app

# Install required tools
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Set environment variables
ENV TEST_NAMESPACE=test-auto

# Set default command
CMD ["sh", "-c", "go test -v -timeout 60s ./..."]

# Alternative: build binary
# RUN go build -o e2e-tests .
