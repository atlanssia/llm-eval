import type {
  Evaluation,
  Dataset,
  ModelConfig,
  CreateEvaluationRequest,
  HealthResponse,
  SSEEvent,
} from './types'

const API_BASE = '/api'

class APIClient {
  private baseURL: string

  constructor(baseURL: string = API_BASE) {
    this.baseURL = baseURL
  }

  async get<T>(path: string): Promise<T> {
    const response = await fetch(`${this.baseURL}${path}`)
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }
    return response.json()
  }

  async post<T>(path: string, body: unknown): Promise<T> {
    const response = await fetch(`${this.baseURL}${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }
    return response.json()
  }
}

const api = new APIClient()

// Health check
export async function getHealth(): Promise<HealthResponse> {
  return api.get<HealthResponse>('/health')
}

// Datasets
export async function getDatasets(): Promise<Dataset[]> {
  return api.get<Dataset[]>('/datasets')
}

// Models
export async function getModels(): Promise<ModelConfig[]> {
  return api.get<ModelConfig[]>('/models')
}

// Evaluations
export async function getEvaluations(): Promise<Evaluation[]> {
  return api.get<Evaluation[]>('/evaluations')
}

export async function getEvaluation(id: string): Promise<Evaluation> {
  return api.get<Evaluation>(`/evaluations/${id}`)
}

export async function createEvaluation(request: CreateEvaluationRequest): Promise<Evaluation> {
  return api.post<Evaluation>('/evaluations', request)
}

// SSE streaming
export function streamEvaluation(id: string, onEvent: (event: SSEEvent) => void): () => void {
  const eventSource = new EventSource(`${API_BASE}/evaluations/${id}/stream`)

  eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data) as SSEEvent
    onEvent(data)
  }

  eventSource.onerror = (error) => {
    console.error('SSE error:', error)
    eventSource.close()
  }

  // Return cleanup function
  return () => {
    eventSource.close()
  }
}

export default api
