# Stage 1: Build the Go binary
FROM golang:latest AS builder

WORKDIR /app

# Copy the Go project files
COPY . .

# Build the Go binary for the desired architecture (amd64 in this case)
RUN CGO_ENABLED=0 go build -o myapp -tags prod

# Stage 2: Create a minimal production image
FROM arm64v8/ubuntu:22.04

RUN apt update && apt install -y bash ca-certificates && rm -rf /var/cache/apt/*
WORKDIR /app

# Copy only the binary from the previous stage
COPY --from=builder /app/myapp .

# Expose the port that your application listens on
EXPOSE 13337

# Command to run the executable
ENTRYPOINT ["./myapp"]