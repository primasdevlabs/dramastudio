# AI Drama Studio — Backend Architecture

## 1. Overview

The backend is the **production control system** for an AI-native drama studio.

It is not an AI model server and it does not perform local GPU inference.

The platform orchestrates external AI model providers to transform an initial story idea into a complete serialized drama production.

A production may span:

- Multiple seasons
- Up to dozens of episodes per season
- Hundreds or thousands of scenes
- Thousands of shots
- Large numbers of generated images, videos, voices, music tracks, and other assets

The backend is therefore responsible for **orchestration, state, continuity, production management, provider integration, persistence, reliability, and observability**.

The core principle is:

> **AI creates the content. The production system owns reality.**

AI providers generate artifacts.

The backend owns:

- Story state
- Canon
- Character state
- World state
- Production state
- Asset versions
- Workflow state
- Agent decisions
- Approvals
- Continuity
- Quality checks
- Budgets
- Provider configuration
- Generation history
- Publishing state

---

# 2. Core Goals

The backend must provide:

### Production orchestration

Coordinate the complete lifecycle:

```text
Idea
→ Series Development
→ Season Planning
→ Episode Planning
→ Script
→ Dialogue
→ Scene Breakdown
→ Storyboard
→ Shot Planning
→ Character / Location Assets
→ Image Generation
→ Video Generation
→ Voice
→ Music / SFX
→ Assembly
→ QA
→ Continuity
→ Render
→ Publishing
```

### Long-running continuity

The system must support stories that span many episodes and seasons without treating episodes as isolated generations.

### Specialized AI models

Different capabilities may use different AI models and providers:

```text
Story Bible
Script
Dialogue
Character Creation
Storyboard
Image Generation
Video / Animation
Voice
Music
Continuity
Quality Evaluation
```

For example, video generation may use **Wan** through an external model provider.

The backend must not assume that Wan, or any other model, is locally hosted.

### Provider independence

AI providers must be replaceable without changing the production domain.

### Autonomous operation

The system must be capable of running production workflows with minimal human intervention.

### Human control

Users must be able to:

- Approve
- Reject
- Revise
- Regenerate
- Pause
- Resume
- Stop
- Override
- Retry
- Replace assets
- Change production settings

### Durable execution

Production jobs may take seconds, minutes, hours, or longer.

The backend must survive:

- Process restarts
- Network failures
- Provider failures
- Worker failures
- Temporary outages
- Long-running generations
- Human approval pauses

---

# 3. Technology Stack

The initial backend stack is:

```text
Language:          Go
Architecture:      DDD Modular Monolith
Workflow Engine:   Temporal
Database:          PostgreSQL
Cache:             Redis
Object Storage:    S3-compatible storage
Media Processing:  FFmpeg
API:               REST
Realtime:          SSE / WebSocket
Observability:     OpenTelemetry
```

External AI providers are integrated through provider adapters.

Conceptually:

```text
                     ┌──────────────────┐
                     │   Next.js Studio │
                     └────────┬─────────┘
                              │
                              ▼
                     ┌──────────────────┐
                     │    Go API        │
                     └────────┬─────────┘
                              │
                ┌─────────────┼─────────────┐
                ▼             ▼             ▼
          PostgreSQL       Redis        Temporal
                                              │
                         ┌────────────────────┼──────────────────┐
                         ▼                    ▼                  ▼
                    Story Workers       AI Workers        Media Workers
                         │                    │                  │
                         │                    ▼                  │
                         │          External AI Providers       │
                         │                    │                  │
                         │        ┌───────────┼──────────┐       │
                         │        ▼           ▼          ▼       │
                         │      LLM       Image       Video      │
                         │                  │          Wan        │
                         │                  │                    │
                         └──────────────────┴────────────────────┘
                                              │
                                              ▼
                                     S3-compatible Storage
```

---

# 4. Architectural Philosophy

## 4.1 Domain First

Business rules belong to domain modules.

The domain must not depend on:

- PostgreSQL
- Redis
- Temporal
- HTTP
- AI providers
- S3
- FFmpeg

Infrastructure implements interfaces defined by the domain/application layer.

---

# 4.2 Agents Are Not the Domain

Agents sit above the domain.

```text
AI Agent
    ↓
Application Service
    ↓
Domain
    ↓
Repository / Infrastructure
```

An AI agent should never directly manipulate database repositories to bypass business rules.

This allows an operation to be performed by:

- An AI agent
- A human operator
- An API request
- An automated workflow

without changing the underlying domain behavior.

---

# 4.3 Models Are Capabilities

The system must not embed model names into business logic.

Bad:

```go
if provider == "wan" {
    ...
}
```

Better:

```text
Capability: video_generation
        ↓
Model Policy
        ↓
Provider Adapter
        ↓
External Model
```

The production system requests a capability.

The AI subsystem determines which configured model/provider should execute it.

---

# 4.4 Workflows Are Durable

Production workflows must not rely on:

```text
goroutines
sleep()
custom retry loops
Redis-only queues
in-memory state
cron-only scheduling
```

Temporal is responsible for durable workflow execution.

---

# 5. Repository Structure

