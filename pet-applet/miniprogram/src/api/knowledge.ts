import { request } from './request'

export interface KnowledgeResult {
  title: string
  content: string
  tags: string[]
  score: number
}

export interface LLMAnswer {
  answer: string
  sources: string[]
}

export interface SearchResponse {
  results: KnowledgeResult[]
  answer: LLMAnswer | null
}

export function searchKnowledge(params: { q?: string; breed?: string }) {
  const query = new URLSearchParams()
  if (params.q) query.set('q', params.q)
  if (params.breed) query.set('breed', params.breed)
  return request<SearchResponse>(`/knowledge/search?${query.toString()}`)
}
