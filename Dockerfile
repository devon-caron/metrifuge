# Use the official Golang image
FROM golang:latest

# Root user is default, so this line is not necessary
# USER root:root

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
# Corrected: specify input package and output file separately
RUN go build -o ./main .

# Expose the port the app runs on
EXPOSE 8000

# Command to run the executable
CMD ["./main"]
