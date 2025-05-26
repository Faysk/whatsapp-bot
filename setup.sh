#!/bin/bash

echo "🛠️ Iniciando verificação e instalação do ambiente WhatsApp Bot..."

# 1. Verifica se Go está instalado
if ! command -v go &> /dev/null; then
    echo "❌ Go não encontrado. Instale o Go para continuar: https://go.dev/doc/install"
    exit 1
fi

# 2. Inicializa o go.mod se não existir
if [ ! -f "go.mod" ]; then
    echo "📦 Inicializando módulo Go..."
    go mod init whatsapp-bot || exit 1
else
    echo "📦 go.mod já existe."
fi

# 3. Instala dependências obrigatórias
echo "⬇️ Instalando pacotes Go necessários..."
go get github.com/joho/godotenv
go get go.mau.fi/whatsmeow
go get modernc.org/sqlite

# 4. Limpeza e ajustes
echo "🧹 Ajustando dependências..."
go mod tidy

# 5. Verifica arquivos importantes
if [ ! -f ".env" ]; then
    echo "⚠️ Arquivo .env não encontrado. Criando exemplo padrão..."
    cp .env.example .env
fi

# 6. Sucesso
echo "✅ Ambiente preparado com sucesso! Você pode rodar: go run ./cmd"
