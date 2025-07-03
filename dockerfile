# Use official Go image
FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Copy files
COPY . .

# Build Go app from cmd/bot directory
RUN go build -o main ./cmd/bot

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./main"]
