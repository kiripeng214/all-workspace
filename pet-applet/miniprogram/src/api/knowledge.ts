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
  const parts: string[] = []
  if (params.q) parts.push('q=' + encodeURIComponent(params.q))
  if (params.breed) parts.push('breed=' + encodeURIComponent(params.breed))
  return request<SearchResponse>(`/knowledge/search?${parts.join('&')}`)
}
