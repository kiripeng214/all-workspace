<template>
  <view class="detail">
    <view v-if="loading" class="loading">加载中...</view>
    <view v-else-if="pet" class="content">
      <view class="header">
        <text class="avatar">{{ pet.avatar }}</text>
        <text class="name" @tap="openRename">{{ pet.name }}</text>
        <view class="tags">
          <text class="tag" v-if="pet.breed">{{ pet.breed }}</text>
          <text class="tag" v-if="pet.birthday">🎂 {{ pet.birthday }}</text>
          <text class="tag" v-if="pet.weight">⚖️ {{ pet.weight }}</text>
        </view>
        <text class="notes" v-if="pet.notes">{{ pet.notes }}</text>
        <view class="actions">
          <button class="btn" @tap="goEdit">编辑</button>
          <button class="btn danger" @tap="confirmDelete">删除</button>
        </view>
      </view>

      <view class="section">
        <view class="section-header">
          <text class="section-title">今日喂养记录</text>
          <text class="section-more" @tap="goRecords">查看全部 ›</text>
        </view>
        <view v-if="todayRecords.length === 0" class="empty-section">暂无今日记录</view>
        <view v-for="r in todayRecords" :key="r.id" class="record-item">
          <text class="record-time">{{ r.time }}</text>
          <text class="record-food">{{ r.foodType }} {{ r.amount }}</text>
          <text class="record-note" v-if="r.notes">{{ r.notes }}</text>
        </view>
        <button class="btn small" @tap="showCreateRecord = true">+ 记录喂养</button>
      </view>

      <view class="section">
        <view class="section-header">
          <text class="section-title">喂养计划</text>
          <text class="section-more" @tap="goSchedules">管理 ›</text>
        </view>
        <view v-if="schedules.length === 0" class="empty-section">暂无喂养计划</view>
        <view v-for="s in schedules" :key="s.id" class="record-item">
          <text class="record-time">{{ s.time }}</text>
          <text class="record-food">{{ s.foodType }} {{ s.amount }}</text>
        </view>
      </view>
    </view>

    <uni-popup :show="showCreateRecord" @close="showCreateRecord = false">
      <view class="popup">
        <text class="popup-title">记录喂养</text>
        <picker mode="time" :value="recordForm.time" @change="onRecordTimeChange">
          <view class="input picker-value">{{ recordForm.time || '选择时间' }}</view>
        </picker>
        <input class="input" v-model="recordForm.foodType" placeholder="食物类型" />
        <input class="input" v-model="recordForm.amount" placeholder="分量" />
        <input class="input" v-model="recordForm.notes" placeholder="备注" />
        <button class="btn" @tap="submitRecord">提交</button>
      </view>
    </uni-popup>

    <uni-popup :show="showRename" @close="cancelRename">
      <view class="popup">
        <text class="popup-title">修改名称</text>
        <input class="input" v-model="renameName" placeholder="宠物名字" />
        <view class="popup-actions">
          <button class="btn cancel" @tap="cancelRename">取消</button>
          <button class="btn confirm" :disabled="renaming" @tap="submitRename">
            {{ renaming ? '保存中...' : '确认' }}
          </button>
        </view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { getPet, deletePet, getTodayRecords, createRecord, updatePet, getSchedules, type FeedingRecord, type FeedingSchedule, type Pet } from '@/api'

const petId = ref('')
const pet = ref<Pet | null>(null)
const loading = ref(false)
const todayRecords = ref<FeedingRecord[]>([])
const schedules = ref<FeedingSchedule[]>([])
const showCreateRecord = ref(false)
const recordForm = ref({ time: '', foodType: '', amount: '', notes: '' })
const showRename = ref(false)
const renameName = ref('')
const renaming = ref(false)

onLoad((options) => {
  petId.value = options?.id || ''
  loadData()
})

onShow(() => {
  if (petId.value) {
    loadData()
  }
})

async function loadData() {
  loading.value = true
  try {
    const [p, records, scheds] = await Promise.all([
      getPet(petId.value),
      getTodayRecords(petId.value),
      getSchedules(petId.value),
    ])
    pet.value = p
    todayRecords.value = records
    schedules.value = scheds
  } finally {
    loading.value = false
  }
}

function goEdit() {
  uni.navigateTo({ url: `/pages/pets/edit?id=${petId.value}` })
}

function confirmDelete() {
  uni.showModal({
    title: '确认删除',
    content: '删除后将同时删除所有喂养计划和记录，确定删除？',
    success: async (res) => {
      if (res.confirm) {
        await deletePet(petId.value)
        uni.navigateBack()
      }
    },
  })
}

