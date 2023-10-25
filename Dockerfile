# Use an official Golang runtime as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /tradewiz

# Install git to clone the repository
RUN apt-get update && apt-get install -y git

# Clone the repository from GitHub
RUN git clone https://github.com/Lunatic921/go-trade-pnl

# Set the working directory to the cloned repository
WORKDIR /tradewiz/go-trade-pnl

# Build the Go program
RUN go build -o ../tw .

RUN mkdir /tradewiz/trades

VOLUME /tradewiz/trades

# Set the entry point for the container
ENTRYPOINT ["/tradewiz/tw", "-trades", "/tradewiz/trades"]