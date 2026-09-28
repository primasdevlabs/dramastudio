export type ProductionMode = "monitored" | "autonomous";

export type ProjectStatus =
  | "PROJECT_CREATED"
  | "DEVELOPING_SERIES"
  | "SERIES_APPROVED"
  | "PLANNING_SEASON"
  | "PRODUCING_EPISODE"
  | "VALIDATING_EPISODE"
  | "POST_PRODUCTION"
  | "READY"
  | "PUBLISHED";

export interface Budget {
  max_cost_per_generation: number;
  total_budget: number;
  current_spent: number;
  currency: string;
}

export interface Project {
  id: string;
  name: string;
  description: string;
  genre: string;
  language: string;
  mode: ProductionMode;
  budget: Budget;
  status: ProjectStatus;
  created_at: string;
  updated_at: string;
}

export interface SeriesBible {
  id: string;
  project_id: string;
  version: number;
  premise: string;
  genre: string;
  themes: string[];
  tone: string;
  world_rules: string[];
  narrative_rules: string[];
  visual_direction: string;
  dialogue_style: string;
  story_constraints: string[];
  created_at: string;
}

export interface CharacterRelationship {
  target_id: string;
  type: string;
  description: string;
}

export interface Character {
  id: string;
  project_id: string;
  version: number;
  name: string;
  role: string;
  bio: string;
  is_locked: boolean;
  appearance?: {
    hair?: string;
    eyes?: string;
    skin?: string;
    features?: string;
  };
  wardrobe?: {
    outfit?: string;
    colors?: string;
  };
  relationships?: CharacterRelationship[];
}

export interface LocationVariant {
  id: string;
  name: string;
  attributes: Record<string, string>;
  asset_url: string;
}

export interface Location {
  id: string;
  project_id: string;
  name: string;
  description: string;
  kind: string;
  parent_id?: string;
  visual_ref?: string;
  variants: LocationVariant[];
}

export interface Season {
  id: string;
  series_id: string;
  number: number;
  title: string;
  summary: string;
  episode_ids: string[];
}

export type EpisodeStatus =
  | "PLANNED"
  | "SCRIPTED"
  | "PRODUCING"
  | "VALIDATING"
  | "COMPLETED"
  | "PUBLISHED";

export interface Episode {
  id: string;
  season_id: string;
  arc_id: string;
  number: number;
  title: string;
  summary: string;
  script?: string;
  status: EpisodeStatus;
  scene_ids: string[];
}

export interface Scene {
  id: string;
  episode_id: string;
  number: number;
  title: string;
  description: string;
  location_id: string;
  time_of_day: string;
  character_ids: string[];
  beat_ids: string[];
}

export interface StoryFact {
  id: string;
  project_id: string;
  entity_id?: string;
  subject: string;
  predicate: string;
  object: string;
  type?: string;
  introduced_episode?: string;
  effective_from?: string;
  effective_until?: string;
  source?: string;
  confidence: number;
  status: "canonical" | "disputed" | "retconned";
  version: number;
  created_at: string;
}

export type MediaType =
  | "image"
  | "video"
  | "voice"
  | "music"
  | "sfx"
  | "subtitle"
  | "storyboard"
  | "render"
  | "character_ref"
  | "location_ref"
  | "prop";

export interface Asset {
  id: string;
  project_id: string;
  character_id?: string;
  location_id?: string;
  episode_id?: string;
  scene_id?: string;
  shot_id?: string;
  type: MediaType;
  provider: string;
  model: string;
  prompt: string;
  url: string;
  version: number;
  status: "PENDING" | "APPROVED" | "REJECTED" | "ARCHIVED";
  cost: number;
  created_at: string;
}

export interface GenerationJob {
  id: string;
  project_id: string;
  asset_id: string;
  capability: string;
  provider: string;
  model: string;
  provider_job_id?: string;
  input: string;
  output_url?: string;
  status: "PENDING" | "RUNNING" | "SUCCEEDED" | "FAILED" | "CANCELLED";
  attempt: number;
  cost: number;
  error?: string;
  started_at: string;
  completed_at?: string;
}

export type ShotStatus =
  | "planned"
  | "generating"
  | "generated"
  | "approved"
  | "rejected";

export interface Shot {
  id: string;
  project_id: string;
  episode_id: string;
  scene_id: string;
  seq: number;
  description: string;
  camera?: Record<string, unknown>;
  characters?: string[];
  location_id?: string;
  duration_sec: number;
  status: ShotStatus;
  approved_asset?: string;
}

