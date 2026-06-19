import { describe, it, expect, vi, beforeEach } from 'vitest'
import { getBreeds } from '@/api/meta'

const mockRequest = vi.fn()

beforeEach(() => {
  vi.clearAllMocks()
  ;(globalThis as any).uni.request = mockRequest
  mockRequest.mockImplementation((options: any) => {
    if (options.success) options.success({ data: {}, statusCode: 200 })
  })
})

describe('meta API', () => {
  it('getBreeds calls GET /meta/breeds', () => {
    getBreeds()
    const args = mockRequest.mock.calls[0][0]
    expect(args.url).toContain('/meta/breeds')
    expect(args.method).toBe('GET')
  })
})
