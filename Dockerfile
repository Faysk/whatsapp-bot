# syntax=docker/dockerfile:1.4

########################################
# 1) Build Stage
########################################
ARG GO_VERSION=1.24-alpine
FROM golang:${GO_VERSION} AS builder

# Usa proxy oficial Go com fallback
ENV GOPROXY=https://proxy.golang.org,direct

# Compilação estática para Linux amd64
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Etapa de cache de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do projeto e compila binário principal
COPY . ./
RUN go build -ldflags="-s -w" -o bot ./cmd

########################################
# 2) Runtime Stage
########################################
FROM alpine:3.18

# Metadados da imagem
LABEL org.opencontainers.image.source="https://github.com/Faysk/whatsapp-bot" \
      org.opencontainers.image.maintainer="Renan Silva <faysk.nan@gmail.com>" \
      org.opencontainers.image.version="v1.0.0"

# Instala certificados e timezone
RUN apk add --no-cache ca-certificates tzdata

# Cria usuário não-root seguro
RUN addgroup -S app && adduser -S app -G app
USER app

WORKDIR /app

# Copia binário e arquivos necessários do build
COPY --from=builder --chown=app:app /app/bot ./bot
COPY --from=builder --chown=app:app /app/authorized.json ./authorized.json
COPY --from=builder --chown=app:app /app/crypto_records.json ./crypto_records.json

# Variáveis de ambiente padrão (podem ser sobrescritas externamente)
ENV DB_DRIVER=postgres \
    DB_PATH=postgres://bot_user:bot_senha@db:5432/whatsapp_bot?sslmode=disable&binary_parameters=true \
    BOT_NAME=FayskBot \
    LOG_LEVEL=INFO \
    LANG=pt-BR \
    TZ=America/Sao_Paulo

# Healthcheck simples com verificação do processo
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s \
  CMD pgrep -f '/app/bot' > /dev/null || exit 1

# Comando principal
ENTRYPOINT ["./bot"]
