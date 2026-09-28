```text
    _    ___   ____  ____    _    __  __    _    ____ _____ _   _ ____ ___ ___  
   / \  |_ _| |  _ \|  _ \  / \  |  \/  |  / \  / ___|_   _| | | |  _ \_ _/ _ \ 
  / _ \  | |  | | | | |_) |/ _ \ | |\/| | / _ \ \___ \ | | | | | | | | | | | | |
 / ___ \ | |  | |_| |  _ </ ___ \| |  | |/ ___ \ ___) || | | |_| | |_| | | |_| |
/_/   \_\___| |____/|_| \_\_/   \_\_|  |_/_/   \_\____/ |_|  \___/|____/___\___/ 
```

# DramaStudio

> **AI-native autonomous drama production engine built with Domain-Driven Design (DDD) in Go.**

---

## Tech Stack & Ecosystem

![Go](https://img.shields.io/badge/Go_1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Temporal Workflow](https://img.shields.io/badge/Temporal-000000?style=for-the-badge&logo=temporal&logoColor=white)
![FFmpeg Engine](https://img.shields.io/badge/FFmpeg-007808?style=for-the-badge&logo=ffmpeg&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![OpenAI](https://img.shields.io/badge/OpenAI-412991?style=for-the-badge&logo=openai&logoColor=white)
![Anthropic](https://img.shields.io/badge/Anthropic-D97706?style=for-the-badge&logo=anthropic&logoColor=white)
![Google Gemini](https://img.shields.io/badge/Google_Gemini-4285F4?style=for-the-badge&logo=google-gemini&logoColor=white)

---

## About

This project is an **AI-native autonomous drama production studio** for creating serialized drama content from an initial story idea.

It is designed around the idea of a **virtual production studio powered by specialized AI model providers**. The system does not run large AI models locally or require users to provision GPUs. Instead, it orchestrates external AI services and model providers through a unified production pipeline.

A production can begin with a simple story idea and progress through:

```text
Idea
  ↓
Series Bible
  ↓
Story Architecture
  ↓
Season / Episode Planning
  ↓
Script
  ↓
Dialogue
  ↓
Character & World Design
  ↓
Storyboard
  ↓
Shot Planning
  ↓
Image Generation
  ↓
Video / Animation
  ↓
Voice / Music / SFX
  ↓
Assembly
  ↓
Continuity & Quality Control
  ↓
Final Render
  ↓
Publishing
```

### Specialized AI Models

The studio does not depend on a single AI model.

Each production capability can use the model or provider best suited for that task.

For example:

```text
Story Bible        → Story / reasoning model
Story Architecture → Reasoning model
Script Writing     → Script model
Dialogue Writing   → Dialogue model
Character Creation → Image / character model
Storyboard         → Vision / storyboard model
Image Generation   → Image model provider
Video Generation   → Video model provider
Animation          → Video / animation model
Voice              → Voice provider
Music              → Music provider
Continuity         → Analysis / reasoning model
```

Video generation may, for example, be powered by **Wan** through an external model provider. The production system does not need to know how the underlying model is hosted or whether the provider changes its infrastructure.

The same principle applies to every other AI capability.

### Provider-Agnostic Architecture

AI providers are accessed through standardized capability interfaces rather than being embedded directly into the production logic.

```text
Production Workflow
        ↓
AI Capability
        ↓
Model Selection / Policy
        ↓
Provider Adapter
        ↓
External AI Provider
        ↓
Generated Artifact
```

This allows the studio to use different providers for different capabilities and replace them without redesigning the production system.

The same ports/adapters rule applies to all infrastructure: bounded contexts depend on vendor-neutral ports (`ObjectStorage`, `WorkflowEngine`, `VideoGenerator`, `ImageGenerator`, `VoiceGenerator`, `ChannelAdapter`, `Authenticator`) and never on vendor names. Adapters are selected by configuration only:

```text
STORAGE_BACKEND=local|s3            # s3 covers AWS S3, Cloudflare R2, MinIO via S3_ENDPOINT
WORKFLOW_ENGINE=local|temporal      # local = in-process pipeline; temporal = durable worker
DRAMASTUDIO_STORE=sqlite|postgres|memory
DATABASE_PATH=data/dramastudio.db   # sqlite file (DRAMASTUDIO_STORE=sqlite)
DATABASE_URL=postgres://...         # postgres DSN (DRAMASTUDIO_STORE=postgres)
MODEL_CATALOG=configs/model-catalog.local.json  # zero-spend dev catalog (all-mock policies)
```

The three deployment profiles:

```text
local        DRAMASTUDIO_ENV=development (default)
             SQLite + local filesystem + local workflow engine + mock-capable providers
             → go run ./apps/api works with no servers and no API spend

test         DRAMASTUDIO_ENV=test
             SQLite (or ephemeral PostgreSQL) + fake providers + deterministic workflows

production   DRAMASTUDIO_ENV=production|staging
             PostgreSQL + S3-compatible storage + Temporal + real AI providers
```

SQLite is an explicit development/test adapter, not a production database:
`internal/platform/database/sqlite` implements the same `postgres.Querier`
port the repositories use, translating the dialect ($N placeholders, schema-
qualified names, JSONB containment) at the boundary — all fourteen repository
implementations run unchanged on both drivers, and `sqlite.Migrate` derives
the SQLite schema from the same `migrations/` files. Use PostgreSQL for
concurrency, multi-worker, and production-like testing.

Self-hosted / cloud modes:

```text
Self-hosted:    Go + PostgreSQL + MinIO        + FFmpeg + external AI APIs
Cloud:          Go + PostgreSQL + S3/R2        + FFmpeg + external AI APIs
```

**Rule: never let an infrastructure vendor become a domain concept.**

For example:

```text
Storyboard
    → Provider A
    → Model X

Character Generation
    → Provider B
    → Model Y

Video Generation
    → Provider C
    → Wan

Voice Generation
    → Provider D
    → Voice Model Z
```

The production system remains independent of those providers.

### Production Intelligence Layer

The pipeline consumes a structured intelligence subsystem rather than
embedding behavior in workflows or prompts. `internal/intelligence/` holds
five versioned, editable concepts — loaded from `catalogs/intelligence/`
as YAML data (see `INTELLIGENCE_DIR`):

- **Agents** (`agents/`) — role definitions: responsibilities, skills,
  policies, rules, permissions, context spec, model capability.
- **Skills** (`skills/`) — reusable methodology + I/O + evaluation criteria
  shared across agents (cinematography serves Video Director, Shot Designer,
  and Quality Supervisor alike).
- **Policies** (`policies/`) — layered operational constraints
  (system → studio → project → series → season → episode → task); lower
  layers cannot weaken keys a higher layer marked `protected`.
- **Rules** (`rules/`) — small deterministic predicates evaluated by the
  rules engine, never left to model interpretation (`severity` INFO→BLOCKING,
  `enforcement` ADVISORY→BLOCK).
- **Evaluators** (`evaluators/`) — independent output assessment:
  deterministic rule checks + optional model review.

Runtime flow: `Task → ResolveExecution → AgentExecutionContext (agent +
skills + effective policies + applicable rules + evaluators + permission
scope) → prompt assembly → model routing → post-rule validation →
evaluation → ExecutionRecord`. Every run audits agent/skill/policy/rule/
evaluator versions + provider/model/version + prompt and context hashes —
episode 40 stays reproducible after skills move to v3.

REST surface: `GET /v1/intelligence/{agents,skills,policies,rules,evaluators}`,
`POST /v1/intelligence/context` (preview), `POST /v1/intelligence/execute`,
`POST /v1/intelligence/policies/{id}/layers` (scoped override),
`GET /v1/projects/{id}/intelligence/executions` (audit).

### The AI Production Team

A **Lead Director Agent** coordinates the production rather than attempting to perform every task itself.

Specialized agents can be responsible for:

* Story development
* Series bible management
* Season planning
* Episode writing
* Dialogue
* Character development
* Visual direction
* Storyboarding
* Shot planning
* Video generation
* Voice direction
* Music and sound
* Editing and assembly
* Continuity
* Quality assurance
* Publishing

Agents are responsible for **planning, decisions, delegation, evaluation, and coordination**.

The underlying AI models are specialized tools used by those agents.

### Persistent Story Continuity

Long-running drama requires more than generating one video at a time.

A series may contain multiple seasons and potentially hundreds of episodes. The studio therefore maintains structured production knowledge including:

* Series Bible
* Characters
* Relationships
* Locations
* Props
* Wardrobe
* Timeline
* Story Graph
* Canonical Facts
* Character Knowledge
* Plot Threads
* Visual References
* Voice Profiles
* Production Assets
* Scene and Shot History

This information becomes the persistent source of truth for future generations.

A dialogue model does not need to receive an entire season to understand a scene. Instead, the production system constructs the relevant context from the canonical state.

### AI Creates. The Production System Controls.

The central architectural principle is:

> **AI creates the content. The production system owns reality.**

AI models generate creative artifacts.

The platform owns:

* Canon
* Continuity
* Production state
* Asset versions
* Workflows
* Agent decisions
* Approvals
* Dependencies
* Budgets
* Provider configuration
* Generation history
* Quality checks

This separation allows the studio to remain consistent even as individual models and providers change.

### No Local GPU Requirement

The platform is designed around **API-based AI model providers**.

It does not require users to install or operate local GPU infrastructure for model inference.

The studio handles the orchestration layer while external providers handle model execution.

That means the platform can focus on what it is actually building:

**an autonomous production system for AI-generated serialized drama.**

---

## Bounded Context Architecture

```text
                                  LEAD DIRECTOR
                                       │
                      ┌────────────────┼────────────────┐
                      │                │                │
                 STORY TEAM       VISUAL TEAM       PRODUCTION
                      │                │                │
                      ▼                ▼                ▼
               ┌─────────────┐  ┌──────────────┐  ┌──────────────┐
               │ Story Bible │  │ Character AI │  │ Animation AI │
               │ Model       │  │ Model        │  │ Model        │
               └─────────────┘  └──────────────┘  └──────────────┘
                      │                │                │
                      ▼                ▼                ▼
               ┌─────────────┐  ┌──────────────┐  ┌──────────────┐
               │ Story / Arc │  │ Storyboard   │  │ Video /      │
               │ Model       │  │ Model        │  │ Motion Model │
               └─────────────┘  └──────────────┘  └──────────────┘
                      │                │                │
                      ▼                │                ▼
               ┌─────────────┐         │         ┌──────────────┐
               │ Script      │         │         │ Voice Model  │
               │ Model       │         │         └──────────────┘
               └─────────────┘         │
                      │                ▼
                      ▼         ┌──────────────┐
               ┌─────────────┐  │ Shot / Camera│
               │ Dialogue    │  │ Model        │
               │ Model       │  └──────────────┘
               └─────────────┘
                      │
                      └────────────────┬────────────────┐
                                       ▼                ▼
                                 ┌───────────┐    ┌────────────┐
                                 │ Assembly  │    │ QA /       │
                                 │ Model     │    │ Continuity │
                                 └───────────┘    └────────────┘
```

---

## Directory Structure

```text
dramastudio/
├── apps/
│   ├── api/                  # Main HTTP API entrypoint
│   ├── worker/               # Background task worker process
│   ├── scheduler/            # Cron & background job scheduler
│   └── web/                  # Next.js 14 + React + TypeScript Web Studio UI
│
├── internal/
│   ├── platform/             # Cross-cutting infrastructure capabilities
│   │   ├── database/         # Postgres & transactional manager
│   │   ├── cache/            # Distributed caching contracts
│   │   ├── storage/          # Object storage abstractions
│   │   ├── events/           # Domain event bus & contracts
│   │   ├── messaging/        # Pub/Sub messaging
│   │   ├── workflow/         # Engine port + adapters (contracts/, temporal/, local/)
│   │   ├── ai/               # Low-level LLM / Embedding adapters
│   │   ├── media/            # Media processing utilities
│   │   ├── observability/    # Structured logging, metrics, tracing
│   │   ├── security/         # Authentication & authorization
│   │   └── clock/            # Time provider
│   │
│   ├── identity/             # Bounded Context: Users, Orgs, Permissions
│   ├── projects/             # Bounded Context: Projects & Production Policy
│   ├── story/                # Bounded Context: Narrative Hierarchy (Series -> Episode -> Beat)
│   ├── canon/                # Bounded Context: Authoritative Facts (Subject-Predicate-Object)
│   ├── characters/           # Bounded Context: Character Aggregates & Profiles
│   ├── world/                # Bounded Context: Locations, Props & Environments
│   ├── production/           # Bounded Context: Production Runs, Stages & Jobs
│   ├── agents/               # Bounded Context: AI Agents & Director Orchestrator
│   ├── ai/                   # Bounded Context: Model Registry & Capability Resolution
│   ├── continuity/           # Bounded Context: Multi-faceted Continuity Validation Engine
│   ├── media/                # Bounded Context: Asset Generation & Provider Adapters
│   ├── postproduction/       # Bounded Context: Composition, Subtitles & FFmpeg Rendering
│   ├── publishing/           # Bounded Context: Channel Distribution (TikTok/YouTube/IG)
│   └── analytics/            # Bounded Context: Production & Performance Analytics
│
├── migrations/               # Database SQL migrations
├── configs/                  # Configuration files
├── scripts/                  # DevOps & maintenance scripts
├── tests/                    # Integration & end-to-end test suite
├── go.mod
└── go.sum
```

---

## Getting Started

### Prerequisites

- **Go 1.22+** installed
- **Docker** & **PostgreSQL** (for local development)
- **FFmpeg** (for video rendering & timeline composition)

### Verification & Testing

Verify that all packages compile and run tests:

```bash
# Verify static code correctness
go vet ./...

# Run the test suite
go test ./...
```

---

## License

Copyright © 2026 DramaStudio. All rights reserved.
