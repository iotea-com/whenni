FROM golang:1.25-bookworm

# Set working directory
WORKDIR /whenni

COPY . .

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Start with Air
CMD ["air", "-c", "/whenni/services/http-api/.air.toml"]
