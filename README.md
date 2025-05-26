# 🤖 WhatsApp Bot - Modular, Inteligente e Escalável com Go

Este é um bot de WhatsApp altamente modular e escalável, construído com [WhatsMeow](https://github.com/tulir/whatsmeow) e linguagem Go. Projetado para ser limpo, poderoso e pronto para crescer com integrações de IA, voz e serviços externos.

---

## 🚀 Funcionalidades

- ✅ Conexão persistente via QR Code com o WhatsApp
- 💬 Respostas automáticas por comandos (ex: `!ping`, `!help`)
- 🧱 Arquitetura limpa e modularizada (eventos, comandos, serviços)
- 🧠 Integração com OpenAI GPT (mensagens com IA)
- 🗣️ **Suporte a áudio com edge-tts (texto para voz)**
- 🎙️ **Transcrição de áudios recebidos (Whisper)**
- 📦 Configuração via `.env` com exemplo incluído
- 🗂️ Banco de dados local via SQLite
- 🛠️ Script de instalação automática (`setup.sh`)
- 🔌 Fácil de integrar com Webhooks, APIs, CRON, etc.

---

## 📁 Estrutura do Projeto

```
whatsapp-bot/
├── cmd/             # Ponto de entrada
├── config/          # Leitura do .env
├── events/          # Dispatcher de eventos WhatsApp
├── handlers/        # Comandos e mensagens
│   └── commands/
├── services/        # Lógica de negócio (envio, IA, mídia)
├── openai/          # ChatGPT e Whisper
├── tts/             # Geração de voz com edge-tts
├── scripts/         # Scripts Python auxiliares (TTS, etc.)
├── store/           # Sessão WhatsApp e banco
├── utils/           # Logger customizado
├── media/           # Áudios gerados dinamicamente (.gitignored)
├── authorized.json  # Lista de números autorizados (exemplo)
├── .env.example     # Exemplo de configuração
├── setup.sh         # Instalador automático
├── go.mod / go.sum  # Dependências
└── README.md
```

---

## 📦 Pré-requisitos

- Go 1.20 ou superior → [https://go.dev/doc/install](https://go.dev/doc/install)
- Git instalado
- Python 3.10+ (para recursos de áudio via edge-tts)
- FFmpeg (para conversão de formatos de áudio, ex: `.opus`)
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
- Configurar variáveis de ambiente
- Instalar bibliotecas (`whatsmeow`, `godotenv`, `sqlite3`)
- Preparar o ambiente com `.env` baseado no `.env.example`

---

## 📄 Exemplo `.env`

```env
DB_PATH=file:session.db?_pragma=foreign_keys(1)
LOG_LEVEL=DEBUG
BOT_NAME=FayskBot
PORT=8080
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4o
MAX_TOKENS=2000
AUTHORIZED_NUMBERS=5511999999999,5511988888888
RESTRICT_TO_GROUP=false
```

---

## ▶️ Como Rodar

### Modo desenvolvimento

```bash
go run ./cmd
```

### Modo produção

```bash
go build -o bot
./bot
```

---

## 💬 Comandos Suportados

- `!ping` → Testa a resposta do bot
- `!help` → Lista os comandos disponíveis
- `bom dia` → Envia uma saudação especial
- (Em breve) `!falar <texto>` → Gera áudio com voz IA
- (Em breve) `!gpt <pergunta>` → Responde com GPT

---

## 🧠 Integrações de IA

### ✅ GPT via OpenAI
- Chat contextual por texto
- Mensagens personalizadas

### ✅ Whisper (áudio → texto)
- Transcrição de mensagens de voz para texto
- Multi-idioma (detectado automaticamente)

### ✅ edge-tts (texto → voz)
- Respostas faladas com voz neural da Microsoft
- Idiomas e vozes configuráveis (ex: Aria, Guy, Ana, etc.)

---

## 📌 Possibilidades Futuras

- [x] Gerar respostas em áudio automaticamente
- [x] Transcrever áudios recebidos
- [ ] Enviar imagens e stickers com IA
- [ ] Suporte a comandos em grupos
- [ ] Sistema de permissões (admin, whitelist, etc.)
- [ ] Integração com agenda e lembretes (Google Calendar, CRON)
- [ ] Dashboard web para gerenciar o bot

---

## 🤝 Contribuição

Pull requests são muito bem-vindos! Para funcionalidades maiores, abra uma issue primeiro para discutirmos juntos.

---

## 🛡️ Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais informações.

---

## 👤 Autor

Desenvolvido por [Renan Silva (Faysk)](https://github.com/faysk)  
📧 Contato: faysk@protonmail.com  
🌐 Projetos: https://faysk.dev