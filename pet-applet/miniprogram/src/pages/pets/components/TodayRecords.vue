<template>
  <view class="section">
    <view class="section-header">
      <text class="section-title">今日喂养记录</text>
      <text class="section-more" @tap="goRecords">查看全部 ›</text>
    </view>
    <view v-if="records.length === 0" class="empty-section">暂无今日记录</view>
    <view v-for="r in records" :key="r.id" class="record-item">
      <text class="record-time">{{ r.time }}</text>
      <text class="record-food">{{ r.foodType }} {{ r.amount }}</text>
      <text class="record-note" v-if="r.notes">{{ r.notes }}</text>
    </view>
    <button class="btn small" @tap="showForm = true">+ 记录喂养</button>

    <view v-if="showForm" class="overlay" @tap="showForm = false">
      <view class="popup" @tap.stop>
        <text class="popup-title">记录喂养</text>
        <picker mode="time" :value="form.time" @change="onTimeChange">
          <view class="input picker-value">{{ form.time || '选择时间' }}</view>
        </picker>
        <input class="input" v-model="form.foodType" placeholder="食物类型" />
        <input class="input" v-model="form.amount" placeholder="分量" />
        <input class="input" v-model="form.notes" placeholder="备注" />
        <view class="popup-actions">
          <button class="popup-btn cancel" @tap="showForm = false">取消</button>
          <button class="popup-btn" @tap="submitRecord">提交</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { createRecord, type FeedingRecord } from '@/api'

const props = defineProps<{
  records: FeedingRecord[]
  petId: string
}>()

const emit = defineEmits<{
  (e: 'view-all'): void
  (e: 'record-created'): void
}>()

const showForm = ref(false)
const form = ref({ time: '', foodType: '', amount: '', notes: '' })

function goRecords() {
  emit('view-all')
}

function onTimeChange(e: any) {
  form.value.time = e.detail.value
}

async function submitRecord() {
  await createRecord(props.petId, form.value)
  showForm.value = false
  form.value = { time: '', foodType: '', amount: '', notes: '' }
  emit('record-created')
}
</script>

<style scoped>
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
.btn {
  background: #4CAF50;
  color: #fff;
  border: none;
  padding: 16rpx 40rpx;
  border-radius: 8rpx;
  font-size: 28rpx;
}
.btn.small {
  padding: 12rpx 30rpx;
  font-size: 26rpx;
  margin-top: 16rpx;
}
.overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
}
.popup {
  background: #fff;
  padding: 60rpx 48rpx;
  border-radius: 24rpx;
  width: 620rpx;
  max-height: 80vh;
  overflow-y: auto;
}
.popup-title {
  font-size: 36rpx;
  font-weight: 700;
  display: block;
  text-align: center;
  margin-bottom: 40rpx;
  color: #333;
}
.input {
  border: 2rpx solid #e0e0e0;
  border-radius: 12rpx;
  padding: 24rpx 28rpx;
  font-size: 30rpx;
  margin-bottom: 24rpx;
  width: 100%;
  box-sizing: border-box;
  min-height: 88rpx;
  color: #333;
  background: #fafafa;
}
.picker-value {
  display: flex;
  align-items: center;
  color: #333;
  min-height: 40rpx;
}
.popup-actions {
  display: flex;
  gap: 20rpx;
  margin-top: 40rpx;
}
.popup-btn {
  flex: 1;
  padding: 28rpx;
  border-radius: 16rpx;
  font-size: 32rpx;
  font-weight: 600;
  text-align: center;
  border: none;
  background: #4CAF50;
  color: #fff;
}
.popup-btn.cancel {
  background: #f5f5f5;
  color: #666;
}
</style>
