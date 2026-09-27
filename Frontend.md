# AI Drama Studio — Frontend Architecture

## 1. Overview

The frontend is the **creative production studio and control room** for the AI Drama Studio.

It is not a traditional SaaS dashboard.

The application provides a persistent workspace where users can:

- Create and develop drama series
- Manage series bibles
- Develop seasons and episodes
- Review scripts and dialogue
- Design characters and locations
- Build and review storyboards
- Manage shots and generated media
- Monitor AI production
- Review continuity
- Approve or reject AI decisions
- Control autonomous production
- Assemble and review episodes
- Manage publishing
- Inspect production history and analytics

The frontend communicates with the Go backend and never communicates directly with AI model providers.

```text
User
 ↓
Next.js Studio
 ↓
Go API
 ↓
Production / Agent System
 ↓
External AI Providers
```

The frontend's responsibility is to provide the **human-facing production environment**.

The backend remains the source of truth.

---

# 2. Core Principle

> **The frontend is the studio. The backend is the production system.**

The frontend should never become the source of truth for:

- Canon
- Production state
- Agent state
- Workflow state
- Asset state
- AI provider configuration
- Generation status
- Approval state
- Production decisions

Those belong to the backend.

The frontend displays, edits, commands, and observes that state.

---

# 3. Technology Stack

```text
Framework:       Next.js
UI:              React
Language:        TypeScript
Styling:         Tailwind CSS
Server State:    TanStack Query
UI State:        Zustand
Validation:      Zod
Forms:           React Hook Form
Realtime:        SSE / WebSocket
Icons:           Lucide
Media:           Browser-native media APIs + dedicated viewers
```

The frontend should remain provider-agnostic.

It should not contain logic such as:

```ts
if (provider === "wan") {
    ...
}
```

Instead, the backend exposes available capabilities and model policies.

---

# 4. Repository Structure

```text
drama-studio-web/
│
├── app/
│   ├── (marketing)/
│   │   ├── page.tsx
│   │   ├── about/
│   │   └── pricing/
│   │
│   ├── (studio)/
│   │   ├── layout.tsx
│   │   │
│   │   ├── projects/
│   │   │   ├── page.tsx
│   │   │   └── new/
│   │   │       └── page.tsx
│   │   │
│   │   └── projects/
│   │       └── [projectId]/
│   │           ├── page.tsx
│   │
│   │           ├── story/
│   │           │   ├── page.tsx
│   │           │   ├── bible/
│   │           │   ├── graph/
│   │           │   ├── timeline/
│   │           │   └── canon/
│   │           │
│   │           ├── characters/
│   │           │   ├── page.tsx
│   │           │   └── [characterId]/
│   │           │
│   │           ├── world/
│   │           │   ├── page.tsx
│   │           │   ├── locations/
│   │           │   └── props/
│   │           │
│   │           ├── seasons/
│   │           │   ├── page.tsx
│   │           │   └── [seasonId]/
│   │           │       └── episodes/
│   │           │           └── [episodeId]/
│   │           │               ├── page.tsx
│   │           │               ├── script/
│   │           │               ├── dialogue/
│   │           │               ├── scenes/
│   │           │               ├── storyboard/
│   │           │               ├── shots/
│   │           │               ├── assets/
│   │           │               ├── audio/
│   │           │               ├── assembly/
│   │           │               ├── continuity/
│   │           │               └── qa/
│   │           │
│   │           ├── production/
│   │           │   ├── page.tsx
│   │           │   ├── activity/
│   │           │   ├── jobs/
│   │           │   └── renders/
│   │           │
│   │           ├── assets/
│   │           │   ├── page.tsx
│   │           │   ├── characters/
│   │           │   ├── locations/
│   │           │   ├── images/
│   │           │   ├── video/
│   │           │   └── audio/
│   │           │
│   │           ├── continuity/
│   │           │   ├── page.tsx
│   │           │   ├── issues/
│   │           │   ├── canon/
│   │           │   ├── timeline/
│   │           │   └── character-state/
│   │           │
│   │           ├── automation/
│   │           │   ├── page.tsx
│   │           │   ├── director/
│   │           │   ├── agents/
│   │           │   ├── workflows/
│   │           │   └── decisions/
│   │           │
│   │           ├── publishing/
│   │           │   ├── page.tsx
│   │           │   └── channels/
│   │           │
│   │           └── settings/
│   │               ├── page.tsx
│   │               ├── production/
│   │               ├── ai-providers/
│   │               ├── models/
│   │               ├── budgets/
│   │               └── team/
│   │
│   └── api/
│
├── features/
│   ├── projects/
│   ├── story/
│   ├── bible/
│   ├── canon/
│   ├── characters/
│   ├── world/
│   ├── seasons/
│   ├── episodes/
│   ├── scripts/
│   ├── dialogue/
│   ├── storyboards/
│   ├── shots/
│   ├── assets/
│   ├── animation/
│   ├── audio/
│   ├── assembly/
│   ├── continuity/
│   ├── production/
│   ├── automation/
│   ├── agents/
│   ├── publishing/
│   └── analytics/
│
├── components/
│   ├── ui/
│   ├── layout/
│   ├── navigation/
│   ├── media/
│   ├── editor/
│   ├── timeline/
│   ├── graph/
│   ├── approvals/
│   └── feedback/
│
├── lib/
│   ├── api/
│   ├── auth/
│   ├── realtime/
│   ├── permissions/
│   ├── validation/
│   └── config/
│
├── stores/
│   ├── studio-store.ts
│   ├── workspace-store.ts
│   ├── timeline-store.ts
│   └── ui-store.ts
│
├── hooks/
│   ├── use-realtime.ts
│   ├── use-media-preview.ts
│   └── use-command-palette.ts
│
├── providers/
│   ├── query-provider.tsx
│   ├── auth-provider.tsx
│   └── realtime-provider.tsx
│
├── types/
│   ├── api.ts
│   ├── production.ts
│   ├── media.ts
│   └── ui.ts
│
├── config/
│   ├── navigation.ts
│   ├── permissions.ts
│   └── environments.ts
│
├── public/
│
└── package.json
```

