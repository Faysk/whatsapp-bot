# 🤖 WhatsApp Bot - Modular, Inteligente e Escalável com Go

Este é um bot de WhatsApp altamente modular, resiliente e expansível, construído com [WhatsMeow](https://github.com/tulir/whatsmeow) em Go 1.24+. Ele é preparado para produção, com PostgreSQL 17+, Docker, GPT-4o, geração de voz e notícias cripto automáticas.

---

## 🚀 Funcionalidades

- ✅ Conexão persistente via QR Code com o WhatsApp
- 💬 Comandos automáticos (`!ping`, `!help`, `!gpt`)
- 🔒 Lista de números autorizados com controle dinâmico
- 🧠 Integração com OpenAI GPT-4o (respostas IA)
- 📰 Notícias de Criptomoedas via API CryptoPanic com tradução automática
- 📊 Monitoramento de ATH (all-time-high) com alertas
- 🔁 Tarefas agendadas (CRON) com Gocron
- 📦 Banco de dados PostgreSQL 100% compatível com WhatsMeow
- 🔌 Arquitetura limpa e modular: comandos, eventos, serviços, handlers
- ⚙️ Instalação automática e verificação de dependências com `setup.sh`
- 🐳 Suporte total a Docker com `run_docker.sh`

---

## 📁 Estrutura do Projeto

```
whatsapp-bot/
├── cmd/              # main.go (entrada)
├── config/           # Carregamento de variáveis do .env
├── events/           # Webhooks e eventos WhatsApp
├── handlers/         # Interpretação de mensagens
│   └── commands/     # Comandos textuais (ex: !ping, !gpt)
├── services/         # Lógica: IA, cripto, notificações
├── openai/           # Integração GPT-4o via OpenAI API
├── store/            # Sessão, usuários autorizados, etc.
├── utils/            # Logger e helpers
├── migrations/       # Scripts SQL do PostgreSQL
├── scripts/          # Shell e auxiliar (setup.sh, run_docker.sh)
├── .env.example      # Modelo do ambiente
├── Dockerfile        # Build do bot
├── Dockerfile.db     # Build do banco PostgreSQL
├── docker-compose.yml
└── README.md
```

---

## 📦 Pré-requisitos

- Go 1.24+ (instalação automática se rodar `setup.sh`)
- PostgreSQL 17+ (local ou via Docker)
- Docker e Docker Compose v2+
- Git instalado
- Conta e API Key OpenAI válidas

---

## ⚙️ Instalação com Script Automático

```bash
git clone https://github.com/Faysk/whatsapp-bot.git
cd whatsapp-bot
chmod +x setup.sh run_docker.sh
./setup.sh
```

Esse script irá:
- Validar ou instalar Go 1.24.3 automaticamente
- Baixar e configurar dependências (whatsmeow, godotenv, gocron...)
- Criar `.env` baseado no `.env.example`
- Fazer o `go mod tidy`

---

## ▶️ Executar com Docker (Recomendado)

```bash
./run_docker.sh
```

Isso irá:
- Subir banco PostgreSQL com schema compatível
- Buildar a imagem Go e iniciar o bot
- Verificar containers até estarem saudáveis

Logs em tempo real:

```bash
docker compose logs -f bot
```

---

## 📄 Exemplo de .env

```env
DB_DRIVER=postgres
DB_PATH=postgres://bot_user:bot_senha@db:5432/whatsapp_bot?sslmode=disable

LOG_LEVEL=INFO
PORT=8080

BOT_NAME=FayskBot
LANG=pt-BR

OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4o
ENABLE_CHATGPT=true
MAX_TOKENS=2000
TEMPERATURE=0.7

AUTHORIZED_NUMBERS=5511999999999,5511988888888
RESTRICT_TO_GROUP=false
```

---

## 💬 Comandos Disponíveis

| Comando     | Descrição |
|-------------|-----------|
| `!ping`     | Testa se o bot está online |
| `!help`     | Lista os comandos disponíveis |
| `!gpt`      | Envia pergunta para GPT-4o |
| `!noticias` | Exibe notícias cripto (CryptoPanic traduzido) |

---

## 🔐 Segurança e Permissões

- Lista `AUTHORIZED_NUMBERS` no `.env` para controle inicial
- Adição/remoção dinâmica via comando do próprio bot
- Suporte futuro a escopos: admin, leitura, grupos restritos

---

## 🔮 Futuro Próximo

- [x] Tradução de notícias automática com fallback
- [x] Verificação de ATH (all-time-high) por criptomoeda
- [ ] Dashboard web com estatísticas e controle
- [ ] Upload e resposta com mídia
- [ ] Webhook para automações externas

---

## 🧠 Integrações IA

- 🔮 GPT-4o da OpenAI com suporte a parâmetros avançados
- 🌍 Tradução de conteúdo automatizada
- 📈 Análise de mercado com feedback otimizado para WhatsApp

---

## 🛠️ Desenvolvimento Local

```bash
go run ./cmd
```

Build manual:
```bash
go build -o bot ./cmd
./bot
```

---

## 🛡️ Licença

Este projeto está sob a Licença MIT.

---

## 👨‍💻 Autor

Desenvolvido por [Renan Silva (Faysk)](https://github.com/faysk)  
📧 faysk.nan@gmail.com  
🌐 https://faysk.top