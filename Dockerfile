# Stage 1: Build the Go binary
FROM golang:1.22.3-alpine as builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /wrapter ./main.go

# Stage 2: Create the final image
FROM alpine:latest

# Copy the compiled binary from the builder stage
COPY --from=builder /wrapter /usr/local/bin/wrapter

# Set the entrypoint to the compiled binary
ENTRYPOINT ["wrapter"]