```text
drama-studio/
│
├── apps/
│   ├── api/
│   │   └── main.go
│   │
│   ├── worker/
│   │   └── main.go
│   │
│   └── scheduler/
│       └── main.go
│
├── internal/
│   │
│   ├── platform/
│   │   ├── database/
│   │   ├── cache/
│   │   ├── storage/
│   │   ├── events/
│   │   ├── messaging/
│   │   ├── workflow/
│   │   ├── ai/
│   │   ├── media/
│   │   ├── observability/
│   │   ├── security/
│   │   ├── clock/
│   │   └── configuration/
│   │
│   ├── identity/
│   ├── projects/
│   ├── story/
│   ├── canon/
│   ├── characters/
│   ├── world/
│   ├── agents/
│   ├── production/
│   ├── continuity/
│   ├── media/
│   ├── postproduction/
│   ├── publishing/
│   └── analytics/
│
├── migrations/
├── configs/
├── scripts/
├── tests/
├── docs/
├── go.mod
└── go.sum
```

---

# 6. Bounded Contexts

The backend is divided into bounded contexts.

```text
Identity
Projects
Story
Canon
Characters
World
Agents
Production
Continuity
Media
Postproduction
Publishing
Analytics
```

Each module owns its business rules and persistence boundary.

---

# 7. Module Structure

Each business module follows:

```text
module/
├── domain/
├── application/
│   ├── commands/
│   ├── queries/
│   └── services/
├── infrastructure/
└── interfaces/
    └── http/
```

Example:

```text
characters/
├── domain/
│   ├── character.go
│   ├── appearance.go
│   ├── personality.go
│   ├── wardrobe.go
│   ├── voice_profile.go
│   ├── relationship.go
│   ├── events.go
│   ├── repository.go
│   └── errors.go
│
├── application/
│   ├── commands/
│   ├── queries/
│   └── services/
│
├── infrastructure/
│   └── persistence/
│
└── interfaces/
    └── http/
```

---

# 8. Identity

Identity handles users, authentication, authorization, organizations, and access control.

Responsibilities:

- Users
- Organizations
- Memberships
- Roles
- Permissions
- Sessions
- API credentials
- Service identities

Example roles:

```text
Owner
Producer
Director
Writer
Editor
Reviewer
Operator
Viewer
```

Authorization should be permission-based rather than relying exclusively on role names.

---

# 9. Projects

A project represents a production.

```text
Project
├── Series
├── Seasons
├── Production Settings
├── AI Policies
├── Budget
├── Team
└── Publishing Configuration
```

A project may contain multiple seasons.

Project configuration includes:

```text
Title
Description
Genre
Language
Target Platforms
Aspect Ratio
Production Mode
Autonomy Level
Creative Rules
AI Policies
Budget
Publishing Rules
```

---

# 10. Story Context

The Story context owns narrative structure.

Hierarchy:

```text
Series
└── Season
    └── Arc
        └── Episode
            └── Scene
                └── Beat
```

The Story context owns:

- Story structure
- Episodes
- Scenes
- Beats
- Story objectives
- Conflicts
- Reveals
- Plot progression

It does not own canonical facts.

Canon is a separate context.

---

# 11. Series Bible

The Series Bible is the canonical creative foundation of a series.

It contains:

```text
Premise
Genre
Themes
Tone
World Rules
Narrative Rules
Character Foundations
Visual Style
Dialogue Style
Story Constraints
Continuity Rules
```

The Bible should be versioned.

```text
Bible v1
Bible v2
Bible v3
```

A production run records which Bible version was used.

This makes generations reproducible.

---

# 12. Canon Context

Canon is the authoritative source of story truth.

Example:

```text
John owns the restaurant.
```

This is a fact.

Canon should also track:

```text
Who knows?
When did they learn it?
When did it become true?
Is it still true?
Which episode established it?
```

Example:

```text
StoryFact
├── id
├── series_id
├── entity_id
├── type
├── value
├── introduced_episode
├── effective_from
├── effective_until
├── source
├── confidence
└── status
```

---

# 13. Knowledge State

Character knowledge must be separate from objective canon.

Example:

```text
Canon:
John is secretly married.

Sarah:
Does not know.

Audience:
Knows.

John:
Knows.
```

This prevents dialogue models from accidentally making characters reveal information they should not know.

---

# 14. Story Graph

The Story Graph represents relationships between narrative events.

Example:

```text
E03
Sarah discovers John's secret
        ↓
Sarah learns Fact X
        ↓
E04
Sarah confronts John
        ↓
E07
Sarah lies about knowing
        ↓
E14
Public reveal
```

Graph nodes can represent:

```text
Event
Reveal
Conflict
Relationship Change
Decision
Foreshadowing
Resolution
```

Edges represent relationships:

```text
causes
reveals
depends_on
contradicts
resolves
foreshadows
follows
```

---

# 15. Characters Context

Characters own persistent character identity.

A character contains:

```text
Identity
Appearance
Face Reference
Hair
Skin
Body
Signature Features
Personality
Relationships
Voice
Wardrobe
Story State
Knowledge
Visual Constraints
History
```

Character versions should be immutable once used in production.

Instead:

```text
Character
    ├── Version 1
    ├── Version 2
    └── Version 3
```

A shot references a specific version.

---

# 16. World Context

World contains:

```text
Locations
Buildings
Rooms
Layouts
Props
Environmental Rules
Visual References
```

A location may contain variants:

```text
Apartment
├── Living Room
│   ├── Day
│   ├── Night
│   └── Rain
```

The location identity remains stable while production variants can change.

---

