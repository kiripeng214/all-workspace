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
        @longpress="openRename(pet)"
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
import { onShow } from '@dcloudio/uni-app'
import { getPets, updatePet, type Pet } from '@/api'

const pets = ref<Pet[]>([])
const loading = ref(false)
const showRename = ref(false)
const renamePet = ref<Pet | null>(null)
const renameName = ref('')
const renaming = ref(false)
let longPressed = false

async function loadPets() {
  loading.value = true
  try {
    pets.value = await getPets()
  } finally {
    loading.value = false
  }
}

function goDetail(id: string) {
  if (longPressed) {
    longPressed = false
    return
  }
  uni.navigateTo({ url: `/pages/pets/detail?id=${id}` })
}

function goCreate() {
  uni.navigateTo({ url: '/pages/pets/edit' })
}

function openRename(pet: Pet) {
  longPressed = true
  renamePet.value = pet
  renameName.value = pet.name
  showRename.value = true
}

async function submitRename() {
  const name = renameName.value.trim()
  if (!name) {
    uni.showToast({ title: '名称不能为空', icon: 'none' })
    return
  }
  if (!renamePet.value) return

  renaming.value = true
  try {
    await updatePet(renamePet.value.id, { name })
    uni.showToast({ title: '修改成功', icon: 'success' })
    showRename.value = false
    renamePet.value = null
    await loadPets()
  } catch (err) {
    uni.showToast({ title: '修改失败', icon: 'error' })
  } finally {
    renaming.value = false
  }
}

function cancelRename() {
  showRename.value = false
  renamePet.value = null
  renameName.value = ''
}

onShow(() => {
  loadPets()
})
</script>

<style scoped>
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
.popup {
  background: #fff;
  padding: 40rpx;
  border-radius: 16rpx;
  width: 560rpx;
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
  margin-bottom: 24rpx;
}
.popup-actions {
  display: flex;
  gap: 20rpx;
}
.btn {
  flex: 1;
  padding: 20rpx;
  border-radius: 8rpx;
  font-size: 28rpx;
  text-align: center;
}
.btn.cancel {
  background: #f5f5f5;
  color: #666;
}
.btn.confirm {
  background: #4CAF50;
  color: #fff;
}
.btn.confirm:disabled {
  opacity: 0.6;
}
</style>
