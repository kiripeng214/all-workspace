<template>
  <view v-if="show" class="overlay" @tap="$emit('close')">
    <view class="popup" @tap.stop>
      <text class="popup-title">{{ editing ? '编辑计划' : '添加计划' }}</text>
      <view class="field">
        <text class="label">时间</text>
        <picker mode="time" :value="localForm.time" @change="onTimeChange">
          <view class="input picker-value">{{ localForm.time || '选择时间' }}</view>
        </picker>
      </view>
      <view class="field">
        <text class="label">食物</text>
        <input class="input" v-model="localForm.foodType" placeholder="粮食" />
      </view>
      <view class="field">
        <text class="label">分量</text>
        <input class="input" v-model="localForm.amount" placeholder="一份" />
      </view>
      <view class="actions">
        <button class="btn cancel" @tap="$emit('cancel')">取消</button>
        <button class="btn" @tap="onSubmit">{{ editing ? '保存' : '添加' }}</button>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'
import type { FeedingSchedule } from '@/api'

function resetForm() {
  localForm.time = ''
  localForm.foodType = ''
  localForm.amount = ''
}

const props = defineProps<{
  show: boolean
  editing: boolean
  initial: FeedingSchedule | null
}>()
const emit = defineEmits<{
  (e: 'submit', form: { time: string; foodType: string; amount: string }): void
  (e: 'cancel'): void
  (e: 'close'): void
}>()

const localForm = reactive({ time: '', foodType: '', amount: '' })

// 弹出框打开时，编辑模式预填数据，添加模式清空
watch(() => props.show, (val) => {
  if (val && props.initial) {
    localForm.time = props.initial.time
    localForm.foodType = props.initial.foodType
    localForm.amount = props.initial.amount
  } else if (val) {
    resetForm()
  }
})

function onTimeChange(e: any) {
  localForm.time = e.detail.value
}

function onSubmit() {
  if (!localForm.time.trim()) {
    uni.showToast({ title: '请输入时间', icon: 'none' })
    return
  }
  if (!localForm.foodType.trim()) {
    uni.showToast({ title: '请输入食物类型', icon: 'none' })
    return
  }
  if (!localForm.amount.trim()) {
    uni.showToast({ title: '请输入分量', icon: 'none' })
    return
  }
  emit('submit', { time: localForm.time, foodType: localForm.foodType, amount: localForm.amount })
}
</script>

<style scoped>
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
.field {
  margin-bottom: 24rpx;
}
.label {
  font-size: 28rpx;
  color: #666;
  display: block;
  margin-bottom: 12rpx;
}
.input {
  border: 2rpx solid #e0e0e0;
  border-radius: 12rpx;
  padding: 24rpx 28rpx;
  font-size: 30rpx;
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
.actions {
  display: flex;
  gap: 20rpx;
  margin-top: 40rpx;
}
.btn {
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
.btn.cancel {
  background: #f5f5f5;
  color: #666;
}
</style>