---

# 5. Routing Architecture

`app/` owns routing.

`features/` owns application functionality.

Routes should remain thin.

Example:

```text
app/
└── (studio)/
    └── projects/
        └── [projectId]/
            └── seasons/
                └── [seasonId]/
                    └── episodes/
                        └── [episodeId]/
                            └── storyboard/
                                └── page.tsx
```

The route should primarily compose:

```text
Route
 ↓
Feature Page
 ↓
Feature Components
 ↓
Queries / Mutations
 ↓
API
```

Business logic should not accumulate inside `page.tsx`.

---

# 6. Feature Architecture

Each feature owns its own frontend behavior.

Example:

```text
features/storyboards/
├── components/
│   ├── storyboard-board.tsx
│   ├── storyboard-panel.tsx
│   ├── storyboard-toolbar.tsx
│   ├── shot-card.tsx
│   └── shot-detail.tsx
│
├── queries/
│   ├── use-storyboard.ts
│   └── use-storyboard-shots.ts
│
├── mutations/
│   ├── use-generate-storyboard.ts
│   ├── use-regenerate-shot.ts
│   └── use-approve-storyboard.ts
│
├── schemas/
│   └── storyboard.schema.ts
│
├── types/
│   └── storyboard.ts
│
└── index.ts
```

This keeps feature logic isolated and reusable.

---

# 7. Feature Boundaries

Frontend features correspond broadly to backend capabilities.

```text
Projects
Story
Bible
Canon
Characters
World
Seasons
Episodes
Scripts
Dialogue
Storyboards
Shots
Assets
Animation
Audio
Assembly
Continuity
Production
Automation
Agents
Publishing
Analytics
```

A feature should not reach directly into another feature's internals.

Use public exports:

```ts
import { StoryboardBoard } from "@/features/storyboards";
```

rather than:

```ts
import { StoryboardBoard } from "@/features/storyboards/components/internal/...";
```

---

# 8. Application State

The frontend uses two major state categories.

## Server State

Managed by TanStack Query.

Examples:

```text
Projects
Series
Seasons
Episodes
Scripts
Characters
Canon
Storyboards
Shots
Assets
Jobs
Agent Tasks
Continuity Issues
Production State
Publishing State
```

This state comes from the backend.

---

# 9. UI State

Zustand manages ephemeral client state.

Examples:

```text
Selected Scene
Selected Shot
Active Panel
Sidebar State
Timeline Position
Preview Mode
Filters
Command Palette
Keyboard Selection
Temporary Workspace State
```

Do not store production truth in Zustand.

Bad:

```ts
useProductionStore().episodeStatus;
```

if that status is authoritative backend state.

Better:

```ts
useEpisodeQuery();
```

---

# 10. API Layer

All backend communication should pass through a centralized API layer.

```text
lib/api/
├── client.ts
├── errors.ts
├── types.ts
└── endpoints/
```

