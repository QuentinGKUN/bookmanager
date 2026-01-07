<template>
  <div class="location-select">
    <el-select v-model="areaId" placeholder="请选择区域" @change="handleAreaChange" style="width: 150px; margin-right: 10px">
      <el-option v-for="area in areas" :key="area.id" :label="area.name" :value="area.id" />
    </el-select>
    <el-select v-model="bookshelfId" placeholder="请选择书架" @change="handleBookshelfChange" style="width: 150px; margin-right: 10px" :disabled="!areaId">
      <el-option v-for="bookshelf in bookshelves" :key="bookshelf.id" :label="bookshelf.name" :value="bookshelf.id" />
    </el-select>
    <el-select v-model="shelfLayerId" placeholder="请选择层数" @change="handleShelfLayerChange" style="width: 150px" :disabled="!bookshelfId">
      <el-option v-for="layer in shelfLayers" :key="layer.id" :label="layer.name" :value="layer.id" />
    </el-select>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { locationApi } from '../api/location'

const props = defineProps({
  modelValue: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['update:modelValue'])

const areaId = ref(null)
const bookshelfId = ref(null)
const shelfLayerId = ref(null)

const areas = ref([])
const bookshelves = ref([])
const shelfLayers = ref([])

const loadAreas = async () => {
  try {
    areas.value = await locationApi.area.list()
  } catch (error) {
    console.error('加载区域失败', error)
  }
}

const loadBookshelves = async () => {
  if (!areaId.value) {
    bookshelves.value = []
    return
  }
  try {
    bookshelves.value = await locationApi.bookshelf.list(areaId.value)
  } catch (error) {
    console.error('加载书架失败', error)
  }
}

const loadShelfLayers = async () => {
  if (!bookshelfId.value) {
    shelfLayers.value = []
    return
  }
  try {
    shelfLayers.value = await locationApi.shelfLayer.list(bookshelfId.value)
  } catch (error) {
    console.error('加载层数失败', error)
  }
}

const handleAreaChange = () => {
  bookshelfId.value = null
  shelfLayerId.value = null
  emit('update:modelValue', null)
  loadBookshelves()
}

const handleBookshelfChange = () => {
  shelfLayerId.value = null
  emit('update:modelValue', null)
  loadShelfLayers()
}

const handleShelfLayerChange = () => {
  emit('update:modelValue', shelfLayerId.value)
}

watch(() => props.modelValue, (val) => {
  if (val !== shelfLayerId.value) {
    // 如果需要根据shelfLayerId反向查找区域和书架，需要加载完整树结构
    // 这里简化处理，只更新shelfLayerId
    shelfLayerId.value = val
  }
})

onMounted(() => {
  loadAreas()
})
</script>

<style scoped>
.location-select {
  display: flex;
  align-items: center;
}
</style>