function goRecords() {
  uni.navigateTo({ url: `/pages/records/index?petId=${petId.value}&petName=${pet.value?.name}` })
}

function goSchedules() {
  uni.navigateTo({ url: `/pages/schedules/index?petId=${petId.value}&petName=${pet.value?.name}` })
}

function openRename() {
  if (!pet.value) return
  renameName.value = pet.value.name
  showRename.value = true
}

async function submitRename() {
  const name = renameName.value.trim()
  if (!name) {
    uni.showToast({ title: '名称不能为空', icon: 'none', duration: 2000 })
    return
  }
  renaming.value = true
  try {
    await updatePet(petId.value, { name })
    uni.showToast({ title: '修改成功', icon: 'success' })
    showRename.value = false
    renameName.value = ''
    try {
      await loadData()
    } catch {
      // 页面刷新失败不影响修改结果
    }
  } catch {
    uni.showToast({ title: '修改失败', icon: 'error' })
  } finally {
    renaming.value = false
  }
}

function cancelRename() {
  showRename.value = false
  renameName.value = ''
}

function onRecordTimeChange(e: any) {
  recordForm.value.time = e.detail.value
}

async function submitRecord() {
  await createRecord(petId.value, recordForm.value)
  showCreateRecord.value = false
  recordForm.value = { time: '', foodType: '', amount: '', notes: '' }
  await loadData()
}
</script>

<style>
.detail {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20rpx;
}
.loading {
  text-align: center;
  padding-top: 300rpx;
  color: #999;
}
.header {
  background: #fff;
  border-radius: 16rpx;
  padding: 40rpx;
  text-align: center;
  margin-bottom: 20rpx;
}
.avatar {
  font-size: 100rpx;
}
.name {
  display: block;
  font-size: 40rpx;
  font-weight: 600;
  color: #333;
  margin-top: 16rpx;
  padding: 8rpx 20rpx;
  border: 2rpx dashed transparent;
  border-radius: 12rpx;
}
.name:active {
  border-color: #4CAF50;
  background: #f0f9f0;
}
.tags {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 10rpx;
  margin-top: 16rpx;
}
.tag {
  background: #f0f9f0;
  color: #4CAF50;
  font-size: 24rpx;
  padding: 6rpx 16rpx;
  border-radius: 20rpx;
}
.notes {
  display: block;
  color: #666;
  font-size: 26rpx;
  margin-top: 16rpx;
}
.actions {
  display: flex;
  gap: 20rpx;
  margin-top: 24rpx;
  justify-content: center;
}
.btn {
  background: #4CAF50;
  color: #fff;
  border: none;
  padding: 16rpx 40rpx;
  border-radius: 8rpx;
  font-size: 28rpx;
}
.btn.danger {
  background: #f44336;
}
.btn.small {
  padding: 12rpx 30rpx;
  font-size: 26rpx;
  margin-top: 16rpx;
}
.section {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}
.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}
.section-more {
  font-size: 26rpx;
  color: #4CAF50;
}
.empty-section {
  text-align: center;
  color: #999;
  font-size: 26rpx;
  padding: 20rpx 0;
}
.record-item {
  display: flex;
  align-items: center;
  padding: 12rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
  gap: 12rpx;
}
.record-time {
  color: #4CAF50;
  font-size: 28rpx;
  font-weight: 600;
  min-width: 100rpx;
}
.record-food {
  color: #333;
  font-size: 26rpx;
}
.record-note {
  color: #999;
  font-size: 24rpx;
  margin-left: auto;
}
.popup {
  background: #fff;
  padding: 40rpx;
  border-radius: 16rpx;
}
.popup-title {
  font-size: 32rpx;
  font-weight: 600;
  display: block;
  margin-bottom: 24rpx;
}
.input {
  border: 1rpx solid #ddd;
  border-radius: 8rpx;
  padding: 16rpx 20rpx;
  font-size: 28rpx;
  margin-bottom: 16rpx;
}
.picker-value {
  color: #333;
  box-sizing: border-box;
}
.popup-actions {
  display: flex;
  gap: 20rpx;
  margin-top: 16rpx;
}
.btn.cancel {
  background: #f5f5f5;
  color: #666;
  flex: 1;
}
.btn.confirm {
  background: #4CAF50;
  color: #fff;
  flex: 1;
}
.btn.confirm:disabled {
  opacity: 0.6;
}
</style>
