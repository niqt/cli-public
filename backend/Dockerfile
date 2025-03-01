# Stage di build
FROM golang:1.24-alpine AS builder
ENV CGO_ENABLED=1
WORKDIR /app

# Install dependencies
RUN apk add --no-cache gcc musl-dev curl

# Copy go.mod e go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the code
COPY . .

# Compile
RUN mkdir -p bin && go build -o bin/backend ./cmd/main.go

# Final stage
FROM alpine:latest
WORKDIR /app

# Install dependencies (curl needs for the health)
RUN apk add --no-cache libc6-compat curl

# Copy the built
COPY --from=builder /app/bin/backend ./backend
RUN chmod +x ./backend

EXPOSE 8080

CMD ["./backend"]