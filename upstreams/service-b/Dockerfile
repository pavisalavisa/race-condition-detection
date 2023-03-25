# Use an official Go runtime as a parent image
FROM golang:1.16-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go app inside the container
RUN go build -o app main.go

# Set the default command to run when the container starts
CMD ["./app"]

# Expose the port that the app will listen on
EXPOSE 9101
