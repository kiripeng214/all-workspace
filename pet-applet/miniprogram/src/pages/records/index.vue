<template>
  <view class="records">
    <view v-if="records.length === 0" class="empty">暂无喂养记录</view>
    <view v-for="r in records" :key="r.id" class="item">
      <view class="item-main">
        <text class="item-time">{{ r.time }}</text>
        <text class="item-food">{{ r.foodType }} {{ r.amount }}</text>
        <text class="item-note" v-if="r.notes">{{ r.notes }}</text>
      </view>
      <text class="delete" @tap="confirmDelete(r.id)">删除</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getRecords, deleteRecord, type FeedingRecord } from '@/api'

const petId = ref('')
const petName = ref('')
const records = ref<FeedingRecord[]>([])

onLoad((options) => {
  petId.value = options?.petId || ''
  petName.value = options?.petName || ''
  uni.setNavigationBarTitle({ title: `${petName.value} - 喂养记录` })
  loadRecords()
})

async function loadRecords() {
  records.value = await getRecords(petId.value)
}

function confirmDelete(id: string) {
  uni.showModal({
    title: '确认删除',
    content: '确定删除该喂养记录？',
    success: async (res) => {
      if (res.confirm) {
        await deleteRecord(id)
        loadRecords()
      }
    },
  })
}
</script>

<style>
.records {
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
  align-items: flex-start;
}
.item-main {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  flex: 1;
}
.item-time {
  font-size: 28rpx;
  font-weight: 600;
  color: #4CAF50;
}
.item-food {
  font-size: 28rpx;
  color: #333;
}
.item-note {
  font-size: 24rpx;
  color: #999;
}
.delete {
  font-size: 26rpx;
  color: #f44336;
  padding: 8rpx;
}
</style>
