import { describe, it, expect, vi, beforeEach } from 'vitest'
import { getPets, getPet, createPet, updatePet, deletePet } from '@/api/pets'

const mockRequest = vi.fn()

beforeEach(() => {
  vi.clearAllMocks()
  // Make uni.request point to our mock
  ;(globalThis as any).uni.request = mockRequest
  // Default: make uni.request succeed
  mockRequest.mockImplementation((options: any) => {
    if (options.success) {
      options.success({ data: {}, statusCode: 200 })
    }
  })
})

describe('pets API', () => {
  it('getPets calls GET /pets', () => {
    getPets()
    expect(mockRequest).toHaveBeenCalledTimes(1)
    const args = mockRequest.mock.calls[0][0]
    expect(args.url).toContain('/pets')
    expect(args.method).toBe('GET')
  })

  it('getPet calls GET /pets/:id', () => {
    getPet('abc123')
    expect(mockRequest).toHaveBeenCalledTimes(1)
    const args = mockRequest.mock.calls[0][0]
    expect(args.url).toContain('/pets/abc123')
    expect(args.method).toBe('GET')
  })

  it('createPet sends POST with data', () => {
    createPet({ name: '旺财', avatar: '🐶' })
    expect(mockRequest).toHaveBeenCalledTimes(1)
    const args = mockRequest.mock.calls[0][0]
    expect(args.url).toContain('/pets')
    expect(args.method).toBe('POST')
    expect(args.data.name).toBe('旺财')
    expect(args.data.avatar).toBe('🐶')
  })

  it('updatePet sends PUT with only changed fields', () => {
    updatePet('abc123', { name: '新名字' })
    expect(mockRequest).toHaveBeenCalledTimes(1)
    const args = mockRequest.mock.calls[0][0]
    expect(args.url).toContain('/pets/abc123')
    expect(args.method).toBe('PUT')
    // Should NOT send avatar/breed/birthday if not provided
    expect(args.data.name).toBe('新名字')
    expect(args.data.avatar).toBeUndefined()
    expect(args.data.breed).toBeUndefined()
  })

  it('deletePet sends DELETE', () => {
    deletePet('abc123')
    expect(mockRequest).toHaveBeenCalledTimes(1)
    const args = mockRequest.mock.calls[0][0]
    expect(args.url).toContain('/pets/abc123')
    expect(args.method).toBe('DELETE')
  })
})
