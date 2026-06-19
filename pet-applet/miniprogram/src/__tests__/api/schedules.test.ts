import { describe, it, expect, vi, beforeEach } from 'vitest'
import { getSchedules, createSchedule, updateSchedule, deleteSchedule } from '@/api/schedules'

const mockRequest = vi.fn()

beforeEach(() => {
  vi.clearAllMocks()
  ;(globalThis as any).uni.request = mockRequest
  mockRequest.mockImplementation((options: any) => {
    if (options.success) options.success({ data: {}, statusCode: 200 })
  })
})

describe('schedules API', () => {
  it('getSchedules calls GET with petId', () => {
    getSchedules('pet_001')
    const args = mockRequest.mock.calls[0][0]
    expect(args.url).toContain('/pets/schedules/pet_001')
    expect(args.method).toBe('GET')
  })

  it('createSchedule sends POST with schedule data', () => {
    createSchedule('pet_001', { time: '08:00', foodType: '猫粮', amount: '一份' })
    const args = mockRequest.mock.calls[0][0]
    expect(args.url).toContain('/pets/schedules/pet_001')
    expect(args.method).toBe('POST')
    expect(args.data.time).toBe('08:00')
    expect(args.data.foodType).toBe('猫粮')
  })

  it('updateSchedule sends PUT with id and data', () => {
    updateSchedule('sched_001', { time: '09:00' })
    const args = mockRequest.mock.calls[0][0]
    expect(args.url).toContain('/schedules/sched_001')
    expect(args.method).toBe('PUT')
    expect(args.data.time).toBe('09:00')
  })

  it('deleteSchedule sends DELETE', () => {
    deleteSchedule('sched_001')
    const args = mockRequest.mock.calls[0][0]
    expect(args.url).toContain('/schedules/sched_001')
    expect(args.method).toBe('DELETE')
  })
})
