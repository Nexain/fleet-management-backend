FROM golang:1.22.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o fleet-management-backend ./cmd/main.go

# Use a lightweight image for the final stage
FROM alpine:latest

WORKDIR /root/

# Install bash and curl for wait-for-it
RUN apk add --no-cache bash curl

# Copy the built binary from the builder stage
COPY --from=builder /app/fleet-management-backend .

# Add wait-for-it script
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /usr/local/bin/wait-for-it
RUN chmod +x /usr/local/bin/wait-for-it

# Use wait-for-it to wait for MQTT broker before starting the app
CMD ["sh", "-c", "wait-for-it mqtt:1883 -- wait-for-it rabbitmq:5672 -- ./fleet-management-backend"]