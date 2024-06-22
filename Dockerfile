# Use the official Go image as the base image
FROM golang:1.22

ENV LOCAL_HOST=0.0.0.0

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o quic-pub-sub .

# Specify the command to run when the container starts
CMD ["./quic-pub-sub"]