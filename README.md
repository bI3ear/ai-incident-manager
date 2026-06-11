# ⚡ AI Incident Manager

An AI-powered IT incident management system that helps engineers diagnose, track, and resolve production incidents faster — built with **Go** and **Claude AI**.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)
![Claude](https://img.shields.io/badge/Claude-claude--sonnet--4--6-7c6cfa?style=flat)
![SQLite](https://img.shields.io/badge/SQLite-GORM-003B57?style=flat&logo=sqlite&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-ready-2496ED?style=flat&logo=docker&logoColor=white)
![CI](https://github.com/bI3ear/ai-incident-manager/actions/workflows/ci.yml/badge.svg)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)

---

## What it does

When production breaks at 3 AM, you need answers fast. This system lets you:

- **Create an incident** → Claude automatically classifies severity (P1–P4) from the description
- **Run AI analysis** → get a structured root cause hypothesis, fix steps, and estimated resolution time
- **Chat with AI** → follow up with new symptoms or error messages and get contextual advice without re-explaining the incident

---

## Demo

> Create incident → AI classifies severity → Analyze → Chat to follow up

```
Engineer: "Payment service OOM crash every 2-3 hours since v2.4.1 deploy"
    ↓ auto-classified as P1

[Analyze with AI]
  → Root Cause: Event listener leak in new logging middleware
  → Fix Steps: 1) Check removeListener calls  2) Heap snapshot  3) ...
  → ETA: 2-4 hours

[Chat]
Engineer: "Disabled the middleware but still leaking, heap shows EventEmitter accumulating"
Claude:   "The leak is upstream of the middleware — check if your HTTP server
           is attaching a new listener on every request without cleanup..."
```

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
| **AI Analysis** | Structured output: root cause · fix steps · estimated resolution time |
| **Live Chat** | Multi-turn conversation — share follow-up errors, get contextual advice |
| **Persistent Chat History** | Conversations saved per incident, restored on revisit |
| **Stats Dashboard** | Real-time counts by severity and open status |
| **Docker Support** | Single-command startup via `docker compose up` |
| **CI/CD Pipeline** | Automated build, test, security scan, and Docker image publish on every push |

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.21+ |
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) |
| Database | SQLite via [modernc/sqlite](https://gitlab.com/cznic/sqlite) — no CGO required |
| AI | [Anthropic Claude API](https://docs.anthropic.com) — `claude-sonnet-4-6` |
| Frontend | Vanilla HTML / CSS / JS — no build step |
| Container | Docker + Docker Compose |
| CI/CD | GitHub Actions |

---

## Getting Started

### Option A — Docker (recommended)

```bash
# 1. Clone
git clone https://github.com/bI3ear/ai-incident-manager.git
cd ai-incident-manager

# 2. Configure
cp .env.example .env
# Edit .env → set ANTHROPIC_API_KEY=sk-ant-...

# 3. Run
docker compose up
```

Open **http://localhost:8080**

The database is persisted in a Docker volume — data survives container restarts.

### Option B — Local

```bash
# Requires Go 1.21+
cp .env.example .env   # set ANTHROPIC_API_KEY
go run main.go
```

---

## CI/CD Pipeline

```
Push / Pull Request
        │
        ├── Build & Test ──────────────────────────────────► go build
        │                                                    go vet
        │                                                    go test
        │
        └── Security Scan ─────────────────────────────────► govulncheck

Merge to main
        │
        └── Docker Build & Push ───────────────────────────► ghcr.io/bl3ear/ai-incident-manager:latest
                                                             ghcr.io/bl3ear/ai-incident-manager:sha-xxxxxxx
```

| Workflow | Trigger | Steps |
|----------|---------|-------|
| `ci.yml` | Every push & PR | `go build` → `go vet` → `go test` → `govulncheck` |
| `docker.yml` | Merge to `main` | Build image → Push to GHCR (tagged with SHA + `latest`) |

---

## API Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/incidents` | List all incidents |
| `GET` | `/api/incidents/:id` | Get single incident |
| `POST` | `/api/incidents` | Create incident — AI classifies severity |
| `PUT` | `/api/incidents/:id` | Update title / description / severity / status |
| `DELETE` | `/api/incidents/:id` | Delete incident |
| `POST` | `/api/incidents/:id/analyze` | Run structured AI analysis |
| `GET` | `/api/incidents/:id/messages` | Get chat history |
| `POST` | `/api/incidents/:id/chat` | Send message, get AI reply |
| `GET` | `/health` | Health check |

---

## Project Structure

```
ai-incident-manager/
├── main.go                        # Entry point, router setup
├── handlers/
│   ├── incident.go                # CRUD handlers
│   ├── ai.go                      # Claude API client, analyze + classify
│   ├── ai_test.go                 # Unit tests
│   └── chat.go                    # Multi-turn chat handler
├── models/
│   ├── incident.go                # Incident model
│   └── message.go                 # Chat message model
├── database/
│   └── db.go                      # SQLite init + AutoMigrate
├── static/
│   └── index.html                 # Single-page dashboard
├── .github/
│   └── workflows/
│       ├── ci.yml                 # Build, test, security scan
│       └── docker.yml             # Build and push Docker image
├── Dockerfile                     # Multi-stage build
├── docker-compose.yml             # Local dev with persistent volume
├── .env.example
└── go.mod
```

---

## How the AI Integration Works

### Severity Classification
When you create an incident, the description is sent to Claude with a prompt asking it to return exactly one of `P1 / P2 / P3 / P4` based on business impact. The response is matched with `strings.Contains` — simple and reliable.

### Structured Analysis
The analysis prompt enforces a specific Markdown format:

```
## Root Cause Hypothesis
## Suggested Fix Steps
## Estimated Resolution Time
```

This makes the output predictable and renders cleanly in the UI every time.

### Contextual Chat
Each message builds a full conversation history sent to Claude on every turn. The **incident context** (title, severity, description) is injected into the first user message — so Claude always knows what it is debugging, even 20 messages deep into a thread.

---

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `ANTHROPIC_API_KEY` | Yes | Your Anthropic API key from [console.anthropic.com](https://console.anthropic.com) |
| `DB_PATH` | No | SQLite file path (default: `incidents.db`) |

---

## License

MIT
