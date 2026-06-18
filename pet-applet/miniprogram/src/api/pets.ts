import { request } from './request'

export interface Pet {
  id: string
  avatar: string
  name: string
  breed: string
  birthday: string
  weight: string
  notes: string
  createdAt: number
}

export interface CreatePetParams {
  avatar?: string
  name: string
  breed?: string
  birthday?: string
  weight?: string
  notes?: string
}

export type UpdatePetParams = Partial<CreatePetParams>

export function getPets() {
  return request<Pet[]>('/pets')
}

export function getPet(id: string) {
  return request<Pet>(`/pets/${id}`)
}

export function createPet(data: CreatePetParams) {
  return request<Pet>({ url: '/pets', method: 'POST', data })
}

export function updatePet(id: string, data: UpdatePetParams) {
  return request<Pet>({ url: `/pets/${id}`, method: 'PUT', data })
}

export function deletePet(id: string) {
  return request<{ message: string }>({ url: `/pets/${id}`, method: 'DELETE' })
}
