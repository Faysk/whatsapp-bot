# 🤖 WhatsApp Bot - Modular & Escalável com Go

Este projeto é um bot de WhatsApp totalmente modularizado usando [WhatsMeow](https://github.com/tulir/whatsmeow), com foco em arquitetura limpa, escalabilidade e facilidade de manutenção.

## 🚀 Recursos

- ✅ Conexão persistente com WhatsApp via QR Code
- 🧠 Respostas automáticas por comandos (`!ping`, `!help`, etc.)
- 🧱 Arquitetura modular (comandos, eventos, config, serviços)
- 📦 Configuração por `.env`
- ⚙️ Script de instalação e preparação do ambiente (`setup.sh`)
- 🗂️ Banco de dados local com SQLite
- 🧩 Pronto para expandir com IA, APIs, Webhooks e mais

---

## 📁 Estrutura do Projeto

```
whatsapp-bot/
├── cmd/               # Ponto de entrada da aplicação
│   └── main.go
├── config/            # Carregamento de variáveis de ambiente
│   └── config.go
├── store/             # Conexão com banco e cliente WhatsApp
│   ├── db.go
│   └── client.go
├── events/            # Escuta de eventos e roteamento
│   └── dispatcher.go
├── handlers/          # Tratamento de comandos de mensagem
│   ├── message.go
│   └── commands/
│       ├── ping.go
│       ├── help.go
│       └── bom_dia.go
├── services/          # Funções auxiliares de envio/resposta
│   └── whatsapp.go
├── utils/             # Logger customizado e utilidades
│   └── logger.go
├── .env               # Arquivo de ambiente (não versionar)
├── .env.example       # Exemplo de configuração
├── .gitignore
├── setup.sh           # Script automático de preparação
├── go.mod
├── go.sum
└── README.md
```

---

## 📦 Pré-requisitos

- Go 1.20 ou superior → [https://go.dev/doc/install](https://go.dev/doc/install)
- Git instalado (opcional)
- Terminal com permissões de execução

---

## ⚙️ Instalação Rápida

```bash
git clone https://github.com/seunome/whatsapp-bot.git
cd whatsapp-bot
chmod +x setup.sh
./setup.sh
```

Esse script irá:
- Verificar o Go instalado
- Inicializar `go.mod` (caso não exista)
- Instalar e organizar dependências (`godotenv`, `whatsmeow`, `sqlite`)
- Criar `.env` com base no `.env.example`

---

## 📄 Exemplo `.env`

```env
DB_PATH=file:session.db?_pragma=foreign_keys(1)
LOG_LEVEL=INFO
```

---

## ▶️ Como Rodar

### Desenvolvimento

```bash
go run ./cmd
```

### Produção

```bash
go build -o bot
./bot
```

---

## 💬 Comandos Suportados

- `!ping` → Testa a resposta do bot
- `!help` → Lista de comandos disponíveis
- `bom dia` → Resposta automática especial

---

## 📌 Possibilidades Futuras

- [ ] Respostas com mídia (imagem, áudio, stickers)
- [ ] Webhook para integração com serviços externos
- [ ] Comandos com IA (ChatGPT, DALL·E, etc.)
- [ ] Integração com banco PostgreSQL ou Redis
- [ ] Controle de permissões (admin, grupos)
- [ ] Sistema de agendamento de mensagens

---

## 🤝 Contribuição

Pull requests são bem-vindos! Para ideias maiores, abra uma issue para discussão antes.

---

## 🛡️ Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais informações.

---

## 👤 Autor

Desenvolvido por [Renan Silva (Faysk)](https://github.com/faysk)  
Contato: faysk@protonmail.com