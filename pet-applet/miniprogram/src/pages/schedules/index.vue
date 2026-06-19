<template>
  <view class="schedules">
    <view v-if="schedules.length === 0" class="empty">暂无喂养计划</view>
    <view v-for="s in schedules" :key="s.id" class="item">
      <view class="item-info">
        <text class="item-time">{{ s.time }}</text>
        <text class="item-food">{{ s.foodType }} {{ s.amount }}</text>
      </view>
      <view class="item-actions">
        <text class="action edit" @tap="editSchedule(s)">编辑</text>
        <text class="action delete" @tap="confirmDelete(s.id)">删除</text>
      </view>
    </view>

    <view class="fab" @tap="addSchedule">
      <text class="fab-icon">+</text>
      <text class="fab-text">添加计划</text>
    </view>

    <ScheduleForm
      :show="showForm"
      :editing="!!editingSchedule"
      :initial="editingSchedule"
      @submit="onFormSubmit"
      @cancel="closeForm"
      @close="closeForm"
    />
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getSchedules, createSchedule, updateSchedule, deleteSchedule, type FeedingSchedule } from '@/api'
import ScheduleForm from './components/ScheduleForm.vue'

const petId = ref('')
const petName = ref('')
const schedules = ref<FeedingSchedule[]>([])
const editingSchedule = ref<FeedingSchedule | null>(null)
const showForm = ref(false)

onLoad((options) => {
  petId.value = options?.petId || ''
  petName.value = options?.petName || ''
  uni.setNavigationBarTitle({ title: `${petName.value} - 喂养计划` })
  loadSchedules()
})

async function loadSchedules() {
  schedules.value = await getSchedules(petId.value)
}

function addSchedule() {
  editingSchedule.value = null
  showForm.value = true
}

function editSchedule(s: FeedingSchedule) {
  editingSchedule.value = s
  showForm.value = true
}

function closeForm() {
  showForm.value = false
  editingSchedule.value = null
}

async function onFormSubmit(form: { time: string; foodType: string; amount: string }) {
  if (editingSchedule.value) {
    await updateSchedule(editingSchedule.value.id, form)
  } else {
    await createSchedule(petId.value, form)
  }
  closeForm()
  loadSchedules()
}

function confirmDelete(id: string) {
  uni.showModal({
    title: '确认删除',
    content: '确定删除该喂养计划？',
    success: async (res) => {
      if (res.confirm) {
        await deleteSchedule(id)
        loadSchedules()
      }
    },
  })
}
</script>

<style>
.schedules {
  background: #f5f5f5;
  min-height: 100vh;
  padding: 20rpx;
}
.empty {
  text-align: center;
  color: #999;
  padding-top: 200rpx;
  font-size: 28rpx;
}
.item {
  background: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.item-info {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}
.item-time {
  font-size: 32rpx;
  font-weight: 600;
  color: #4CAF50;
}
.item-food {
  font-size: 26rpx;
  color: #666;
}
.item-actions {
  display: flex;
  gap: 16rpx;
}
.action {
  font-size: 26rpx;
  padding: 8rpx 16rpx;
  border-radius: 6rpx;
}
.action.edit { color: #4CAF50; }
.action.delete { color: #f44336; }
.fab {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  background: #4CAF50;
  color: #fff;
  padding: 20rpx;
  border-radius: 12rpx;
  font-size: 28rpx;
  margin-top: 16rpx;
}
.fab:active {
  opacity: 0.8;
}
.fab-icon {
  font-size: 36rpx;
  font-weight: 600;
}
.fab-text {
  font-size: 28rpx;
}
</style>
