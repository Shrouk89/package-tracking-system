# Use a lightweight base image for Go
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Install any necessary dependencies (if using any external packages)
RUN go mod tidy

# Build the Go application
RUN go build -o main .

# Expose the port the app will run on
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]
