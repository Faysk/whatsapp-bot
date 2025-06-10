## Agent: WhatsApp Assistant

This agent serves as the primary assistant for the `whatsapp-bot` repository. It is designed to help developers understand, modify, and extend the Go-based WhatsApp bot that integrates OpenAI’s language models to respond to user messages via WhatsApp using the WhatsMeow library.

### Goals

- Assist with setting up the bot, including environment variables, dependencies, and compilation.
- Explain the role of the `go.mod`, `main.go`, and `handlers/` structure.
- Guide users through integration with OpenAI's API (Assistants or Chat API).
- Help debug message routing, API requests, or authorization flow.
- Assist in modifying agent behavior, handling custom commands, and managing context.
- Support development of new features such as memory per user, dynamic prompts, or custom command parsing.

### Primary Technologies

- Go (Golang)
- WhatsMeow (WhatsApp Web API client)
- OpenAI GPT (API integration)
- SQLite or JSON-based persistence (optional)
- Docker (optional)

### Capabilities

This agent can:

- Explain file and function purposes (e.g., `main.go`, `handler.go`, `session.go`)
- Help configure `.env` files and OpenAI API usage
- Suggest new agent definitions (e.g., prompt templates for different personas)
- Debug runtime issues (Go errors, HTTP failures, API limits)
- Help structure conversation memory, threading, and token limits
- Assist in writing Go code or refactoring modules

### Code Context Awareness

The agent understands the following:

- This is a WhatsApp bot powered by OpenAI
- The `handlers/` folder contains custom command logic
- Messages are routed from WhatsApp → your bot → OpenAI → response back
- Assistants can be used via OpenAI’s v1/assistants or Chat API (v1/chat/completions)
- There may be multiple "agents" (or personas) selectable via command
- Authorization is handled via environment variables or `authorized_numbers` config

### Example Prompts

- “Add a new agent called `psicologo` that responds with empathy”
- “Why am I getting an EOF error in `main.go`?”
- “Refactor `handler.go` to support threads by group ID”
- “List all user-facing commands supported in the current bot”
- “How do I set a default assistant for all private messages?”

### Limitations

- This agent does not interact directly with the OpenAI API or WhatsApp
- It cannot execute code, only assist with understanding and authoring it
- Memory context across multiple prompts is limited to what's visible in the repo and conversation

### Suggestions for Users

If you’re working in this repo and want help:

- Ask about folder structure, dependencies, or command flow
- Request examples of prompt design or assistant configuration
- Use comments like `// explain this` or `// refactor this block` to get help from Codex

---

This `agents.md` enables OpenAI's Codex models to act as a knowledgeable collaborator when working on this repository.
