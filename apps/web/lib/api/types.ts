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
}

export interface LocationVariant {
  id: string;
  time_of_day: string;
  asset_url: string;
}

export interface Location {
  id: string;
  project_id: string;
  name: string;
  description: string;
  type: string;
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

export interface Episode {
  id: string;
  season_id: string;
  arc_id: string;
  number: number;
  title: string;
  summary: string;
  script?: string;
  status: "PLANNED" | "SCRIPTED" | "PRODUCING" | "VALIDATING" | "COMPLETED";
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
}

export interface StoryFact {
  id: string;
  subject: string;
  predicate: string;
  object: string;
  introduced: string;
  valid_from: string;
  status: "canonical" | "disputed" | "retconned";
}

export interface Asset {
  id: string;
  project_id: string;
  character_id?: string;
  scene_id?: string;
  shot_id?: string;
  type: "image" | "video" | "voice" | "music" | "sfx" | "render";
  provider: string;
  model: string;
  prompt: string;
  url: string;
  version: number;
  status: "PENDING" | "APPROVED" | "REJECTED" | "ARCHIVED";
  cost: number;
  created_at: string;
}

export type ContinuitySeverity = "INFO" | "WARNING" | "ERROR" | "BLOCKING";

export interface ContinuityIssue {
  id: string;
  project_id: string;
  episode_id: string;
  scene_id: string;
  category: string;
  severity: ContinuitySeverity;
  entity: string;
  expected_state: string;
  actual_state: string;
  cause: string;
  evidence: string;
  resolution?: string;
  is_resolved: boolean;
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

export interface ApprovalRequest {
  id: string;
  project_id: string;
  episode_id: string;
  stage: string;
  target_id: string;
  decision: "APPROVE" | "REJECT" | "REQUEST_REVISION" | "REGENERATE" | "PAUSE" | "STOP" | "OVERRIDE";
  notes: string;
  decided_by: string;
  submitted_at: string;
  decided_at: string;
}

export interface RenderTask {
  id: string;
  episode_id: string;
  format: string;
  resolution: string;
  status: "PENDING" | "RENDERING" | "COMPLETED" | "FAILED";
  output_url: string;
  created_at: string;
}

export interface Publication {
  id: string;
  episode_id: string;
  channel_id: string;
  metadata: {
    title: string;
    caption: string;
    tags: string[];
  };
  published_at: string;
}
