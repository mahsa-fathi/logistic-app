FROM golang:1.22.0-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY start.sh ./
RUN chmod +x ./start.sh

# Set the entrypoint for the container
ENTRYPOINT ["./start.sh"]
