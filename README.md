# 🤖 WhatsApp Bot - Modular, Inteligente e Escalável com Go

Este é um bot de WhatsApp altamente modular e escalável, construído com [WhatsMeow](https://github.com/tulir/whatsmeow) e linguagem Go. Projetado para ser limpo, poderoso e pronto para crescer com integrações de IA, voz, banco de dados e serviços externos.

---

## 🚀 Funcionalidades

- ✅ Conexão persistente via QR Code com o WhatsApp
- 💬 Respostas automáticas por comandos (ex: `!ping`, `!help`)
- 🧱 Arquitetura limpa e modularizada (eventos, comandos, serviços)
- 🧠 Integração com OpenAI GPT (mensagens com IA)
- 🗣️ Suporte a áudio com edge-tts (texto para voz)
- 🎙️ Transcrição de áudios recebidos (Whisper)
- 🗂️ Suporte a banco de dados PostgreSQL (recomendado) ou SQLite
- 📦 Configuração via `.env` com exemplo incluído
- 🛠️ Script de instalação automática (`setup.sh`)
- 🔌 Pronto para integrar com APIs, CRON, Webhooks e mais

---

## 📁 Estrutura do Projeto

```
whatsapp-bot/
├── cmd/              # Ponto de entrada principal (main.go)
├── config/           # Carregamento de .env e variáveis globais
├── events/           # Dispatcher de eventos do WhatsApp
├── handlers/         # Comandos e mensagens recebidas
│   └── commands/     # Comandos como !ping, !help, etc.
├── services/         # Lógica de negócio (envio, IA, TTS)
├── openai/           # Integração com ChatGPT
├── scheduler/        # Tarefas agendadas (cron)
├── store/            # Persistência (auth, DB, sessions)
├── utils/            # Logger e utilitários
├── media/            # Áudios e arquivos temporários (.gitignored)
├── authorized.json   # Lista persistente de números autorizados
├── .env.example      # Exemplo de variáveis de ambiente
├── setup.sh          # Instalador automático do projeto
├── go.mod / go.sum   # Dependências do Go
└── README.md
```

---

## 📦 Pré-requisitos

- Go 1.20 ou superior
- Git instalado
- PostgreSQL (recomendado) ou SQLite3
- Python 3.10+ (para edge-tts)
- FFmpeg (para conversão de formatos de áudio)
- API Key da OpenAI (para GPT e Whisper)

---

## ⚙️ Instalação Rápida

```bash
git clone https://github.com/seunome/whatsapp-bot.git
cd whatsapp-bot
chmod +x setup.sh
./setup.sh
```

Esse script irá:
- Validar e instalar dependências do Go
- Instalar bibliotecas (whatsmeow, godotenv, pq/sqlite)
- Criar `.env` com base no `.env.example`
- Preparar o projeto para execução

---

## 📄 Exemplo `.env`

```env
# Tipo de banco: postgres ou sqlite
DB_DRIVER=postgres

# PostgreSQL: postgres://usuario:senha@localhost:5432/whatsapp_bot?sslmode=disable
DB_PATH=postgres://bot_user:bot_senha@localhost:5432/whatsapp_bot?sslmode=disable

LOG_LEVEL=INFO
PORT=8080

BOT_NAME=FayskBot
LANG=pt-BR

OPENAI_API_KEY=sk-...
ENABLE_CHATGPT=true
OPENAI_MODEL=gpt-4o
MAX_TOKENS=2000
TEMPERATURE=0.7

AUTHORIZED_NUMBERS=5511999999999,5511988888888
RESTRICT_TO_GROUP=false
```

---

## ▶️ Como Rodar

### Desenvolvimento:

```bash
go run ./cmd
```

### Produção:

```bash
go build -o bot
./bot
```

---

## 💬 Comandos Suportados

- `!ping` → Testa a resposta do bot
- `!help` → Lista os comandos disponíveis
- `bom dia` → Envia uma saudação especial
- (Em breve) `!falar <texto>` → Converte texto em áudio com voz IA
- (Em breve) `!gpt <pergunta>` → Responde com ChatGPT

---

## 🧠 Integrações de IA

### ✅ GPT via OpenAI
- Geração de respostas contextuais
- Controle por número autorizado

### ✅ Whisper (áudio → texto)
- Transcrição automática de áudios recebidos
- Multi-idioma com detecção automática

### ✅ edge-tts (texto → voz)
- Geração de resposta falada com vozes neurais (Microsoft)
- Configuração de idioma e voz no futuro

---

## 🔮 Possibilidades Futuras

- [x] Geração de áudio por texto com IA
- [x] Transcrição de áudio automática
- [ ] Envio de imagens/stickers com IA
- [ ] Comandos personalizados em grupos
- [ ] Módulo de permissões (admin, whitelist)
- [ ] Dashboard web para gerenciamento
- [ ] Integração com agenda e lembretes

---

## 📑 Observações Técnicas

- O pareamento via QR usa o pacote [`qrcode-terminal-go`](https://github.com/Baozisoftware/qrcode-terminal-go)
- O bot é compatível com PostgreSQL e SQLite, mas recomenda-se PostgreSQL para produção (evita `SQLITE_BUSY`)
- O schema do PostgreSQL deve ser aplicado manualmente com o arquivo `00-latest-schema.sql` do WhatsMeow

---

## 🤝 Contribuição

Pull requests são bem-vindos! Para grandes funcionalidades, crie uma issue para discutirmos.

---

## 🛡️ Licença

Este projeto está sob a licença MIT.

---

## 👤 Autor

Desenvolvido por [Renan Silva (Faysk)](https://github.com/faysk)  
📧 Contato: faysk.nan@gmail.com  
🌐 Projetos: https://faysk.top
