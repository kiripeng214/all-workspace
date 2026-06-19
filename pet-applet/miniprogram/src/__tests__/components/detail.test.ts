import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'

// Mock uni-app lifecycle hooks
vi.mock('@dcloudio/uni-app', () => ({
  onLoad: vi.fn((cb: any) => {
    setTimeout(() => cb({ id: 'pet_001' }), 0)
  }),
  onShow: vi.fn(),
}))

const mockPet = {
  id: 'pet_001',
  avatar: '🐶',
  name: '旺财',
  breed: '金毛',
  birthday: '2020-03-15',
  weight: '15kg',
  notes: '可爱',
  createdAt: 1700000000000,
}

const mockUpdatePet = vi.fn()
vi.mock('@/api', () => ({
  getPet: vi.fn(() => Promise.resolve(mockPet)),
  getTodayRecords: vi.fn(() => Promise.resolve([])),
  getSchedules: vi.fn(() => Promise.resolve([])),
  createRecord: vi.fn(() => Promise.resolve({})),
  updatePet: (...args: any[]) => {
    mockUpdatePet(...args)
    return Promise.resolve({ ...mockPet, ...args[1] })
  },
  deletePet: vi.fn(() => Promise.resolve({ message: '已删除' })),
  getBreeds: vi.fn(() => Promise.resolve({ petEmojis: [], breedOptions: {} })),
}))

let DetailPage: any

beforeEach(async () => {
  vi.clearAllMocks()
  ;(globalThis as any).uni.request = vi.fn((options: any) => {
    if (options.success) options.success({ data: mockPet, statusCode: 200 })
  })
  // 延迟 import 以便 mock 先生效
  if (!DetailPage) {
    DetailPage = (await import('@/pages/pets/detail.vue')).default
  }
})

async function createWrapper() {
  return mount(DetailPage, {
    global: {
      stubs: {
        'uni-popup': {
          template: '<view class="uni-popup-stub"><slot /></view>',
          props: ['show', 'modelValue'],
        },
      },
    },
  })
}

async function waitForRender(wrapper: any) {
  await wrapper.vm.$nextTick()
  await new Promise(r => setTimeout(r, 100))
  await wrapper.vm.$nextTick()
}

describe('改名弹窗交互', () => {
  it('点击宠物姓名弹出改名弹窗', async () => {
    const wrapper = await createWrapper()
    await waitForRender(wrapper)

    const nameEl = wrapper.find('.name')
    expect(nameEl.text()).toContain('旺财')

    await nameEl.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.overlay').exists()).toBe(true)
    expect(wrapper.find('.rename-title').text()).toBe('修改姓名')

    const input = wrapper.find('.rename-input').element as HTMLInputElement
    expect(input.value).toBe('旺财')
  })

  it('取消按钮关闭改名弹窗', async () => {
    const wrapper = await createWrapper()
    await waitForRender(wrapper)

    await wrapper.find('.name').trigger('click')
    await wrapper.vm.$nextTick()
    expect(wrapper.find('.overlay').exists()).toBe(true)

    await wrapper.find('.rename-btn.cancel').trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.overlay').exists()).toBe(false)
  })

  it('确认改名调用 updatePet 并刷新', async () => {
    const wrapper = await createWrapper()
    await waitForRender(wrapper)

    await wrapper.find('.name').trigger('click')
    await wrapper.vm.$nextTick()

    const input = wrapper.find('.rename-input')
    await input.setValue('小强')

    await wrapper.find('.rename-btn.confirm').trigger('click')
    await new Promise(r => setTimeout(r, 50))
    await wrapper.vm.$nextTick()

    expect(mockUpdatePet).toHaveBeenCalledWith('pet_001', { name: '小强' })
    expect(wrapper.find('.overlay').exists()).toBe(false)
  })

  it('姓名为空时提示不调用 API', async () => {
    const wrapper = await createWrapper()
    await waitForRender(wrapper)

    await wrapper.find('.name').trigger('click')
    await wrapper.vm.$nextTick()

    const input = wrapper.find('.rename-input')
    await input.setValue('')

    await wrapper.find('.rename-btn.confirm').trigger('click')
    await wrapper.vm.$nextTick()

    expect(mockUpdatePet).not.toHaveBeenCalled()
    expect(wrapper.find('.overlay').exists()).toBe(true)
  })
})

describe('生日弹窗交互', () => {
  it('点击🎂标签弹出生日编辑弹窗', async () => {
    const wrapper = await createWrapper()
    await waitForRender(wrapper)

    // 用文本匹配找到生日标签（.tag 有 breed/birthday/weight 三个）
    const tags = wrapper.findAll('.tag')
    const birthdayTag = tags.find(t => t.text().includes('🎂'))!
    expect(birthdayTag).toBeDefined()
    expect(birthdayTag.text()).toContain('2020-03-15')

    await birthdayTag.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.overlay').exists()).toBe(true)
    expect(wrapper.find('.rename-title').text()).toBe('修改生日')
  })

  it('不改日期直接确认关闭弹窗不发请求', async () => {
    const wrapper = await createWrapper()
    await waitForRender(wrapper)

    const tags = wrapper.findAll('.tag')
    const birthdayTag = tags.find(t => t.text().includes('🎂'))!
    await birthdayTag.trigger('click')
    await wrapper.vm.$nextTick()

    // 弹窗已打开
    expect(wrapper.find('.overlay').exists()).toBe(true)

    await wrapper.find('.rename-btn.confirm').trigger('click')
    await new Promise(r => setTimeout(r, 50))
    await wrapper.vm.$nextTick()

    expect(mockUpdatePet).not.toHaveBeenCalled()
    expect(wrapper.find('.overlay').exists()).toBe(false)
  })
})
