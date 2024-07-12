# Start with a base image containing Go
FROM golang:1.20-alpine

# Install make and other dependencies
RUN apk add --no-cache make git bash postgresql-client

# Install migrate CLI tool using go install
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Ensure the startup script is executable
RUN chmod +x /app/startup.sh

# Build the Go app
RUN go build -o main .

# Copy migration scripts and seed data
# COPY migrations ./migrations
# COPY seeds ./seeds

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
