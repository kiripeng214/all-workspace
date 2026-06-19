<template>
  <view class="header">
    <text class="avatar">{{ pet.avatar }}</text>
    <view class="name" @click="openRename">{{ pet.name }}</view>
    <view class="tags">
      <text class="tag" v-if="pet.breed">{{ pet.breed }}</text>
      <text class="tag" v-if="pet.birthday" @click="openBirthdayEdit">🎂 {{ pet.birthday }}</text>
      <text class="tag" v-if="pet.weight">⚖️ {{ pet.weight }}</text>
    </view>
    <text class="notes" v-if="pet.notes">{{ pet.notes }}</text>
    <view class="actions">
      <button class="btn" @tap="goEdit">编辑</button>
      <button class="btn danger" @tap="confirmDelete">删除</button>
    </view>

    <view v-if="showRename" class="overlay" @click="cancelRename">
      <view class="overlay-box" @click.stop>
        <text class="overlay-title">修改姓名</text>
        <input class="overlay-input" v-model="renameName" placeholder="宠物姓名" />
        <view class="overlay-actions">
          <button class="overlay-btn cancel" @click="cancelRename">取消</button>
          <button class="overlay-btn confirm" :disabled="renaming" @click="submitRename">
            {{ renaming ? '保存中...' : '确认' }}
          </button>
        </view>
      </view>
    </view>

    <view v-if="showBirthday" class="overlay" @click="cancelBirthday">
      <view class="overlay-box" @click.stop>
        <text class="overlay-title">修改生日</text>
        <picker mode="date" :value="birthdayValue" @change="onBirthdayChange">
          <view class="overlay-input picker-value">{{ birthdayValue || '选择日期' }}</view>
        </picker>
        <view class="overlay-actions">
          <button class="overlay-btn cancel" @click="cancelBirthday">取消</button>
          <button class="overlay-btn confirm" :disabled="savingBirthday" @click="submitBirthday">
            {{ savingBirthday ? '保存中...' : '确认' }}
          </button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { updatePet, type Pet } from '@/api'

const props = defineProps<{ pet: Pet }>()
const emit = defineEmits<{
  (e: 'edit'): void
  (e: 'delete'): void
  (e: 'name-changed'): void
  (e: 'birthday-changed'): void
}>()

const showRename = ref(false)
const renameName = ref('')
const renaming = ref(false)
const showBirthday = ref(false)
const birthdayValue = ref('')
const savingBirthday = ref(false)

function goEdit() {
  emit('edit')
}

function confirmDelete() {
  uni.showModal({
    title: '确认删除',
    content: '删除后将同时删除所有喂养计划和记录，确定删除？',
    success: (res) => {
      if (res.confirm) {
        emit('delete')
      }
    },
  })
}

function openRename() {
  renameName.value = props.pet.name
  showRename.value = true
}

async function submitRename() {
  const name = renameName.value.trim()
  if (!name) {
    uni.showToast({ title: '姓名不能为空', icon: 'none' })
    return
  }
  renaming.value = true
  try {
    await updatePet(props.pet.id, { name })
    uni.showToast({ title: '修改成功', icon: 'success' })
    showRename.value = false
    emit('name-changed')
  } catch {
    uni.showToast({ title: '修改失败', icon: 'error' })
  } finally {
    renaming.value = false
  }
}

function cancelRename() {
  showRename.value = false
}

function openBirthdayEdit() {
  birthdayValue.value = props.pet.birthday
  showBirthday.value = true
}

function onBirthdayChange(e: any) {
  birthdayValue.value = e.detail.value
}

async function submitBirthday() {
  if (!birthdayValue.value) {
    uni.showToast({ title: '请选择日期', icon: 'none' })
    return
  }
  if (birthdayValue.value === props.pet.birthday) {
    showBirthday.value = false
    return
  }
  savingBirthday.value = true
  try {
    await updatePet(props.pet.id, { birthday: birthdayValue.value })
    uni.showToast({ title: '修改成功', icon: 'success' })
    showBirthday.value = false
    emit('birthday-changed')
  } catch {
    uni.showToast({ title: '修改失败', icon: 'error' })
  } finally {
    savingBirthday.value = false
  }
}

function cancelBirthday() {
  showBirthday.value = false
}
</script>

<style scoped>
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
.tag:active {
  background: #d0ebd0;
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
.overlay-box {
  background: #fff;
  padding: 48rpx;
  border-radius: 20rpx;
  width: 560rpx;
}
.overlay-title {
  font-size: 32rpx;
  font-weight: 600;
  display: block;
  text-align: center;
  margin-bottom: 30rpx;
}
.overlay-input {
  display: block;
  border: 2rpx solid #ddd;
  border-radius: 12rpx;
  padding: 0 24rpx;
  font-size: 30rpx;
  height: 80rpx;
  line-height: 80rpx;
  width: 100%;
  box-sizing: border-box;
  margin-bottom: 30rpx;
}
.picker-value {
  display: flex;
  align-items: center;
  color: #333;
}
.overlay-actions {
  display: flex;
  gap: 20rpx;
}
.overlay-btn {
  flex: 1;
  padding: 22rpx;
  border-radius: 12rpx;
  font-size: 30rpx;
  text-align: center;
  border: none;
}
.overlay-btn.cancel {
  background: #f5f5f5;
  color: #666;
}
.overlay-btn.confirm {
  background: #4CAF50;
  color: #fff;
}
.overlay-btn.confirm:disabled {
  opacity: 0.5;
}
</style>