# 17. Wardrobe

Wardrobe is part of continuity.

Example:

```text
Sarah
Episode 12
Scene 03

Black blouse
Blue jeans
White sneakers
```

If Scene 04 suddenly generates:

```text
Red evening dress
```

without a wardrobe change event, the continuity system should detect it.

---

# 18. Agents Context

The Agents context manages the virtual production team.

The hierarchy is:

```text
Lead Director
│
├── Story Director
├── Visual Director
├── Production Director
├── Continuity Supervisor
├── Quality Supervisor
└── Publishing Director
```

Specialized capabilities include:

```text
Bible Creation
Story Architecture
Season Planning
Episode Planning
Script Writing
Dialogue Writing
Character Creation
Storyboard Creation
Shot Planning
Image Generation
Video Generation
Voice
Music
Editing
Continuity
QA
Publishing
```

An agent is an orchestration entity.

A model is an execution resource.

---

# 19. Agent Definition

An agent definition contains:

```text
Identity
System Instructions
Skills
Tools
Permissions
Context Policy
Model Policy
Budget
Evaluation Criteria
```

Example:

```text
Story Director

Skills:
- story_analysis
- story_planning
- plot_development

Tools:
- read_canon
- read_story_graph
- create_story_plan
- request_approval

Model Policy:
- reasoning model

Permissions:
- story.read
- story.propose
- story.approve
```

---

# 20. Lead Director

The Lead Director is the primary orchestrator.

It should not be a giant class.

Its responsibilities are split into:

```text
Observe
Assess
Plan
Delegate
Monitor
Evaluate
Decide
Interrupt
```

Conceptual loop:

```text
Observe
   ↓
Assess
   ↓
Plan
   ↓
Delegate
   ↓
Monitor
   ↓
Review
   ↓
Pass?
 ┌─┴─┐
Yes  No
 │    │
 ▼    ▼
Next  Revise
       │
       └──→ Review
```

The Lead Director does not directly generate everything.

It delegates specialized work.

---

# 21. Agent Tasks

Every task must have explicit boundaries.

```text
AgentTask
├── Objective
├── Input
├── Expected Output
├── Constraints
├── Budget
├── Max Iterations
├── Timeout
├── Approval Policy
└── Success Criteria
```

Example:

```text
Task:
Generate storyboard for Episode 04 Scene 08

Input:
Approved script
Character references
Location reference
Visual bible

Expected Output:
Shot sequence

Constraints:
9:16
Canonical character appearance
Canonical location
Maximum 8 shots

Max Iterations:
3
```

If the task repeatedly fails, it escalates instead of looping indefinitely.

---

# 22. AI Capability System

AI capabilities provide the abstraction between production logic and external models.

Example capabilities:

```text
story_bible
story_architecture
season_planning
episode_planning
script_writing
dialogue_writing
character_creation
location_creation
storyboard_generation
shot_planning
image_generation
video_generation
voice_generation
music_generation
sfx_generation
continuity_analysis
quality_evaluation
```

---

# 23. Model Registry

The model registry maps capabilities to available providers/models.

```text
Capability
    ↓
Model Policy
    ↓
Provider
    ↓
Model
```

Example:

```text
video_generation
    ↓
provider_x
    ↓
wan
```

Another project may use:

```text
video_generation
    ↓
provider_y
    ↓
another_video_model
```

No production domain code changes.

---

# 24. Provider Adapters

External providers must be isolated.

```text
internal/platform/ai/
├── contracts/
├── registry/
├── routing/
└── providers/
    ├── llm/
    ├── image/
    ├── video/
    ├── voice/
    └── music/
```

Provider interface:

```go
type VideoGenerator interface {
    Generate(ctx context.Context, request VideoGenerationRequest) (
        VideoGenerationResult,
        error,
    )
}
```

The application layer should not know provider-specific request structures.

---

# 25. External AI Execution

A generation request follows:

```text
Agent
 ↓
Application Service
 ↓
AI Capability
 ↓
Model Resolver
 ↓
Provider Adapter
 ↓
External Provider API
 ↓
Generation Result
 ↓
Object Storage
 ↓
Asset Record
```

External providers own model inference.

The platform owns the resulting production artifact.

---

# 26. Generation Jobs

Every AI generation should be represented by a persistent job.

```text
GenerationJob
├── id
├── project_id
├── capability
├── provider
├── model
├── input
├── output
├── status
├── attempt
├── cost
├── started_at
├── completed_at
└── error
```

Status:

```text
PENDING
RUNNING
SUCCEEDED
FAILED
CANCELLED
```

---

# 27. Asset System

Assets are never silently overwritten.

Every generated artifact is versioned.

```text
Asset
├── id
├── type
├── project_id
├── character_id
├── scene_id
├── provider
├── model
├── prompt
├── reference_assets
├── parameters
├── version
├── status
├── cost
└── created_at
```

Asset types:

```text
Character Reference
Location Reference
Prop
Image
Video
Voice
Music
SFX
Subtitle
Storyboard
Render
```

---

# 28. Structured Generation Inputs

Prompts should not be the source of truth.

Instead, generation requests use structured specifications.

Example:

```json
{
  "subject": {
    "character": "sarah"
  },
  "location": "apartment_living_room",
  "time": "night",
  "emotion": "suspicious",
  "action": "looking_at_phone",
  "camera": {
    "shot": "medium_close_up",
    "angle": "eye_level"
  },
  "style": "series_canonical_style",
  "continuity": {
    "wardrobe": "wardrobe_sarah_ep12_v2",
    "props": ["phone_001"]
  }
}
```

