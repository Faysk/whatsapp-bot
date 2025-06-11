#!/bin/bash
set -e

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}🐳 Inicializando WhatsApp Bot com Docker...${NC}"

# 1. Verifica Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker não está instalado. Instale em: https://docs.docker.com/get-docker/${NC}"
    exit 1
fi

# 2. Verifica Docker Compose V2
if ! docker compose version &> /dev/null; then
    echo -e "${RED}❌ Docker Compose V2 não encontrado. Atualize para versão que suporta 'docker compose'.${NC}"
    exit 1
fi

# 3. Verifica arquivo docker-compose.yml
if [ ! -f "docker-compose.yml" ]; then
    echo -e "${RED}❌ Arquivo docker-compose.yml não encontrado em $(pwd)${NC}"
    exit 1
fi

# 4. Opção --clean
if [[ "$1" == "--clean" ]]; then
    echo -e "${YELLOW}🧹 Limpando containers e volumes anteriores...${NC}"
    docker compose down --volumes || true
fi

# 5. Build das imagens
echo -e "${YELLOW}🔨 Construindo imagens...${NC}"
docker compose build

# 6. Subindo os containers
echo -e "${YELLOW}🚀 Subindo containers...${NC}"
docker compose up -d

# 7. Função para aguardar container saudável
function wait_for_health() {
    local container=$1
    local timeout=60
    local elapsed=0

    echo -e "${YELLOW}⏱️ Aguardando '${container}' ficar saudável...${NC}"
    while [[ "$(docker inspect --format='{{.State.Health.Status}}' "$container" 2>/dev/null)" != "healthy" ]]; do
        if (( elapsed >= timeout )); then
            echo -e "${RED}❌ Tempo esgotado aguardando '${container}'. Verifique com:${NC}"
            echo -e "${RED}   docker compose logs -f $container${NC}"
            exit 1
        fi
        echo "⏳ $container ainda não está pronto... (${elapsed}s)"
        sleep 2
        elapsed=$((elapsed + 2))
    done
    echo -e "${GREEN}✅ $container está saudável!${NC}"
}

wait_for_health "whatsapp-db"
wait_for_health "whatsapp-bot"

# 8. Finalização
echo ""
echo -e "${GREEN}✅ Ambiente iniciado com sucesso!${NC}"
echo -e "${YELLOW}📜 Use: docker compose logs -f bot${NC}"