export type ContinuitySeverity = "INFO" | "WARNING" | "ERROR" | "BLOCKING";

export type IssueStatus = "open" | "resolved" | "wontfix";

export interface ContinuityIssue {
  id: string;
  check_id?: string;
  project_id: string;
  episode_id?: string;
  scene_id?: string;
  category: string;
  severity: ContinuitySeverity;
  entity: string;
  expected_state: string;
  actual_state: string;
  cause: string;
  evidence: string;
  resolution?: string;
  status: IssueStatus;
  created_at: string;
  resolved_at?: string;
}

export interface ContinuityCheck {
  id: string;
  project_id: string;
  episode_id: string;
  check_type: string;
  target_id?: string;
  status: string;
  issue_count: number;
  created_at: string;
  finished_at?: string;
}

export interface LeadDirectorDecision {
  id: string;
  project_id: string;
  episode_id: string;
  decision: string;
  reason: string;
  decision_maker: "USER" | "LEAD_DIRECTOR" | "SPECIALIZED_AGENT" | "SYSTEM";
  mode: string;
  timestamp: string;
}

export type ApprovalDecision =
  | "APPROVE"
  | "REJECT"
  | "REQUEST_REVISION"
  | "REGENERATE"
  | "PAUSE"
  | "STOP"
  | "OVERRIDE";

export interface ApprovalRequest {
  id: string;
  project_id: string;
  episode_id: string;
  stage: string;
  target_id: string;
  decision: ApprovalDecision;
  notes: string;
  decided_by: string;
  submitted_at: string;
  decided_at?: string;
}

export type RunStatus =
  | "pending"
  | "running"
  | "paused"
  | "awaiting_approval"
  | "completed"
  | "failed"
  | "stopped";

export interface ProductionRun {
  id: string;
  production_id: string;
  project_id: string;
  episode_id: string;
  stage: string;
  status: RunStatus;
  bible_version: number;
  workflow_id?: string;
  jobs?: ProductionJob[];
  created_at: string;
  completed_at?: string;
}

export type JobStatus =
  | "queued"
  | "running"
  | "succeeded"
  | "failed"
  | "blocked"
  | "skipped";

export interface ProductionJob {
  id: string;
  production_id?: string;
  project_id: string;
  run_id?: string;
  episode_id?: string;
  scene_id?: string;
  shot_id?: string;
  kind: string;
  status: JobStatus;
  attempt: number;
  result_url?: string;
  error?: string;
  created_at: string;
  completed_at?: string;
}

export type TaskStatus =
  | "pending"
  | "running"
  | "succeeded"
  | "failed"
  | "escalated"
  | "cancelled";

export interface AgentTask {
  id: string;
  project_id: string;
  agent_id: string;
  objective: string;
  status: TaskStatus;
  iteration: number;
  error?: string;
  created_at: string;
  completed_at?: string;
}

export interface RenderTask {
  id: string;
  timeline_id: string;
  project_id: string;
  episode_id: string;
  format: string;
  resolution: string;
  status: "queued" | "rendering" | "succeeded" | "failed";
  output_url?: string;
  error?: string;
  created_at: string;
  finished_at?: string;
}

export interface TrackItem {
  id: string;
  shot_id: string;
  start_time: number;
  duration: number;
  asset_url: string;
}

export interface EpisodeTimeline {
  id: string;
  project_id: string;
  episode_id: string;
  version: number;
  status: "draft" | "approved" | "rendered";
  video_tracks: TrackItem[];
  audio_tracks: TrackItem[];
  created_at: string;
}

export interface Channel {
  id: string;
  project_id: string;
  platform: "youtube" | "tiktok" | "instagram";
  name: string;
  account_ref?: string;
  enabled: boolean;
  created_at: string;
}

export interface PublishMetadata {
  title: string;
  caption: string;
  description?: string;
  tags: string[];
}

export interface Publication {
  id: string;
  project_id: string;
  episode_id: string;
  channel_id: string;
  metadata: PublishMetadata;
  video_url: string;
  scheduled_at?: string;
  status: "draft" | "scheduled" | "publishing" | "published" | "failed";
  external_id?: string;
  published_at?: string;
  created_at: string;
}

export interface MetricSummary {
  metric: string;
  total: number;
  count: number;
}

// Every list endpoint returns this envelope.
export interface ListResponse<T> {
  items: T[] | null;
}
