<template>
  <view class="avatar-picker">
    <text class="current-avatar">{{ avatar || '🐾' }}</text>
    <view class="emoji-grid">
      <text
        v-for="emoji in emojis"
        :key="emoji"
        class="emoji-option"
        :class="{ active: avatar === emoji }"
        @tap="$emit('pick', emoji)"
      >{{ emoji }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
defineProps<{ avatar: string; emojis: string[] }>()
defineEmits<{ (e: 'pick', emoji: string): void }>()
</script>

<style scoped>
.avatar-picker {
  text-align: center;
}
.current-avatar {
  font-size: 100rpx;
  display: block;
  margin-bottom: 20rpx;
}
.emoji-grid {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 16rpx;
}
.emoji-option {
  font-size: 44rpx;
  padding: 10rpx 14rpx;
  border-radius: 16rpx;
  border: 2rpx solid transparent;
}
.emoji-option.active {
  border-color: #4CAF50;
  background: #f0f9f0;
}
</style>
