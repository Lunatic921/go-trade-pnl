# Use an official Golang runtime as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Install git to clone the repository
RUN apt-get update && apt-get install -y git

# Clone the repository from GitHub
RUN git clone https://github.com/Lunatic921/go-trade-pnl

# Set the working directory to the cloned repository
WORKDIR /app/go-trade-pnl

# Build the Go program
RUN go build -o app .

# Set the entry point for the container
ENTRYPOINT ["./app"]

# Default command to run the program with the -trades argument and path
CMD ["-trades", "/app/volume_path"]