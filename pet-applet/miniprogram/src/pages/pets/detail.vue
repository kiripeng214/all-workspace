<template>
  <view class="detail">
    <view v-if="loading" class="loading">加载中...</view>
    <view v-else-if="pet" class="content">
      <PetInfoCard
        :pet="pet"
        @edit="goEdit"
        @delete="handleDelete"
        @name-changed="loadData"
        @birthday-changed="loadData"
        @knowledge="goKnowledge"
      />
      <TodayRecords
        :records="todayRecords"
        :pet-id="petId"
        @view-all="goRecords"
        @record-created="loadData"
      />
      <ScheduleList
        :schedules="schedules"
        @manage="goSchedules"
      />
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { getPet, deletePet, getTodayRecords, getSchedules, type FeedingRecord, type FeedingSchedule, type Pet } from '@/api'
import PetInfoCard from './components/PetInfoCard.vue'
import TodayRecords from './components/TodayRecords.vue'
import ScheduleList from './components/ScheduleList.vue'

const petId = ref('')
const pet = ref<Pet | null>(null)
const loading = ref(false)
const todayRecords = ref<FeedingRecord[]>([])
const schedules = ref<FeedingSchedule[]>([])

onLoad((options) => {
  petId.value = options?.id || ''
  loadData()
})

onShow(() => {
  if (petId.value) {
    loadData()
  }
})

async function loadData() {
  loading.value = true
  try {
    const [p, records, scheds] = await Promise.all([
      getPet(petId.value),
      getTodayRecords(petId.value),
      getSchedules(petId.value),
    ])
    pet.value = p
    todayRecords.value = records
    schedules.value = scheds
  } finally {
    loading.value = false
  }
}

function goEdit() {
  uni.navigateTo({ url: `/pages/pets/edit?id=${petId.value}` })
}

async function handleDelete() {
  await deletePet(petId.value)
  uni.navigateBack()
}

function goRecords() {
  uni.navigateTo({ url: `/pages/records/index?petId=${petId.value}&petName=${pet.value?.name}` })
}

function goSchedules() {
  uni.navigateTo({ url: `/pages/schedules/index?petId=${petId.value}&petName=${pet.value?.name}` })
}

function goKnowledge() {
  const breed = encodeURIComponent(pet.value?.breed || '')
  const name = encodeURIComponent(pet.value?.name || '')
  uni.navigateTo({
    url: `/pages/knowledge/index?breed=${breed}&name=${name}`,
    fail(err) { console.error('跳转知识库失败:', err) },
  })
}
</script>

<style>
.detail {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20rpx;
}
.loading {
  text-align: center;
  padding-top: 300rpx;
  color: #999;
}
</style>
