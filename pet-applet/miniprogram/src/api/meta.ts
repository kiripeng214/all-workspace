import { request } from './request'

export interface MetaData {
  petEmojis: string[]
  breedOptions: Record<string, string[]>
}

export function getBreeds() {
  return request<MetaData>('/meta/breeds')
}
