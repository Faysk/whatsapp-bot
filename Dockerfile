# syntax=docker/dockerfile:1.4

########################################
# 1) Build Stage
########################################
ARG GO_VERSION=1.24-alpine
FROM golang:${GO_VERSION} AS builder

# Corrige falhas com IPv6 e proxy.golang.org
ENV GODEBUG=netdns=go+1 \
    GOPROXY=https://proxy.golang.org,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Etapa de cache de dependências
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copia o restante do projeto e compila binário
COPY . ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -o bot ./cmd

########################################
# 2) Runtime Stage
########################################
FROM alpine:3.18

LABEL org.opencontainers.image.source="https://github.com/Faysk/whatsapp-bot" \
      org.opencontainers.image.maintainer="Renan Silva <faysk.nan@gmail.com>" \
      org.opencontainers.image.version="v1.0.0"

# Instala certificados e timezone
RUN apk add --no-cache ca-certificates tzdata

# Cria usuário seguro
RUN addgroup -S app && adduser -S app -G app
USER app

WORKDIR /app

# Copia binário e arquivos necessários
COPY --from=builder --chown=app:app /app/bot .
COPY --from=builder --chown=app:app /app/authorized.json .
COPY --from=builder --chown=app:app /app/crypto_records.json .

# Variáveis padrão (podem ser sobrescritas por .env)
ENV DB_DRIVER=postgres \
    DB_PATH=postgres://bot_user:bot_senha@db:5432/whatsapp_bot?sslmode=disable&binary_parameters=true \
    BOT_NAME=FayskBot \
    LOG_LEVEL=INFO \
    LANG=pt-BR \
    TZ=America/Sao_Paulo

# Healthcheck simples
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s \
  CMD pgrep -f '/app/bot' > /dev/null || exit 1

ENTRYPOINT ["./bot"]
