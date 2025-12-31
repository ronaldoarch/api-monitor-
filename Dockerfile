# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-monitor .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binário do builder
COPY --from=builder /app/api-monitor .
COPY --from=builder /app/web ./web

# Expor porta (Railway vai definir via variável PORT)
EXPOSE 8080

# Comando para executar
CMD ["./api-monitor"]