The API client handles:

```text
Authentication
Headers
Request IDs
Error normalization
Serialization
Response parsing
```

Feature code should not duplicate `fetch()` logic everywhere.

---

# 11. API Client

Conceptually:

```ts
api.get(...)
api.post(...)
api.patch(...)
api.delete(...)
```

Feature-specific query hooks wrap these calls.

Example:

```text
useEpisode()
    ↓
api.get("/v1/projects/.../episodes/...")
```

---

# 12. Validation

Zod validates frontend input before requests are sent.

Examples:

```text
Project Creation
Episode Settings
Character Editing
Storyboard Editing
Shot Configuration
Production Settings
Publishing Metadata
```

Frontend validation improves UX.

Backend validation remains authoritative.

Never assume frontend validation provides security.

---

# 13. Forms

React Hook Form should manage complex forms.

Example:

```text
Character Editor
├── Identity
├── Appearance
├── Personality
├── Voice
├── Wardrobe
└── Visual Constraints
```

Use schemas:

```text
RHF
 ↓
Zod
 ↓
API
```

---

# 14. Studio Shell

The authenticated application uses a persistent studio shell.

Conceptually:

```text
┌─────────────────────────────────────────────────────────────┐
│ Logo     Project / Season / Episode       Search   Profile │
├───────────────┬─────────────────────────────────────────────┤
│               │                                             │
│ Navigation    │                                             │
│               │              Workspace                      │
│ Story         │                                             │
│ Characters    │                                             │
│ World         │                                             │
│ Seasons       │                                             │
│ Production    │                                             │
│ Assets        │                                             │
│ Continuity    │                                             │
│ Automation    │                                             │
│ Publishing    │                                             │
│               │                                             │
│               │                                             │
└───────────────┴─────────────────────────────────────────────┘
```

The shell should remain stable while the production workspace changes.

---

# 15. Project Workspace

A project is the primary workspace.

The project header should expose:

```text
Project Name
Production Status
Current Season
Current Episode
Autonomy Mode
Active Workflow
Budget Status
Continuity Status
```

The user should always understand:

> What is the system doing right now?

---

# 16. Series Development

The initial project experience should be conversational and structured.

User provides:

```text
Story Idea
```

The studio then guides the development process.

Possible stages:

```text
Idea
 ↓
Premise
 ↓
Series Bible
 ↓
Characters
 ↓
World
 ↓
Story Architecture
 ↓
Season Plan
```

The frontend presents AI-generated proposals as editable production artifacts.

---

# 17. Series Bible UI

The Bible interface should expose structured sections:

```text
Premise
Themes
Tone
Genre
World
Narrative Rules
Characters
Visual Direction
Dialogue Style
Story Constraints
Continuity Rules
```

Users should be able to:

```text
Edit
Approve
Version
Compare
Restore
Lock
```

---

# 18. Character Studio

The Character Studio is a dedicated workspace.

```text
Character
├── Identity
├── Appearance
├── Personality
├── Relationships
├── Story State
├── Knowledge
├── Wardrobe
├── Voice
├── Reference Images
└── Versions
```

The UI should visually connect:

```text
Character
 ↓
Reference Images
 ↓
Scenes
 ↓
Shots
 ↓
Generated Assets
```

---

# 19. Character Generation

Character generation is provider-independent.

The frontend asks the backend:

```text
What character-generation capability is available?
```

The UI does not assume a specific model.

Example:

```text
Character Generation

Reference
[ existing character state ]

Direction
[ structured creative instructions ]

Model
[ configured model ]

[ Generate ]
```

The backend decides how the request is executed.

---

# 20. Story Workspace

The story workspace should expose multiple representations of the same story.

```text
Story
├── Overview
├── Bible
├── Characters
├── Timeline
├── Story Graph
├── Canon
└── Plot Threads
```

The user should be able to move between narrative and structural views without duplicating data.

---

# 21. Story Graph

The Story Graph provides a visual representation of narrative dependencies.

Nodes may represent:

```text
Event
Reveal
Conflict
Decision
Relationship Change
Foreshadowing
Resolution
```

Edges represent:

```text
causes
reveals
depends_on
contradicts
resolves
foreshadows
```

The graph is a visualization of backend state, not an independent source of truth.

---

# 22. Timeline

The timeline should show story chronology.

Example:

```text
Episode 01
│
├── 08:00
├── 11:30
└── 22:15

Episode 02
│
├── 07:00
├── 13:00
└── 21:00
```