A provider-specific adapter can transform this into whatever request format its API requires.

---

# 29. Prompt Construction

Prompt construction should happen near the AI capability layer.

The system should build prompts from:

```text
Series Bible
+
Character State
+
World State
+
Scene
+
Shot
+
Continuity Constraints
+
Creative Direction
```

This prevents giant hardcoded prompts from becoming the actual application architecture.

---

# 30. Context Retrieval

Models should receive only the context required for their task.

For example:

### Dialogue model

```text
Episode
Scene
Characters
Character Knowledge
Relationships
Scene Objective
Canon Facts
Previous Relevant Dialogue
Tone
Constraints
```

### Storyboard model

```text
Approved Script
Scene
Characters
Locations
Visual Bible
Shot Requirements
Camera Language
Continuity Constraints
```

### Video model

```text
Approved Shot
Storyboard
Character Reference
Location Reference
Motion Specification
Camera Movement
Duration
Aspect Ratio
Previous/Next Shot References
```

This is critical for long-running productions.

---

# 31. Production Context

Production owns executable production state.

It tracks:

```text
Episodes
Scenes
Shots
Production Jobs
Generation Jobs
Approvals
Production Runs
Dependencies
```

Example:

```text
Episode
├── Scene
│   ├── Shot
│   │   ├── Image Generation
│   │   ├── Video Generation
│   │   └── Audio
│   └── ...
└── ...
```

---

# 32. Production State Machine

A project may progress through:

```text
PROJECT_CREATED
        ↓
DEVELOPING_SERIES
        ↓
SERIES_APPROVED
        ↓
PLANNING_SEASON
        ↓
PLANNING_EPISODES
        ↓
PRODUCING_EPISODE
        ↓
VALIDATING_EPISODE
        ↓
POST_PRODUCTION
        ↓
READY
        ↓
PUBLISHED
```

Individual episodes can have their own state machine.

---

# 33. Temporal Workflows

Temporal owns long-running production workflows.

Example:

```text
ProduceEpisodeWorkflow
│
├── GenerateEpisodePlan
├── ValidateContinuity
├── GenerateScript
├── GenerateDialogue
├── GenerateScenePlans
├── GenerateStoryboard
├── GenerateCharacterAssets
├── GenerateLocationAssets
├── GenerateVideo
├── GenerateVoice
├── GenerateMusic
├── AssembleEpisode
├── RunContinuityChecks
├── RunQualityChecks
└── RequestApproval
```

Temporal provides:

- Durable execution
- Retry policies
- Timers
- Signals
- Queries
- Workflow history
- Failure recovery
- Human approval pauses

---

# 34. Human Approval

Approval is a first-class workflow primitive.

Example:

```text
Workflow
   ↓
Approval Requested
   ↓
Workflow Paused
   ↓
User Decision
   ↓
Signal Workflow
   ↓
Continue
```

Possible decisions:

```text
APPROVE
REJECT
REQUEST_REVISION
REGENERATE
PAUSE
STOP
OVERRIDE
```

---

# 35. Production Modes

## Monitored

AI performs production but requests human intervention at configured decision points.

Example:

```text
Story Direction
      ↓
Approval Required
      ↓
User Approves
      ↓
Production Continues
```

## Autonomous

The Lead Director makes decisions according to:

```text
Creative Rules
Production Rules
Budget
Quality Thresholds
Approval Policy
```

The decision is recorded.

---

# 36. Decision Records

Important decisions must be auditable.

```text
Decision
├── id
├── project_id
├── episode_id
├── decision
├── reason
├── decision_maker
├── mode
└── timestamp
```

Decision makers can be:

```text
USER
LEAD_DIRECTOR
SPECIALIZED_AGENT
SYSTEM
```

---

# 37. Continuity Context

Continuity validates the production against the source of truth.

It should not own the original story data.

It reads from:

```text
Story
Canon
Characters
World
Production
Assets
Timeline
```

and produces:

```text
ContinuityIssue
```

---

# 38. Continuity Checks

The system should validate:

### Story

- Contradictory events
- Broken plot dependencies
- Missing resolutions
- Incorrect reveals

### Character

- Character state
- Knowledge
- Relationships
- Personality
- Age
- Injuries
- Appearance

### Visual

- Face
- Hair
- Clothing
- Body
- Location
- Props
- Visual style

### Timeline

- Event order
- Time of day
- Travel
- Duration
- Scene sequencing

### Production

- Missing assets
- Incorrect references
- Missing shots
- Broken dependencies

---

# 39. Continuity Issues

Avoid a single vague continuity score.

Use actionable issues:

```text
ContinuityIssue
├── category
├── severity
├── entity
├── evidence
├── expected_state
├── actual_state
├── resolution
└── status
```

Severity:

```text
INFO
WARNING
ERROR
BLOCKING
```

Example:

```text
Category:
Wardrobe

Expected:
Sarah — black blouse

Actual:
Sarah — red dress

Cause:
No wardrobe change event found.

Severity:
BLOCKING
```

---

# 40. Timeline Engine

The timeline tracks narrative time independently from production time.

