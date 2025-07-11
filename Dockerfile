# Etapa de build
FROM golang:1.23 AS builder

WORKDIR /app

# Copia arquivos go
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário
RUN go build -o main ./cmd

# Etapa final
FROM alpine:latest

WORKDIR /root/

# Copia o binário da etapa de build
COPY --from=builder /app/main .

# Porta da API
EXPOSE 3000

# Comando de execução
CMD ["./main"]