The interface can surface conflicts identified by the continuity system.

---

# 23. Episode Workspace

The episode is the primary production workspace.

Navigation:

```text
Overview
Script
Dialogue
Scenes
Storyboard
Shots
Assets
Audio
Assembly
Continuity
QA
```

A production status bar should indicate:

```text
Script       ✓
Dialogue     ✓
Storyboard   ✓
Assets       ✓
Animation    ●
Audio        ○
Assembly     ○
QA           ○
```

Statuses come from backend state.

---

# 24. Script Workspace

The Script workspace provides:

```text
Episode Structure
Scenes
Beats
Action
Dialogue
Character Presence
Location
Story Objectives
```

The script should support:

```text
Generate
Edit
Regenerate
Compare Versions
Approve
Reject
```

Every version remains associated with its generation history.

---

# 25. Dialogue Workspace

Dialogue is a separate production capability.

The UI should make clear that:

```text
Script
```

and:

```text
Dialogue
```

are related but independently generated artifacts.

The dialogue view can show:

```text
Character
Line
Emotion
Intent
Delivery
Voice
Scene Objective
```

---

# 26. Storyboard Workspace

Storyboard is a visual production board.

```text
Scene
 ├── Shot 01
 ├── Shot 02
 ├── Shot 03
 └── Shot 04
```

Each shot displays:

```text
Thumbnail
Shot Type
Camera
Character
Action
Location
Duration
Status
```

Actions:

```text
Generate
Regenerate
Approve
Replace
Reorder
Edit
```

---

# 27. Shot Workspace

A shot is the smallest major visual production unit.

The shot UI can show:

```text
Shot Specification
Character References
Location Reference
Storyboard
Motion
Camera
Duration
Aspect Ratio
Generation History
Continuity Constraints
```

Generated media should be shown alongside the specification that produced it.

---

# 28. Animation / Video Generation

The frontend does not directly call Wan or another video provider.

Flow:

```text
Shot
 ↓
Generate Video
 ↓
Backend
 ↓
Configured Video Capability
 ↓
External Provider
 ↓
Generation Job
 ↓
Realtime Status
 ↓
Asset
 ↓
Frontend Preview
```

The UI should display:

```text
Queued
Generating
Processing
Completed
Failed
```

Provider-specific details can be shown as metadata when appropriate, but the UI should remain capability-oriented.

---

# 29. Generation History

Every creative artifact should have a generation history.

Example:

```text
Shot 07

Generation 01
Rejected

Generation 02
Rejected

Generation 03
Approved
```

Each generation can expose:

```text
Model
Provider
Timestamp
Input
References
Parameters
Cost
Result
Status
```

The user can select which version becomes active.

---

# 30. Asset Library

The Asset Library provides a centralized view of project media.

Categories:

```text
Characters
Locations
Props
Images
Video
Audio
Storyboards
Renders
```

Filters:

```text
Episode
Scene
Character
Location
Asset Type
Status
Version
```

Assets should never appear as anonymous files.

The UI should preserve their production context.

---

# 31. Production Control Tower

The Production page provides a high-level operational view.

```text
Production
────────────────────────────────────

Episode 12

Current Phase
Animation

Active Tasks
8

Waiting
2

Completed
42

Continuity Issues
3

Approvals
1

Budget
...
```

It should answer:

> What is happening across the production?

---

# 32. Agent Activity

The user should be able to observe the AI production team.

Example:

```text
Lead Director
Planning Episode 12

Story Director
Reviewing unresolved plot threads

Storyboard Agent
Generating Scene 04 storyboard

Animation Agent
Waiting for video provider

Continuity Supervisor
Checking wardrobe continuity
```

Activity is streamed from the backend.

The frontend does not simulate agent activity.

---

# 33. Lead Director Console

The Lead Director console is the primary automation interface.

It should show:

```text
Current Objective
Current Phase
Active Tasks
Completed Tasks
Blocked Tasks
Pending Decisions
Recent Decisions
Production Budget
Quality Status
```

The user can:

```text
Pause
Resume
Stop
Approve
Reject
Override
Regenerate
Inspect
```

---

# 34. Autonomy Controls

Production mode should be visible at all times.

```text
Monitored
```

or:

```text
Autonomous
```

The UI should show why the system has paused.

Example:

```text
Production paused

Reason:
Human approval required before Episode 04
story direction is committed.

[ Review Decision ]
```

---

# 35. Approval Interface

Approval requests should be actionable.

Example:

```text
Lead Director Decision

Decision:
Use ending variant B for Episode 12.

Reason:
Variant B resolves the established plot thread
without contradicting Canon v18.

────────────────────────

[ Approve ]
[ Request Revision ]
[ Choose Alternative ]
[ Pause Production ]
```

The approval action is sent to the backend.

Temporal then resumes or redirects the workflow.

---

# 36. Continuity Center

Continuity is a first-class area.

Dashboard:

```text
Continuity
────────────────────────────

Blocking       1
Errors         2
Warnings       7
Resolved      18
```

Categories:

```text
Story
Timeline
Character
Knowledge
Relationship
Wardrobe
Props
Location
Visual
Dialogue
```

---

# 37. Continuity Issue View

Each issue should show:

```text
Issue
Expected State
Actual State
Evidence
Affected Episode
Affected Scene
Affected Character
Suggested Resolution
```

Actions:

```text
Resolve
Accept
Reject
Regenerate
Create Canon Change
```

The frontend should not decide whether something is canonically correct.

It presents the evidence and sends the user's chosen action to the backend.

---

# 38. Canon Interface

Canon should be treated as authoritative.

The UI should distinguish:

```text
Canonical
Proposed
Deprecated
Contradicted
```

Users should be able to inspect:

```text
Fact
Source Episode
Effective Time
Affected Characters
Character Knowledge
History
```

---

# 39. World Studio

World management includes:

```text
Locations
Rooms
Layouts
Props
Environmental References
Lighting Variants
Time-of-Day Variants
```

A location should have a persistent identity across episodes.

---

# 40. Assembly Workspace

Assembly is a timeline-oriented workspace.

```text
Video Track
────────────────────────
[Shot 01][Shot 02][Shot 03]

Dialogue
────────────────────────
[████][██████][████]

Music
────────────────────────
[██████████████████]

SFX
────────────────────────
    [██]       [██]

Subtitles
────────────────────────
[████████████████████]
```

The frontend manipulates a structured edit representation.

Actual rendering occurs on the backend.

---

# 41. Rendering

The frontend requests:

```text
Render Episode
```

The backend creates a rendering job.

The frontend receives:

```text
Queued
Processing
Encoding
Completed
Failed
```

The final render becomes a versioned asset.

---

# 42. Publishing Workspace

Publishing manages the transition from production to distribution.

```text
Episode
 ↓
Final QA
 ↓
Ready
 ↓
Publishing Configuration
 ↓
Schedule / Publish
```

The interface should support:

```text
Title
Description
Caption
Thumbnail
Tags
Channel
Schedule
Publication Status
```

---

# 43. Realtime Architecture

The frontend receives production events through SSE or WebSocket.

```text
Temporal / Backend
        ↓
Domain Event
        ↓
Realtime Gateway
        ↓
SSE / WebSocket
        ↓
Next.js
        ↓
TanStack Query
        ↓
UI
```

Examples:

```text
GenerationStarted
GenerationCompleted
AgentTaskCreated
AgentTaskCompleted
ContinuityIssueDetected
ApprovalRequested
RenderCompleted
```

Realtime events should trigger cache updates or invalidation.

They should not become an alternative database.

---

# 44. Reconnection

Realtime connections are not guaranteed.

When disconnected:

```text
Realtime
   ↓
Disconnected
   ↓
Reconnect
   ↓
Fetch Current State
   ↓
Resume Live Updates
```

The UI must remain correct even if realtime events were missed.

---

# 45. Optimistic Updates

Use optimistic updates selectively.

Good candidates:

```text
Panel preferences
Filters
Temporary UI changes
Non-critical local interactions
```

Avoid optimistic updates for irreversible production actions such as:

```text
Publishing
Canon changes
Production state transitions
Asset approval
AI generation
Stopping workflows
```

Those should wait for backend confirmation.

---

# 46. Loading States

The application should use intentional loading states.

Use:

```text
Skeletons
Progressive Loading
Streaming
Placeholder Previews
```

Avoid blank screens.

For long-running generation:

```text
Queued
 ↓
Submitted
 ↓
Provider Processing
 ↓
Downloading
 ↓
Validating
 ↓
Available
```

The UI should reflect actual backend state.

---

# 47. Error States

Errors should be contextual.

Bad:

```text
Something went wrong.
```

Better:

```text
Video generation failed.

The configured video provider rejected the
generation request.

Generation ID:
gen_123

[View Details]
[Retry]
[Change Configuration]
```

Do not expose sensitive provider details.

---

# 48. Permissions

The frontend should respect backend permissions.

Example permissions:

```text
project.read
project.update

story.read
story.edit
story.approve

character.read
character.edit
character.generate

production.read
production.control

agent.read
agent.control

continuity.read
continuity.resolve

publishing.read
publishing.manage
```

The backend remains authoritative.

Frontend permission checks are for UX.

---

# 49. Command Palette

A global command palette provides fast studio navigation.

Examples:

```text
Open Episode 12
Open Character Sarah
Generate Storyboard
Review Continuity Issues
Pause Production
View Agent Activity
Open Canon
Search Assets
```

Commands should respect user permissions and backend state.

---

# 50. Keyboard Navigation

A production application should support professional keyboard workflows.

Examples:

```text
Cmd/Ctrl + K
Command Palette

Space
Play / Pause Preview

Arrow Keys
Navigate Shots

Enter
Open Selected Artifact

Escape
Close Panel
```

Keyboard shortcuts should be centralized and configurable.

---

# 51. Media Preview System

Media previews should support:

```text
Images
Video
Audio
Storyboard Frames
Rendered Episodes
```

The preview system should understand production context.

For example:

```text
Video Preview

Episode 12
Scene 04
Shot 07
Generation 03
```

rather than showing a raw filename.

---

# 52. Responsive Design

The application should be desktop-first because production workflows require:

```text
Timelines
Story Graphs
Multi-panel Workspaces
Media Comparison
Dense Production Information
```

However, important operational views should remain usable on smaller screens:

```text
Production Status
Approvals
Agent Activity
Continuity Alerts
Notifications
Episode Status
```

Mobile should prioritize monitoring and decision-making rather than attempting to reproduce the complete desktop editing environment.

---

# 53. Design System

The frontend should use a centralized design system.

```text
components/ui/
├── button
├── input
├── select
├── dialog
├── drawer
├── tabs
├── tooltip
├── dropdown
├── command
├── table
├── badge
├── progress
└── ...
```

Product-specific components belong inside features.

Do not create dozens of near-identical UI components.

---

# 54. Iconography

Use **Lucide** consistently.

Do not mix unrelated icon libraries.

Avoid:

```text
Emoji as icons
Random SVG icons
Multiple icon libraries
Provider-specific icon conventions
```

Icons should communicate actions consistently across the studio.

---

# 55. Design Principles

The studio should avoid:

```text
Generic SaaS dashboard aesthetics
AI-slop visual language
Excessive gradients
Neon colors
Purple / indigo defaults
Generic blue interfaces
Glowing borders
Excessive glass effects
Unnecessary animations
Oversized decorative elements
Emoji-based interface icons
```

The interface should prioritize:

```text
Clarity
Hierarchy
Information density
Production context
Professional creative tooling
Consistency
Accessibility
Responsiveness
```

---

# 56. No Hardcoded Production Configuration

Avoid:

```ts
const MODEL = "wan";
const MAX_EPISODES = 50;
const DEFAULT_PROVIDER = "...";
```

when these values are backend-controlled configuration.

Instead:

```text
Frontend
 ↓
Backend configuration
 ↓
Capabilities / policies
 ↓
UI
```

The frontend should render what the backend says is available.

---

# 57. No Provider-Specific UI Coupling

Avoid designing the entire video-generation UI around one provider.

Bad:

```text
Wan Settings
Wan Steps
Wan Controls
```

Better:

```text
Video Generation
Model
Duration
Resolution
Aspect Ratio
Motion
References
Generation Settings
```

Provider-specific capabilities can appear dynamically when supported.

---

# 58. Data Fetching

TanStack Query should organize server state.

Example:

```text
queries/
├── use-project.ts
├── use-season.ts
├── use-episode.ts
├── use-character.ts
├── use-storyboard.ts
├── use-production.ts
└── use-continuity.ts
```

Query keys should be centralized and predictable.

Example:

```text
["projects", projectId]
["episodes", episodeId]
["storyboards", storyboardId]
```

---

# 59. Mutations

Production mutations should be explicit.

Examples:

```text
useApproveEpisode()
useApproveStoryboard()
useGenerateVideo()
useRegenerateShot()
usePauseProduction()
useResumeProduction()
useStopProduction()
useResolveContinuityIssue()
```

Mutations should invalidate or update the relevant query state after successful backend confirmation.

---

# 60. Frontend Types

Types should reflect API contracts.

```text
types/
├── api.ts
├── production.ts
├── media.ts
└── ui.ts
```

Avoid creating duplicate representations of backend entities without a reason.