```text
TimelineEvent
├── world_time
├── episode
├── scene
├── event_order
├── participants
├── location
└── duration
```

This allows the system to detect impossible sequences.

Example:

```text
10:00 — John leaves Lagos
10:05 — John appears in Cotonou
```

If the established travel constraints make that impossible, the system can flag the inconsistency.

---

# 41. Postproduction

Postproduction converts approved production assets into a final episode.

AI may propose:

```text
Cuts
Transitions
Timing
Shot Selection
Music Placement
Audio Levels
Caption Timing
```

But actual media assembly should be deterministic where possible.

Recommended architecture:

```text
Editor Agent
      ↓
Edit Decision List / Timeline
      ↓
Postproduction Engine
      ↓
FFmpeg / Media Processing
      ↓
Render
```

The AI does not need to directly manipulate raw media.

---

# 42. Assembly

Assembly should be represented as structured data.

Example:

```text
Timeline
├── Video Track
│   ├── Shot 01
│   ├── Shot 02
│   └── Shot 03
│
├── Dialogue Track
├── Music Track
├── SFX Track
└── Subtitle Track
```

The timeline becomes the source of truth for rendering.

---

# 43. Media Processing

FFmpeg handles deterministic media operations:

```text
Trim
Concat
Scale
Crop
Mux
Audio Mix
Subtitles
Transcode
Thumbnail
Waveform
Extract Frames
Generate Preview
```

Heavy media processing should run in workers.

---

# 44. Object Storage

Large artifacts should not be stored in PostgreSQL.

PostgreSQL stores metadata.

S3-compatible storage stores:

```text
Images
Videos
Audio
Storyboards
Renders
Thumbnails
Subtitles
Preview Files
```

Example:

```text
projects/{project_id}/
    characters/
    locations/
    episodes/
        {episode_id}/
            scenes/
            shots/
            audio/
            renders/
```

The exact storage layout should be generated from centralized configuration rather than duplicated throughout the codebase.

---

# 45. Redis

Redis is an infrastructure service, not the source of production truth.

Use Redis for:

```text
Caching
Rate Limiting
Short-lived coordination
Distributed locks where required
Realtime support
Performance optimization
```

Do not use Redis as the authoritative state for:

```text
Episodes
Canon
Assets
Agent decisions
Production state
Approvals
```

Those belong in PostgreSQL and/or Temporal workflow state.

---

# 46. PostgreSQL

PostgreSQL stores durable business state.

Major data areas include:

```text
identity
projects
story
canon
characters
world
agents
production
continuity
media
postproduction
publishing
analytics
```

Module ownership must remain clear.

A module must not directly manipulate another module's tables.

---

# 47. Module Communication

Modules communicate through:

### Commands

Request a state change.

```text
CreateEpisode
ApproveStoryboard
GenerateCharacter
ResolveContinuityIssue
```

### Queries

Read information.

```text
GetEpisode
GetCharacter
GetCanonFacts
GetStoryboard
```

### Domain Events

Communicate completed facts.

```text
EpisodeCreated
CharacterLocked
StoryFactEstablished
AssetGenerated
ContinuityIssueDetected
```

Avoid:

```text
characters repository → production repository
```

Prefer:

```text
Characters
    ↓
CharacterLocked event
    ↓
Production reacts
```

---

# 48. Domain Events

Core events include:

```text
SeriesCreated
SeasonCreated
ArcCreated
EpisodeCreated
EpisodeApproved
EpisodeCompleted

CharacterCreated
CharacterUpdated
CharacterLocked

StoryFactEstablished
StoryFactChanged

SceneCreated
SceneApproved

AssetGenerated
AssetApproved

ProductionJobCreated
ProductionJobCompleted
ProductionJobFailed

ContinuityIssueDetected
ContinuityIssueResolved

AgentTaskCreated
AgentTaskCompleted
AgentDecisionMade

ApprovalRequested
ApprovalGranted
ApprovalRejected
```

Events should be versioned.

---

# 49. Idempotency

External providers and workflows can retry.

Every externally visible operation that can safely be retried should support idempotency.

Example:

```text
Idempotency-Key:
production-job-123-generation-2
```

A retry must not accidentally:

- Generate duplicate production records
- Charge twice
- Create duplicate assets
- Publish twice
- Advance a workflow incorrectly

---

# 50. Failure Handling

Failures are expected.

Categories:

```text
Provider Failure
Network Failure
Timeout
Rate Limit
Invalid Request
Content Rejection
Worker Failure
Database Failure
Storage Failure
Workflow Failure
Quality Failure
Continuity Failure
```

Not every failure should be retried.

Example:

```text
Timeout
→ retry

429 rate limit
→ retry with provider policy

Invalid prompt
→ do not blindly retry

Continuity failure
→ revise / escalate

Human rejection
→ wait for new instruction
```

---

# 51. Retry Policy

Retries belong to workflow/application infrastructure.

Every retry should have:

```text
Maximum Attempts
Backoff
Timeout
Retryable Errors
Non-Retryable Errors
```

Never create infinite autonomous loops.

---

# 52. Budget Management

AI generation can become expensive.

The backend should track:

```text
Project Budget
Season Budget
Episode Budget
Task Budget
Provider Cost
Generation Cost
Media Processing Cost
```

Every generation records its cost where the provider makes cost information available.

The Lead Director should receive budget information when making production decisions.

---

# 53. Model Routing

