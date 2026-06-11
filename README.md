# ⚡ AI Incident Manager

An AI-powered IT incident management system that helps engineers diagnose, track, and resolve production incidents faster — built with **Go** and **Claude AI**.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)
![Claude](https://img.shields.io/badge/Claude-claude--sonnet--4--6-7c6cfa?style=flat)
![SQLite](https://img.shields.io/badge/SQLite-GORM-003B57?style=flat&logo=sqlite&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)

---

## What it does

When production breaks at 3 AM, you need answers fast. This system lets you:

- **Create an incident** → Claude automatically classifies severity (P1–P4) from the description
- **Run AI analysis** → get a structured root cause hypothesis, fix steps, and estimated resolution time
- **Chat with AI** → follow up with new symptoms or error messages and get contextual advice without re-explaining the incident

---

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                    Browser  (HTML / JS)                      │
│         Dashboard · Create Form · AI Chat Thread            │
└────────────────────────┬─────────────────────────────────────┘
                         │  REST API
┌────────────────────────▼─────────────────────────────────────┐
│                  Go Backend  (Gin)                           │
│                                                              │
│  ┌─────────────────┐   ┌──────────────────────────────────┐ │
│  │ IncidentHandler │   │          AI Handler              │ │
│  │  CRUD · Status  │   │  AnalyzeIncident()               │ │
│  └────────┬────────┘   │  ClassifySeverity()              │ │
│           │            │  Chat() — multi-turn conv.       │ │
│  ┌────────▼────────────▼──────────────────────────────┐   │ │
│  │                GORM  +  SQLite                      │   │ │
│  │         incidents.db  (incidents · messages)        │   │ │
│  └─────────────────────────────────────────────────────┘   │ │
└────────────────────────────────────────┬─────────────────────┘
                                         │  HTTPS
                           ┌─────────────▼──────────────┐
                           │    Anthropic Claude API     │
                           │    claude-sonnet-4-6        │
                           └────────────────────────────┘
```

---

## Features

| Feature | Description |
|---------|-------------|
| **Incident CRUD** | Create, update, close, and delete incidents |
| **Auto Severity Classification** | Claude reads the description and assigns P1–P4 on creation |
| **AI Analysis** | Structured output: root cause · fix steps · ETA |
| **Live Chat** | Multi-turn conversation — share follow-up errors and get contextual advice |
| **Persistent History** | Chat threads saved per incident, restored on revisit |
| **Stats Dashboard** | Real-time counts by severity and status |

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.21+ |
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) |
| Database | SQLite via [modernc/sqlite](https://gitlab.com/cznic/sqlite) (no CGO required) |
| AI | [Anthropic Claude API](https://docs.anthropic.com) — `claude-sonnet-4-6` |
| Frontend | Vanilla HTML / CSS / JS — no build step |

---

## Getting Started

### Prerequisites

- Go 1.21+
- An [Anthropic API key](https://console.anthropic.com)

### Installation

```bash
# 1. Clone
git clone https://github.com/your-username/ai-incident-manager.git
cd ai-incident-manager

# 2. Configure
cp .env.example .env
# Edit .env → set ANTHROPIC_API_KEY=sk-ant-...

# 3. Run
go run main.go
```

Open **http://localhost:8080**

---

## API Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/incidents` | List all incidents |
| `GET` | `/api/incidents/:id` | Get single incident |
| `POST` | `/api/incidents` | Create incident (AI classifies severity) |
| `PUT` | `/api/incidents/:id` | Update title / description / severity / status |
| `DELETE` | `/api/incidents/:id` | Delete incident |
| `POST` | `/api/incidents/:id/analyze` | Run structured AI analysis |
| `GET` | `/api/incidents/:id/messages` | Get chat history |
| `POST` | `/api/incidents/:id/chat` | Send message, get AI reply |

---

## Project Structure

```
ai-incident-manager/
├── main.go               # Entry point, router setup
├── handlers/
│   ├── incident.go       # CRUD handlers
│   ├── ai.go             # Claude API client, analyze + classify
│   └── chat.go           # Multi-turn chat handler
├── models/
│   ├── incident.go       # Incident model
│   └── message.go        # Chat message model
├── database/
│   └── db.go             # SQLite init + AutoMigrate
├── static/
│   └── index.html        # Single-page dashboard
├── .env.example
└── go.mod
```

---

## How the AI Integration Works

### Severity Classification
When you create an incident, the description is sent to Claude with a prompt asking it to return exactly one of `P1 / P2 / P3 / P4` based on business impact. No parsing gymnastics — the response is matched with a simple `strings.Contains`.

### Structured Analysis
The analysis prompt enforces a specific Markdown format in the response:

```
## Root Cause Hypothesis
## Suggested Fix Steps
## Estimated Resolution Time
```

This makes the output predictable and easy to render in the UI.

### Contextual Chat
Each chat message builds a full conversation history sent to Claude. The **system context** (incident title, severity, description) is injected into the very first user turn — so Claude always knows what incident it's debugging, even 20 messages deep.

---

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `ANTHROPIC_API_KEY` | Yes | Your Anthropic API key |

---

## License

MIT
