# Start from the latest golang base image
FROM golang:latest
ENV GO111MODULE=on

# Add Maintainer Info
LABEL maintainer="Jun Makino <jun.makino@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /config
COPY config/config.yaml config.yaml
# Copy go mod and sum files
#COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
#RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
WORKDIR /
COPY greeter_server greeter_server
COPY greeter_client greeter_client

# Build the Go app
WORKDIR /greeter_server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
WORKDIR /greeter_client
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 

# Expose port 8080 to the outside world
EXPOSE 50051

# Command to run the executable
#CMD ["./greeter_server"]

# Back to root dir
WORKDIR /