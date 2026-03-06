// Evaluation types
export type Status = 'pending' | 'running' | 'completed' | 'failed' | 'canceled'

export interface Evaluation {
  id: string
  created_at: string
  updated_at: string
  status: Status
  config: EvalConfig
  total_cases: number
  completed_cases: number
  error?: string
}

export interface EvalConfig {
  models: string[]
  datasets: string[]
  sample_size: number
  max_workers: number
  ephemeral: boolean
}

export interface Metrics {
  accuracy: number
  f1: number
  bleu: number
  rouge_l: number
  avg_latency: number
  avg_tokens_per_second: number
}

export interface ModelResult {
  model_name: string
  metrics: Metrics
  error_count: number
}

// Dataset types
export interface Dataset {
  name: string
  source: string
  task_type: string
  total_cases: number
  description: string
}

// Model types
export interface ModelConfig {
  name: string
  endpoint: string
  timeout: number
}

// SSE Event types
export type EventType = 'progress' | 'model_complete' | 'evaluation_complete' | 'error'

export interface SSEEvent {
  type: EventType
  data: Record<string, unknown>
}

// API Response types
export interface HealthResponse {
  status: string
  version: string
}

export interface CreateEvaluationRequest {
  models: string[]
  datasets: string[]
  config: {
    sample_size?: number
    max_workers?: number
    ephemeral?: boolean
  }
}
