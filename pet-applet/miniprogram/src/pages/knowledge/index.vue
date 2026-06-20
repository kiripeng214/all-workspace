<template>
  <view class="knowledge">
    <view class="search-bar">
      <input class="search-input" v-model="query" placeholder="搜索宠物知识..." @confirm="doSearch" />
      <view class="search-btn" @tap="doSearch">搜索</view>
    </view>

    <view v-if="petBreed" class="breed-tip">
      <text class="breed-tag">{{ petBreed }}相关知识</text>
    </view>

    <view v-if="loading" class="loading">正在搜索知识库...</view>

    <view v-else-if="answer" class="answer-card">
      <text class="answer-title">AI 回答</text>
      <text class="answer-text">{{ answer.answer }}</text>
      <view v-if="answer.sources && answer.sources.length" class="sources">
        <text class="sources-title">参考来源：</text>
        <text v-for="s in answer.sources" :key="s" class="source-item">{{ s }}</text>
      </view>
    </view>

    <view v-if="results.length" class="results">
      <view class="results-title">相关知识条目</view>
      <view v-for="r in results" :key="r.title" class="result-item">
        <text class="result-title">{{ r.title }}</text>
        <text class="result-content">{{ r.content }}</text>
        <view class="result-tags">
          <text v-for="t in r.tags" :key="t" class="tag">{{ t }}</text>
        </view>
      </view>
    </view>

    <view v-else-if="!loading && searched" class="empty">
      <text>未找到相关知识，试试换个关键词</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { searchKnowledge, type KnowledgeResult, type LLMAnswer } from '@/api/knowledge'

const query = ref('')
const petBreed = ref('')
const petName = ref('')
const results = ref<KnowledgeResult[]>([])
const answer = ref<LLMAnswer | null>(null)
const loading = ref(false)
const searched = ref(false)

onLoad((options: any) => {
  const breed = options?.breed ? decodeURIComponent(options.breed) : ''
  const name = options?.name ? decodeURIComponent(options.name) : ''
  if (breed) {
    petBreed.value = breed
    query.value = breed
  }
  if (name) {
    petName.value = name
  }
  if (petBreed.value) {
    doSearch()
  }
})

async function doSearch() {
  if (!query.value.trim()) return
  loading.value = true
  searched.value = true
  try {
    const data = await searchKnowledge({
      q: query.value.trim(),
      breed: petBreed.value || undefined,
    })
    results.value = data.results || []
    answer.value = data.answer || null
  } catch (e) {
    console.error('搜索异常:', e)
    results.value = []
    answer.value = null
    uni.showToast({ title: '搜索失败，请检查后端是否运行', icon: 'none' })
  } finally {
    loading.value = false
  }
}
</script>

<style>
.knowledge {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20rpx;
}
.search-bar {
  display: flex;
  gap: 16rpx;
  margin-bottom: 20rpx;
}
.search-input {
  flex: 1;
  background: #fff;
  border-radius: 12rpx;
  padding: 20rpx 24rpx;
  font-size: 28rpx;
  border: 1rpx solid #e0e0e0;
}
.search-btn {
  background: #4CAF50;
  color: #fff;
  border: none;
  padding: 20rpx 40rpx;
  border-radius: 12rpx;
  font-size: 28rpx;
}
.breed-tip {
  margin-bottom: 20rpx;
}
.breed-tag {
  display: inline-block;
  background: #e8f5e9;
  color: #4CAF50;
  font-size: 24rpx;
  padding: 8rpx 20rpx;
  border-radius: 20rpx;
}
.loading {
  text-align: center;
  color: #999;
  padding-top: 100rpx;
  font-size: 28rpx;
}
.answer-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
  border-left: 6rpx solid #4CAF50;
}
.answer-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  display: block;
  margin-bottom: 16rpx;
}
.answer-text {
  font-size: 28rpx;
  color: #444;
  line-height: 1.7;
  white-space: pre-wrap;
}
.sources {
  margin-top: 20rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid #f0f0f0;
}
.sources-title {
  font-size: 24rpx;
  color: #999;
  display: block;
  margin-bottom: 8rpx;
}
.source-item {
  font-size: 24rpx;
  color: #4CAF50;
  display: block;
  margin-bottom: 4rpx;
}
.results {
  margin-top: 10rpx;
}
.results-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #666;
  margin-bottom: 16rpx;
}
.result-item {
  background: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
}
.result-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  display: block;
  margin-bottom: 8rpx;
}
.result-content {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
  display: block;
}
.result-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8rpx;
  margin-top: 12rpx;
}
.tag {
  background: #f0f9f0;
  color: #4CAF50;
  font-size: 22rpx;
  padding: 4rpx 12rpx;
  border-radius: 12rpx;
}
.empty {
  text-align: center;
  color: #999;
  padding-top: 200rpx;
  font-size: 28rpx;
}
</style>
