# Use the official golang image as a base
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Install psql
RUN apt-get update && apt-get install -y postgresql-client

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Copy the wait-for script into the container
COPY wait-for-postgres.sh /wait-for-postgres.sh
RUN chmod +x /wait-for-postgres.sh

# Expose port 4040 to the outside world
EXPOSE 4040

# Command to run the executable
CMD /wait-for-postgres.sh db ./main