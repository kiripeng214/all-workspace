<template>
  <view class="pet-list">
    <view v-if="loading" class="loading">
      <text>加载中...</text>
    </view>
    <view v-else-if="pets.length === 0" class="empty">
      <text class="empty-icon">🐾</text>
      <text class="empty-text">还没有宠物，点击下方添加</text>
    </view>
    <view v-else class="list">
      <view
        v-for="pet in pets"
        :key="pet.id"
        class="card"
        @tap="goDetail(pet.id)"
      >
        <text class="avatar">{{ pet.avatar }}</text>
        <view class="info">
          <text class="name">{{ pet.name }}</text>
          <text class="breed" v-if="pet.breed">{{ pet.breed }}</text>
        </view>
        <text class="arrow">›</text>
      </view>
    </view>
    <view class="fab" @tap="goCreate">
      <text class="fab-icon">+</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getPets, type Pet } from '@/api'

const pets = ref<Pet[]>([])
const loading = ref(false)

async function loadPets() {
  loading.value = true
  try {
    pets.value = await getPets()
  } finally {
    loading.value = false
  }
}

function goDetail(id: string) {
  uni.navigateTo({ url: `/pages/pets/detail?id=${id}` })
}

function goCreate() {
  uni.navigateTo({ url: '/pages/pets/edit' })
}

onShow(() => {
  loadPets()
})
</script>

<style>
.pet-list {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20rpx;
}
.loading, .empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding-top: 300rpx;
}
.empty-icon {
  font-size: 100rpx;
  margin-bottom: 20rpx;
}
.empty-text {
  color: #999;
  font-size: 28rpx;
}
.list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}
.card {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  display: flex;
  align-items: center;
  box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.06);
}
.avatar {
  font-size: 60rpx;
  margin-right: 24rpx;
}
.info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}
.name {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
}
.breed {
  font-size: 26rpx;
  color: #999;
}
.arrow {
  font-size: 40rpx;
  color: #ccc;
}
.fab {
  position: fixed;
  right: 40rpx;
  bottom: 80rpx;
  width: 100rpx;
  height: 100rpx;
  background: #4CAF50;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4rpx 16rpx rgba(76,175,80,0.4);
}
.fab-icon {
  font-size: 50rpx;
  color: #fff;
}
</style>
