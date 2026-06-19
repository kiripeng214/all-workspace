import { describe, it, expect, vi, beforeEach } from 'vitest'
import { getTodayRecords, createRecord, deleteRecord } from '@/api/records'

const mockRequest = vi.fn()

beforeEach(() => {
  vi.clearAllMocks()
  ;(globalThis as any).uni.request = mockRequest
  mockRequest.mockImplementation((options: any) => {
    if (options.success) options.success({ data: {}, statusCode: 200 })
  })
})

describe('records API', () => {
  it('getTodayRecords calls GET with petId', () => {
    getTodayRecords('pet_001')
    const args = mockRequest.mock.calls[0][0]
    expect(args.url).toContain('/pets/records/today/pet_001')
    expect(args.method).toBe('GET')
  })

  it('createRecord sends POST with feed data', () => {
    createRecord('pet_001', { time: '08:00', foodType: '猫粮', amount: '一份' })
    const args = mockRequest.mock.calls[0][0]
    expect(args.url).toContain('/pets/records/pet_001')
    expect(args.method).toBe('POST')
    expect(args.data.foodType).toBe('猫粮')
  })

  it('deleteRecord sends DELETE', () => {
    deleteRecord('rec_001')
    const args = mockRequest.mock.calls[0][0]
    expect(args.url).toContain('/records/rec_001')
    expect(args.method).toBe('DELETE')
  })
})
