import { describe, it, expect, vi, beforeEach } from 'vitest'
import { request } from '@/api/request'

const mockRequest = vi.fn()

beforeEach(() => {
  vi.clearAllMocks()
  ;(globalThis as any).uni.request = mockRequest
})

describe('request wrapper', () => {
  it('GET request with string param', async () => {
    mockRequest.mockImplementation((options: any) => {
      options.success({ data: { id: '1' }, statusCode: 200 })
    })

    const result = await request('/pets')
    expect(result).toEqual({ id: '1' })
    expect(mockRequest.mock.calls[0][0].url).toContain('/pets')
    expect(mockRequest.mock.calls[0][0].method).toBe('GET')
  })

  it('POST request with data', async () => {
    mockRequest.mockImplementation((options: any) => {
      options.success({ data: { id: 'new' }, statusCode: 201 })
    })

    const result = await request({ url: '/pets', method: 'POST', data: { name: '旺财' } })
    expect(result).toEqual({ id: 'new' })
    expect(mockRequest.mock.calls[0][0].method).toBe('POST')
    expect(mockRequest.mock.calls[0][0].data.name).toBe('旺财')
  })

  it('throws on 4xx error', async () => {
    mockRequest.mockImplementation((options: any) => {
      options.success({ data: { error: 'not found' }, statusCode: 404 })
    })

    await expect(request('/pets/404')).rejects.toBeDefined()
  })

  it('throws on network error', async () => {
    mockRequest.mockImplementation((options: any) => {
      options.fail({ errMsg: 'request:fail' })
    })

    await expect(request('/pets')).rejects.toBeDefined()
  })
})