Generated API types may be introduced later if the API specification supports it.

---

# 61. Server vs Client Components

Use Server Components where they provide meaningful benefits.

Use Client Components for:

```text
Interactive Editors
Media Players
Timelines
Story Graphs
Realtime Views
Forms
Command Palette
Drag and Drop
Interactive Production Controls
```

Do not mark entire route trees `"use client"` unnecessarily.

---

# 62. Caching

Caching should respect production state.

Static or slowly changing information can be cached.

Highly dynamic information should use:

```text
TanStack Query
Realtime Events
Backend Revalidation
```

Avoid stale production state in critical controls.

---

# 63. Production Safety

Dangerous actions should require explicit confirmation.

Examples:

```text
Stop Production
Delete Project
Delete Character
Replace Canon
Publish Episode
Cancel Generation
```

The confirmation should clearly explain the consequence.

---

# 64. Audit Visibility

Users should be able to inspect production history.

Example:

```text
Production History

16:42
Lead Director approved storyboard v4.

16:47
Video generation started.

16:52
Generation completed.

16:53
Continuity issue detected.

16:55
David approved character reference v3.
```

The frontend consumes the backend's audit/event history.

---

# 65. Accessibility

The studio should meet strong accessibility standards.

Requirements:

```text
Keyboard navigation
Focus management
Screen-reader labels
Semantic HTML
Visible focus states
Accessible dialogs
Accessible media controls
Color-independent status communication
Reduced motion support
```

Never rely solely on color to communicate:

```text
success
failure
warning
active
blocked
```

---

# 66. Performance

The application may display thousands of production artifacts.

Use:

```text
Pagination
Virtualized Lists
Lazy Loading
Thumbnail Generation
Progressive Media Loading
Query Caching
Code Splitting
Route-Level Loading
```

Do not load an entire project's assets into the browser.

---

# 67. Large Media

The browser should not proxy large video files through the Go API unnecessarily.

Prefer:

```text
Frontend
 ↓
Backend requests signed URL
 ↓
S3-compatible storage
 ↓
Browser
```

The backend remains responsible for authorization.

---

# 68. Security

The frontend must never receive:

```text
AI Provider API Keys
Database Credentials
Internal Service Credentials
Storage Master Credentials
Temporal Credentials
```

Authentication tokens and sessions should follow the selected authentication architecture.

All sensitive actions are authorized by the backend.

---

# 69. Environment Configuration

Frontend configuration should be explicit.

Example:

```text
NEXT_PUBLIC_API_URL
NEXT_PUBLIC_REALTIME_URL
```

Secrets must never use `NEXT_PUBLIC_`.

Do not hardcode production URLs.

---

# 70. Testing

## Unit Tests

Test:

```text
Utilities
Schemas
State Machines
Transformations
Feature Logic
```

## Component Tests

Test:

```text
Forms
Dialogs
Editors
Production Controls
Approval Interfaces
Media Components
```

## Integration Tests

Test:

```text
API Queries
Mutations
Realtime Updates
Authentication
Permissions
```

## End-to-End Tests

Critical workflows:

```text
Create Project
Create Series
Approve Bible
Create Season
Create Episode
Generate Script
Approve Script
Generate Storyboard
Generate Video
Resolve Continuity
Assemble Episode
Approve Episode
Publish
```

---

# 71. Error Monitoring

Frontend errors should be observable.

Capture:

```text
Runtime Errors
API Errors
Realtime Disconnects
Failed Mutations
Media Playback Errors
Unexpected State Transitions
```

Every error should include a request or trace identifier when available.

---

# 72. Realtime Event Handling

Realtime events should be mapped to application behavior.

Example:

```text
GenerationCompleted
        ↓
Update Asset Query
        ↓
Invalidate Shot Query
        ↓
Update Production Status
```

Avoid putting large event-processing systems directly inside UI components.

Use a dedicated realtime layer.

---

# 73. Frontend-to-Backend Contract

The frontend should think in terms of **production commands**, not implementation details.

Example:

```text
User:
Generate storyboard
```

Frontend sends:

```text
CreateStoryboardGeneration
```

The backend decides:

```text
Agent
Capability
Model
Provider
Workflow
Retry Policy
```

The frontend then observes the result.

---

# 74. Example: Video Generation

```text
User
 ↓
Clicks "Generate Video"
 ↓
Frontend validates shot configuration
 ↓
POST /video-generations
 ↓
Go Backend
 ↓
Production Application Service
 ↓
Temporal Workflow
 ↓
Video Capability
 ↓
Model Policy
 ↓
Provider Adapter
 ↓
External Video Provider
 ↓
Generation Complete
 ↓
Asset Created
 ↓
Realtime Event
 ↓
TanStack Query Update
 ↓
Video Preview
```

