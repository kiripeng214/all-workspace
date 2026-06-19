<template>
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
</template>

<script setup lang="ts">
import type { FeedingSchedule } from '@/api'

defineProps<{ schedules: FeedingSchedule[] }>()
const emit = defineEmits<{ (e: 'manage'): void }>()

function goSchedules() {
  emit('manage')
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
</style>
