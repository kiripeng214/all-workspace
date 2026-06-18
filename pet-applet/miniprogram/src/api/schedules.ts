import { request } from './request'

export interface FeedingSchedule {
  id: string
  petId: string
  time: string
  foodType: string
  amount: string
}

export interface CreateScheduleParams {
  time: string
  foodType?: string
  amount?: string
}

export type UpdateScheduleParams = Partial<CreateScheduleParams>

export function getSchedules(petId: string) {
  return request<FeedingSchedule[]>(`/pets/schedules/${petId}`)
}

export function createSchedule(petId: string, data: CreateScheduleParams) {
  return request<FeedingSchedule>({ url: `/pets/schedules/${petId}`, method: 'POST', data })
}

export function updateSchedule(id: string, data: UpdateScheduleParams) {
  return request<FeedingSchedule>({ url: `/schedules/${id}`, method: 'PUT', data })
}

export function deleteSchedule(id: string) {
  return request<{ message: string }>({ url: `/schedules/${id}`, method: 'DELETE' })
}
