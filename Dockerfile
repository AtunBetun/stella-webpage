# Start from the official golang base image
FROM golang:1.18

# Set the working directory in the container
WORKDIR /app

# Copy the local package files to the container's working directory
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