Model routing should support policy.

Example:

```text
Capability:
video_generation

Policy:
provider = X
model = Wan
max_cost_per_generation = ...
allowed_duration = ...
aspect_ratios = ...
```

Policies may be scoped to:

```text
System
Organization
Project
Season
Episode
Task
```

The more specific policy overrides broader policy where explicitly configured.

No silent fallback should be introduced.

If a required model is unavailable, the system should report the failure or request an explicit policy change.

---

# 54. Configuration

Configuration should be explicit.

Avoid:

```go
timeout := config.Timeout
if timeout == 0 {
    timeout = 30
}
```

when `0` could represent a meaningful configuration state.

Configuration should be validated during startup.

Required configuration should fail clearly rather than silently falling back.

---

# 55. API Design

The API is REST-oriented.

Example:

```text
/v1/projects
/v1/projects/{projectId}

/v1/projects/{projectId}/story
/v1/projects/{projectId}/seasons
/v1/projects/{projectId}/episodes

/v1/projects/{projectId}/characters
/v1/projects/{projectId}/locations

/v1/projects/{projectId}/production
/v1/projects/{projectId}/jobs

/v1/projects/{projectId}/continuity

/v1/projects/{projectId}/agents
/v1/projects/{projectId}/decisions

/v1/projects/{projectId}/assets

/v1/projects/{projectId}/publishing
```

---

# 56. API Principles

APIs should:

- Validate input
- Authorize access
- Execute application commands
- Return domain-safe responses
- Never expose infrastructure internals
- Never expose provider-specific structures unnecessarily
- Support pagination
- Support filtering
- Support idempotency where appropriate
- Return consistent errors

---

# 57. API Error Format

Use a consistent structure:

```json
{
  "error": {
    "code": "STORYBOARD_NOT_APPROVED",
    "message": "The storyboard must be approved before video generation.",
    "details": {},
    "request_id": "req_123"
  }
}
```

Internal errors should not expose:

- Provider credentials
- Database details
- Stack traces
- Internal service topology

---

# 58. Realtime

The studio requires realtime production updates.

Use:

```text
SSE
```

for primarily server → client event streams.

Use:

```text
WebSocket
```

when bidirectional realtime communication is actually required.

Events include:

```text
AgentStarted
AgentThinking / Progress
TaskStarted
TaskCompleted
GenerationStarted
GenerationCompleted
ContinuityIssueDetected
ApprovalRequested
RenderCompleted
```

Realtime events should not replace durable database state.

The UI can reconnect and rebuild its state from the API.

---

# 59. Security

Security is part of the architecture.

Requirements include:

```text
Authentication
Authorization
Project Isolation
Organization Isolation
Secret Management
Provider Credential Encryption
Signed Storage URLs
Rate Limiting
Audit Logs
Input Validation
Webhook Verification
Idempotency
Request Tracing
```

AI provider credentials must never be exposed to the browser.

---

# 60. Provider Credentials

Credentials belong exclusively to the backend.

```text
Browser
    ✗
Provider API Key

Browser
    ↓
Go Backend
    ↓
Encrypted Credential
    ↓
Provider
```

Frontend receives capability information, not provider secrets.

---

# 61. Webhooks

External providers may use asynchronous generation.

Flow:

```text
Backend
 ↓
Provider
 ↓
Generation Started
 ↓
Provider Webhook
 ↓
Webhook Verification
 ↓
Generation Job Updated
 ↓
Temporal Signal/Event
 ↓
Workflow Continues
```

Webhook handlers must be:

- Authenticated
- Signature verified
- Idempotent
- Fast
- Observable

Heavy work should happen asynchronously.

---

# 62. Observability

Every production operation should be traceable.

Use OpenTelemetry for:

```text
Traces
Metrics
Logs
```

Trace hierarchy:

```text
Production
 └── Episode
      └── Workflow
           └── Agent Task
                └── AI Generation
                     └── Provider Request
```

This makes it possible to answer:

> Why is Episode 12 stuck?

or:

> Which provider generated this shot?

or:

> How much did this episode cost?

---

# 63. Audit Logging

Record important actions:

```text
User approvals
User overrides
Agent decisions
Model changes
Provider changes
Asset approvals
Asset replacements
Continuity resolutions
Publishing actions
Configuration changes
```

Audit logs should be append-only from the application perspective.

---

# 64. Quality System

Quality is not one AI score.

The system should use multiple validators.

```text
Narrative QA
Visual QA
Audio QA
Continuity QA
Technical QA
Policy / Provider QA
```

Example:

```text
Episode
 ↓
Narrative QA
 ↓
Visual QA
 ↓
Continuity QA
 ↓
Technical QA
 ↓
Ready
```

Each check produces structured findings.

---

# 65. Quality Evaluation

AI evaluators may assess:

```text
Script coherence
Dialogue quality
Character consistency
Storyboard completeness
Visual consistency
Audio synchronization
Subtitle accuracy
Scene continuity
```

However, deterministic checks should be preferred where possible.

For example:

```text
Video duration
Resolution
Frame rate
Audio presence
Subtitle timing
Missing assets
```

should not require an LLM.

---

# 66. Publishing

Publishing is a separate context.

It manages:

```text
Channels
Accounts
Metadata
Captions
Thumbnails
Publication Status
Scheduling
Platform Responses
```

The production system should not assume that every generated episode is automatically published.

