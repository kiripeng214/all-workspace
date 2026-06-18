import { request } from './request'

export interface FeedingRecord {
  id: string
  petId: string
  scheduleId: string
  time: string
  foodType: string
  amount: string
  notes: string
  createdAt: number
}

export interface CreateRecordParams {
  scheduleId?: string
  time?: string
  foodType?: string
  amount?: string
  notes?: string
}

export function getRecords(petId: string) {
  return request<FeedingRecord[]>(`/pets/records/${petId}`)
}

export function getTodayRecords(petId: string) {
  return request<FeedingRecord[]>(`/pets/records/today/${petId}`)
}

export function createRecord(petId: string, data: CreateRecordParams) {
  return request<FeedingRecord>({ url: `/pets/records/${petId}`, method: 'POST', data })
}

export function deleteRecord(id: string) {
  return request<{ message: string }>({ url: `/records/${id}`, method: 'DELETE' })
}
