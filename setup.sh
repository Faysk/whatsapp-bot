#!/bin/bash

set -e

echo "🛠️ Iniciando verificação e preparação do ambiente WhatsApp Bot..."

GO_REQUIRED_VERSION="1.24.3"
GO_INSTALL_DIR="/usr/local/go"
GO_DOWNLOAD_URL="https://go.dev/dl/go${GO_REQUIRED_VERSION}.linux-amd64.tar.gz"

# 1. Função para comparar versões
version_lt() {
    [ "$(printf '%s\n' "$1" "$2" | sort -V | head -n1)" != "$2" ]
}

# 2. Verifica se Go está instalado e se a versão é suficiente
if ! command -v go &> /dev/null; then
    echo "⚠️ Go não encontrado. Instalando Go ${GO_REQUIRED_VERSION}..."
    INSTALL_GO=true
else
    GO_CURRENT_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if version_lt "$GO_CURRENT_VERSION" "$GO_REQUIRED_VERSION"; then
        echo "⚠️ Go $GO_CURRENT_VERSION é menor que a exigida ($GO_REQUIRED_VERSION). Atualizando..."
        INSTALL_GO=true
    else
        echo "✅ Go $GO_CURRENT_VERSION está instalado e é compatível."
        INSTALL_GO=false
    fi
fi

# 3. Instala ou atualiza o Go se necessário
if [ "$INSTALL_GO" = true ]; then
    cd /tmp
    wget -q --show-progress "$GO_DOWNLOAD_URL"
    sudo rm -rf "$GO_INSTALL_DIR"
    sudo tar -C /usr/local -xzf "go${GO_REQUIRED_VERSION}.linux-amd64.tar.gz"
    echo "✅ Go ${GO_REQUIRED_VERSION} instalado com sucesso."
    export PATH="/usr/local/go/bin:$PATH"
    echo 'export PATH="/usr/local/go/bin:$PATH"' >> ~/.bashrc
    source ~/.bashrc || true
fi

# 4. Exibe versão Go ativa
go version

# 5. Inicializa módulo se necessário
if [ ! -f "go.mod" ]; then
    echo "📦 Inicializando módulo Go..."
    go mod init whatsapp-bot || exit 1
else
    echo "📦 Módulo Go já existe."
fi

# 6. Instala dependências reais do projeto
echo "⬇️ Instalando dependências do projeto..."
go get github.com/joho/godotenv
go get github.com/gorilla/websocket
go get github.com/go-co-op/gocron
go get github.com/Baozisoftware/qrcode-terminal-go
go get go.mau.fi/whatsmeow@latest

# 7. Limpa pacotes não usados
echo "🧹 Limpando dependências órfãs..."
go mod tidy

# 8. Gera .env se necessário
if [ ! -f ".env" ]; then
    echo "⚠️ .env não encontrado. Copiando de .env.example..."
    cp .env.example .env
fi

echo ""
echo "✅ Ambiente preparado com sucesso!"
echo "ℹ️ Para rodar localmente: go run ./cmd"
echo "ℹ️ Para rodar com Docker: ./run_docker.sh"