Publishing should obey project policy and approval requirements.

---

# 67. Analytics

Analytics is separate from production state.

It can collect:

```text
Views
Watch Time
Retention
Engagement
Shares
Comments
Follower Growth
Publication Performance
```

Analytics may later influence production decisions, but it should not directly mutate story canon.

---

# 68. Production Workflow Example

A user creates:

```text
Series:
"The Last Promise"
```

The Lead Director begins:

```text
Series Development
```

### Step 1 — Bible

The Bible model creates:

```text
Premise
Characters
World
Themes
Rules
Visual Direction
```

The result is stored and versioned.

### Step 2 — Season

The Story Director plans:

```text
Season 1
├── Arc 1
├── Arc 2
└── Arc 3
```

### Step 3 — Episodes

The season planner creates:

```text
Episode 01
Episode 02
...
Episode 50
```

### Step 4 — Script

The Script capability generates the episode script.

### Step 5 — Dialogue

The Dialogue capability works from:

```text
Script
Character State
Knowledge State
Relationships
Canon
```

### Step 6 — Storyboard

The Storyboard capability creates:

```text
Scene
→ Shot
→ Camera
→ Action
→ Character
→ Location
```

### Step 7 — Character/Location References

The appropriate image/character models generate reference assets.

### Step 8 — Animation

The Video capability submits shots to the configured provider.

For example:

```text
Video Capability
    ↓
Provider Adapter
    ↓
Wan
```

### Step 9 — Audio

Voice, music, and SFX providers generate their assets.

### Step 10 — Assembly

The Editor Agent creates a structured edit plan.

```text
Edit Plan
 ↓
FFmpeg
 ↓
Episode Render
```

### Step 11 — QA

The system performs:

```text
Technical QA
Narrative QA
Continuity QA
Visual QA
Audio QA
```

### Step 12 — Approval

If monitored mode requires approval:

```text
Approval Requested
```

Otherwise:

```text
Lead Director Decision
```

### Step 13 — Publishing

The Publishing workflow prepares the final artifact for configured channels.

---

# 69. Continuity Across Seasons

The most important architectural rule is:

> **Never treat an episode as an isolated AI generation.**

Every episode reads from persistent canonical state.

After Episode 01:

```text
Character State
Canon
Timeline
Relationships
Plot Threads
```

are updated.

Episode 02 consumes the resulting state.

Episode 50 consumes the accumulated canonical state.

Season 2 consumes Season 1's established state.

This enables:

```text
Season 1
  ↓
Canon
  ↓
Season 2
  ↓
Canon
  ↓
Season 3
```

---

# 70. Production Snapshots

For reproducibility, important production stages should create snapshots.

Example:

```text
Episode 12 Production Snapshot

Bible: v4
Canon: v18
Characters: v31
World: v12
Storyboard: v7
Model Policy: v9
```

If an issue occurs later, the production can determine exactly which state generated an artifact.

---

# 71. Reproducibility

A generation should be traceable to:

```text
Project
Episode
Scene
Shot
Task
Agent
Capability
Model
Provider
Prompt/Input Specification
Reference Assets
Canon Version
Character Version
World Version
Model Policy Version
```

This is essential for debugging and regeneration.

---

# 72. Regeneration

Regeneration should create a new version.

Never overwrite the previous result.

```text
Shot 07

Generation 1 → rejected
Generation 2 → rejected
Generation 3 → approved
```

All generations remain auditable.

The production graph points to the currently approved artifact.

---

# 73. Cancellation

Cancellation must be supported at multiple levels:

```text
Generation
Task
Scene
Episode
Season
Project
```

Temporal workflows should receive cancellation signals where appropriate.

External provider jobs should be cancelled where supported.

If a provider does not support cancellation, the system should mark the generation as cancelled locally and prevent its result from automatically advancing the production.

---

# 74. Human Override

A user can override an AI decision.

Example:

```text
Lead Director:
Use Character Reference v4.

Human:
Use v3 instead.
```

The override becomes a recorded decision.

The workflow continues using the explicitly selected state.

---

# 75. No God Modules

Avoid:

```text
production_service.go
agent_service.go
ai_service.go
```

containing hundreds or thousands of lines.

Split responsibilities.

For example:

```text
agents/director/
├── observe.go
├── assess.go
├── plan.go
├── delegate.go
├── evaluate.go
├── decide.go
└── interrupt.go
```

Likewise:

```text
production/
├── episode/
├── scene/
├── shot/
├── jobs/
└── workflows/
```

---

# 76. No Hardcoded Provider Logic

Avoid:

```go
if model == "wan" {
    ...
}

if provider == "provider-x" {
    ...
}
```

Provider-specific behavior belongs inside adapters.

```text
AI Capability
    ↓
Provider Contract
    ↓
Provider Adapter
```

---

# 77. No Hardcoded Creative Defaults

Creative behavior should come from:

```text
Series Bible
Project Configuration
Production Policy
Agent Configuration
Model Policy
Task Constraints
```

Avoid silently inventing:

```text
genre
tone
aspect ratio
language
model
provider
style
budget
```

when the system has not been configured to do so.

---

# 78. Testing Strategy

Testing should occur at multiple levels.

## Unit Tests

Test:

```text
Domain Entities
Value Objects
Business Rules
State Transitions
Continuity Rules
Authorization Rules
```

## Application Tests

