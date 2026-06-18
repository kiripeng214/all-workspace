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

    <view class="add-area">
      <text class="add-title">{{ editingSchedule ? '编辑计划' : '添加计划' }}</text>
      <view class="add-row">
        <text class="add-label">时间</text>
        <picker mode="time" :value="form.time" @change="onTimeChange">
          <view class="add-input picker-value">{{ form.time || '选择时间' }}</view>
        </picker>
      </view>
      <view class="add-row">
        <text class="add-label">食物</text>
        <input class="add-input" v-model="form.foodType" placeholder="粮食" />
      </view>
      <view class="add-row">
        <text class="add-label">分量</text>
        <input class="add-input" v-model="form.amount" placeholder="一份" />
      </view>
      <view class="add-btns">
        <button class="btn cancel" v-if="editingSchedule" @tap="cancelEdit">取消</button>
        <button class="btn" @tap="onSubmit">{{ editingSchedule ? '保存' : '添加' }}</button>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getSchedules, createSchedule, updateSchedule, deleteSchedule, type FeedingSchedule } from '@/api'

const petId = ref('')
const petName = ref('')
const schedules = ref<FeedingSchedule[]>([])
const editingSchedule = ref<FeedingSchedule | null>(null)
const form = reactive({ time: '', foodType: '', amount: '' })

onLoad((options) => {
  petId.value = options?.petId || ''
  petName.value = options?.petName || ''
  uni.setNavigationBarTitle({ title: `${petName.value} - 喂养计划` })
  loadSchedules()
})

async function loadSchedules() {
  schedules.value = await getSchedules(petId.value)
}

function editSchedule(s: FeedingSchedule) {
  editingSchedule.value = s
  form.time = s.time
  form.foodType = s.foodType
  form.amount = s.amount
}

function cancelEdit() {
  editingSchedule.value = null
  form.time = ''
  form.foodType = ''
  form.amount = ''
}

function onTimeChange(e: any) {
  form.time = e.detail.value
}

async function onSubmit() {
  if (!form.time.trim()) {
    uni.showToast({ title: '请输入时间', icon: 'none' })
    return
  }
  if (editingSchedule.value) {
    await updateSchedule(editingSchedule.value.id, form)
  } else {
    await createSchedule(petId.value, form)
  }
  cancelEdit()
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
.add-area {
  background: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
  margin-top: 20rpx;
}
.add-title {
  font-size: 30rpx;
  font-weight: 600;
  display: block;
  margin-bottom: 16rpx;
}
.add-row {
  display: flex;
  align-items: center;
  margin-bottom: 12rpx;
}
.add-label {
  font-size: 26rpx;
  color: #666;
  width: 100rpx;
}
.add-input {
  flex: 1;
  border: 1rpx solid #ddd;
  border-radius: 6rpx;
  padding: 12rpx 16rpx;
  font-size: 26rpx;
}
.picker-value {
  color: #333;
  box-sizing: border-box;
}
.add-btns {
  display: flex;
  gap: 16rpx;
  margin-top: 16rpx;
}
.btn {
  flex: 1;
  background: #4CAF50;
  color: #fff;
  border: none;
  padding: 16rpx;
  border-radius: 8rpx;
  font-size: 28rpx;
}
.btn.cancel {
  background: #999;
}
</style>
