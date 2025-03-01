# Stage di build
FROM golang:1.24-alpine AS builder
ENV CGO_ENABLED=1
WORKDIR /app

# Installa dipendenze necessarie
RUN apk add --no-cache gcc musl-dev curl

# Copia go.mod e go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copia il codice sorgente
COPY . .

# Compila l'applicazione direttamente
RUN mkdir -p bin && go build -o bin/backend ./cmd/main.go
# Oppure, se il main è in una sottocartella diversa:
# RUN mkdir -p bin && go build -o bin/backend ./cmd/api/main.go

# Stage finale dell'immagine
FROM alpine:latest
WORKDIR /app

# Installa dipendenze di runtime se necessarie
RUN apk add --no-cache libc6-compat curl

# Copia l'eseguibile compilato
COPY --from=builder /app/bin/backend ./backend
RUN chmod +x ./backend

# Espone la porta utilizzata dall'applicazione
EXPOSE 8080

# Comando predefinito per avviare l'applicazione
CMD ["./backend"]