The frontend does not know how the provider executes the generation.

---

# 75. Example: Monitored Production

```text
Lead Director
 ↓
Creates decision
 ↓
Backend
 ↓
ApprovalRequested
 ↓
Realtime
 ↓
Frontend
 ↓
Approval Panel
 ↓
User Approves
 ↓
API Command
 ↓
Temporal Signal
 ↓
Workflow Continues
```

---

# 76. Example: Autonomous Production

```text
Lead Director
 ↓
Decision
 ↓
Decision Record
 ↓
Workflow Continues
 ↓
Frontend receives AgentDecisionMade
 ↓
Activity Feed updates
```

The frontend remains an observer and control surface.

It does not need to participate in every autonomous decision.

---

# 77. Workspace Philosophy

The application should feel like a **professional virtual production studio**.

Users should be able to move naturally between:

```text
Story
 ↓
Character
 ↓
Scene
 ↓
Storyboard
 ↓
Shot
 ↓
Generated Media
 ↓
Assembly
 ↓
Final Episode
```

Every artifact should retain its relationships.

A user looking at a video should be able to discover:

```text
Which shot?
Which scene?
Which episode?
Which character?
Which generation?
Which model?
Which provider?
Which references?
Which canon state?
```

---

# 78. Production Graph UI

The frontend can expose a production dependency graph:

```text
Series Bible
     ↓
Season
     ↓
Episode
     ↓
Scene
     ↓
Shot
 ┌───┼────┐
 ▼   ▼    ▼
Image Video Audio
 └───┼────┘
     ▼
   Assembly
     ↓
   Render
     ↓
  Publishing
```

This is a visualization of backend relationships.

---

# 79. Version Comparison

Important artifacts should support comparison.

Examples:

```text
Bible v3 vs v4
Script v7 vs v8
Storyboard v2 vs v3
Character v4 vs v5
Shot Generation 2 vs 3
```

The UI should show what changed without losing historical versions.

---

# 80. Frontend Architecture Boundary

The complete frontend architecture is:

```text
                     USER
                       │
                       ▼
              ┌─────────────────┐
              │   NEXT.JS APP    │
              └────────┬────────┘
                       │
             ┌─────────┴─────────┐
             ▼                   ▼
        Studio UI            Realtime
             │                   │
             ▼                   │
         Features                │
             │                   │
             ▼                   │
       TanStack Query            │
             │                   │
             ▼                   ▼
             └────── API / Events
                       │
                       ▼
                  GO BACKEND
                       │
              ┌────────┴────────┐
              ▼                 ▼
          Production          Agents
              │                 │
              └────────┬────────┘
                       ▼
                 AI Capabilities
                       │
                       ▼
              External Providers
```

---

# 81. What the Frontend Owns

The frontend owns:

```text
Presentation
Interaction
Navigation
Workspace State
Forms
Validation UX
Media Preview
Realtime Visualization
Command Interfaces
Approval Interfaces
Production Monitoring
```

---

# 82. What the Frontend Does Not Own

The frontend does not own:

```text
Canon
Production Truth
Workflow Execution
AI Model Selection
Provider Credentials
Provider APIs
Agent Decisions
Generation State
Budget Authority
Continuity Authority
Publishing Authority
```

Those belong to the backend.

---

# 83. Final Architecture Principle

The frontend is not a CRUD interface wrapped around an AI API.

It is the **control room of an autonomous digital production studio**.

The user can enter a story and watch it progress through:

```text
IDEA
 ↓
BIBLE
 ↓
STORY
 ↓
CHARACTERS
 ↓
EPISODES
 ↓
SCRIPT
 ↓
DIALOGUE
 ↓
STORYBOARD
 ↓
SHOTS
 ↓
VISUAL GENERATION
 ↓
VIDEO / ANIMATION
 ↓
AUDIO
 ↓
ASSEMBLY
 ↓
CONTINUITY
 ↓
QA
 ↓
FINAL RENDER
 ↓
PUBLISHING
```

The frontend makes that entire process **visible, understandable, controllable, and editable**.

The backend remains responsible for making it **durable, consistent, autonomous, and correct**.

> **The frontend is the studio interface.
> The backend is the production system.
> The agents are the production team.
> The external AI providers are the creative engines.
> Canon and production state are the source of truth.**
