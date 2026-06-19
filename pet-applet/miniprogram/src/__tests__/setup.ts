import { vi } from 'vitest'

// Mock uni-app global APIs used by our code
const uniMock: Record<string, any> = {
  request: vi.fn(),
  navigateTo: vi.fn(),
  navigateBack: vi.fn(),
  showToast: vi.fn(),
  showModal: vi.fn((options: any) => {
    if (options.success) options.success({ confirm: true, cancel: false })
  }),
  onLoad: vi.fn(),
  onShow: vi.fn(),
  getSystemInfoSync: vi.fn(() => ({
    windowWidth: 375,
    windowHeight: 812,
    pixelRatio: 2,
  })),
}

// Attach to global scope
;(globalThis as any).uni = uniMock
;(globalThis as any).getApp = vi.fn(() => ({}))
