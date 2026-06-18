<template>
  <view></view>
</template>

<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
onLoad(() => {
  uni.redirectTo({ url: '/pages/pets/index' })
})
</script>