Test:

```text
Commands
Queries
Services
Module Interactions
```

## Integration Tests

Test:

```text
PostgreSQL
Redis
S3
Provider Adapters
Temporal
```

## Workflow Tests

Temporal workflow tests should cover:

```text
Success
Retry
Timeout
Provider Failure
Human Approval
Rejection
Cancellation
Resume
Escalation
```

## Contract Tests

Provider adapters should verify that they satisfy the platform's AI capability contracts.

---

# 79. Development Environment

Local development should provide infrastructure services without requiring local AI inference.

Example:

```text
PostgreSQL
Redis
Temporal
S3-compatible storage
FFmpeg
Go API
Temporal Workers
Next.js frontend
```

AI calls can target configured external providers.

Mock providers should be available for automated tests.

---

# 80. Mock AI Providers

A mock provider should implement the same contract as a real provider.

Example:

```text
VideoGenerator
├── WanProvider
├── ProviderXVideoProvider
└── MockVideoProvider
```

Tests can therefore run without spending money on external AI APIs.

---

# 81. Deployment Architecture

Initial deployment can remain a modular monolith.

```text
                   Load Balancer
                        │
                        ▼
                  Go API Instances
                        │
           ┌────────────┼────────────┐
           ▼            ▼            ▼
      PostgreSQL      Redis       Temporal
                                     │
                                     ▼
                              Worker Instances
                                     │
                   ┌─────────────────┼────────────────┐
                   ▼                 ▼                ▼
              AI Workers       Media Workers    Story Workers
                   │
                   ▼
           External Providers
```

The application remains logically modular while infrastructure can scale independently.

---

# 82. Scaling

Scale based on workload.

AI generation workers may need more capacity than story workers.

Media workers may need different resources than API workers.

Temporal workers can be split by task queues:

```text
story
ai
image
video
audio
media
qa
publishing
```

Example:

```text
video task queue
    ↓
video workers
    ↓
video providers
```

This prevents video workloads from starving unrelated workflows.

---

# 83. Task Queue Isolation

Recommended Temporal task queues:

```text
story-tasks
visual-tasks
video-tasks
audio-tasks
media-tasks
qa-tasks
publishing-tasks
```

Worker pools can scale independently.

---

# 84. Rate Limiting

External provider limits must be respected.

Rate limiting should consider:

```text
Provider
Model
Capability
Organization
Project
```

The provider adapter should expose normalized errors such as:

```text
ProviderRateLimited
ProviderUnavailable
ProviderTimeout
ProviderRejected
```

The workflow decides whether and when to retry.

---

# 85. Cost Control

Before expensive generation:

```text
Task
 ↓
Estimate
 ↓
Budget Check
 ↓
Generate
```

A task may be blocked if:

```text
Episode budget exceeded
Project budget exceeded
Provider spending limit reached
```

The system should request an explicit decision rather than silently switching providers or models.

---

# 86. Architecture Dependency Rules

The following dependency direction should be enforced:

```text
Interfaces
    ↓
Application
    ↓
Domain

Infrastructure
    ↓
implements interfaces/contracts
```

Domain must never depend on:

```text
HTTP
PostgreSQL
Redis
Temporal
S3
Provider SDKs
FFmpeg
```

Agents may depend on application interfaces.

---

# 87. Backend Architectural Boundary

The final conceptual architecture is:

```text
┌─────────────────────────────────────────────┐
│                 NEXT.JS STUDIO              │
└──────────────────────┬──────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────┐
│                   GO API                    │
├─────────────────────────────────────────────┤
│ Identity                                     │
│ Projects                                     │
│ Story                                        │
│ Canon                                        │
│ Characters                                   │
│ World                                        │
│ Agents                                       │
│ Production                                   │
│ Continuity                                   │
│ Media                                        │
│ Postproduction                               │
│ Publishing                                   │
│ Analytics                                    │
└──────────────────────┬──────────────────────┘
                       │
             ┌─────────┴─────────┐
             ▼                   ▼
       PostgreSQL            Temporal
             │                   │
             │             ┌─────┼─────┐
             │             ▼     ▼     ▼
             │          Story   AI    Media
             │         Workers Workers Workers
             │                   │
             │                   ▼
             │           External AI Providers
             │             ┌─────┼─────┐
             │             ▼     ▼     ▼
             │           LLM   Image Video
             │                         │
             │                        Wan
             │
             └──────────────┬──────────────┐
                            ▼              ▼
                          Redis       S3 Storage
```

---

# 88. Core Principle

The backend is not an AI chatbot with a video-generation button.

It is a **durable production operating system** that coordinates specialized AI capabilities and external model providers.

The AI models are replaceable.

The providers are replaceable.

The agents are replaceable.

The production truth is not.

```text
                EXTERNAL AI
                    │
         ┌──────────┼──────────┐
         ▼          ▼          ▼
       Story      Visual     Audio
         │          │          │
         └──────────┼──────────┘
                    ▼
             AGENT SYSTEM
                    │
                    ▼
             PRODUCTION CORE
                    │
          ┌─────────┼─────────┐
          ▼         ▼         ▼
        CANON    CONTINUITY  STATE
          │         │         │
          └─────────┼─────────┘
                    ▼
              FINAL PRODUCTION
```

**The models generate.
The agents coordinate.
Temporal executes.
The domain owns truth.
The production system remembers.**
