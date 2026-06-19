<template>
  <view class="edit">
    <view class="form">
      <AvatarPicker :avatar="form.avatar" :emojis="emojis" @pick="onPickEmoji" />

      <view class="field">
        <text class="label">名字 *</text>
        <input class="input" v-model="form.name" placeholder="宠物名字" />
      </view>

      <view class="field">
        <text class="label">品种</text>
        <picker :range="breedOptions[form.avatar] || []" @change="onBreedChange">
          <view class="picker">{{ form.breed || '选择品种' }}</view>
        </picker>
      </view>

      <view class="field">
        <text class="label">生日</text>
        <picker mode="date" :value="form.birthday" @change="onDateChange">
          <view class="picker">{{ form.birthday || '选择生日' }}</view>
        </picker>
      </view>

      <view class="field">
        <text class="label">体重</text>
        <view class="weight-row">
          <input class="input weight-input" type="digit" v-model="weightValue" placeholder="数值" />
          <picker class="unit-picker" :range="weightUnits" @change="onWeightUnitChange">
            <view class="picker unit-value">{{ weightUnit }}</view>
          </picker>
        </view>
      </view>

      <view class="field">
        <text class="label">备注</text>
        <input class="input" v-model="form.notes" placeholder="备注信息" />
      </view>

      <button class="submit" @tap="onSubmit">{{ isEdit ? '保存' : '添加' }}</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getPet, createPet, updatePet, getBreeds, type Pet } from '@/api'
import AvatarPicker from './components/AvatarPicker.vue'

const isEdit = ref(false)
const petId = ref('')
const form = reactive({ avatar: '', name: '', breed: '', birthday: '', weight: '', notes: '' })
const emojis = ref<string[]>([])
const breedOptions = ref<Record<string, string[]>>({})
const weightValue = ref('')
const weightUnit = ref('kg')
const weightUnits = ['kg', 'g']

const fallbackEmojis = ['🐶', '🐱', '🐰', '🐹', '🐦', '🐟', '🐢', '🦜', '🦊', '🐻']
emojis.value = fallbackEmojis

async function loadBreeds(emoji: string) {
  try {
    const meta = await getBreeds()
    emojis.value = meta.petEmojis
    breedOptions.value = meta.breedOptions
  } catch (err) {
    console.error('获取品种列表失败:', err)
  }
}

function onPickEmoji(emoji: string) {
  form.avatar = emoji
  form.breed = ''
  loadBreeds(emoji)
}

onLoad(async (options) => {
  if (options?.id) {
    isEdit.value = true
    petId.value = options.id
    const pet = await getPet(petId.value)
    Object.assign(form, {
      avatar: pet.avatar,
      name: pet.name,
      breed: pet.breed,
      birthday: pet.birthday,
      weight: pet.weight,
      notes: pet.notes,
    })
    // 解析体重：如 "15kg" → value=15, unit=kg
    const m = pet.weight?.match(/^([\d.]+)(kg|g)$/)
    if (m) {
      weightValue.value = m[1]
      weightUnit.value = m[2]
    }
    if (pet.avatar) {
      await loadBreeds(pet.avatar)
    }
  }
})

function onBreedChange(e: any) {
  form.breed = (breedOptions.value[form.avatar] || [])[e.detail.value] || ''
}

function onDateChange(e: any) {
  form.birthday = e.detail.value
}

function onWeightUnitChange(e: any) {
  weightUnit.value = weightUnits[e.detail.value] || 'kg'
}

async function onSubmit() {
  if (!form.name.trim()) {
    uni.showToast({ title: '请输入宠物名字', icon: 'none' })
    return
  }
  // 拼接体重数值 + 单位
  form.weight = weightValue.value ? `${weightValue.value}${weightUnit.value}` : ''
  if (isEdit.value) {
    await updatePet(petId.value, form)
    uni.showToast({ title: '保存成功', icon: 'success' })
  } else {
    await createPet(form)
    uni.showToast({ title: '添加成功', icon: 'success' })
  }
  uni.navigateBack()
}
</script>

<style>
.edit {
  background: #f5f5f5;
  min-height: 100vh;
  padding: 30rpx;
}
.form {
  background: #fff;
  border-radius: 24rpx;
  padding: 32rpx;
  display: flex;
  flex-direction: column;
  gap: 32rpx;
}
.field {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}
.label {
  font-size: 28rpx;
  font-weight: 600;
  color: #5D4037;
}
.input, .picker {
  border: 2rpx solid #e0e0e0;
  border-radius: 16rpx;
  padding: 24rpx 28rpx;
  font-size: 30rpx;
  color: #333;
  width: 100%;
  box-sizing: border-box;
  min-height: 88rpx;
}
.picker {
  display: flex;
  align-items: center;
  color: #999;
}
.weight-row {
  display: flex;
  gap: 16rpx;
  align-items: center;
}
.weight-input {
  flex: 1;
}
.unit-picker {
  min-width: 100rpx;
}
.unit-value {
  border: 2rpx solid #e0e0e0;
  border-radius: 16rpx;
  padding: 24rpx 28rpx;
  font-size: 30rpx;
  color: #333;
  min-height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fafafa;
  box-sizing: border-box;
}
.submit {
  background: linear-gradient(135deg, #FF8A65, #FF7043);
  color: #fff;
  border: none;
  padding: 28rpx;
  border-radius: 50rpx;
  font-size: 32rpx;
  font-weight: 600;
  width: 100%;
  box-shadow: 0 8rpx 24rpx rgba(255, 112, 67, 0.35);
}
.submit:active {
  opacity: 0.85;
}
</style>
