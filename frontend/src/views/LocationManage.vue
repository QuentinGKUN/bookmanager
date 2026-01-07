<template>
  <div class="location-manage">
    <el-container>
      <el-header>
        <h1>位置管理</h1>
      </el-header>
      <el-main>
        <el-tabs v-model="activeTab">
          <!-- 区域管理 -->
          <el-tab-pane label="区域管理" name="area">
            <el-card>
              <div style="margin-bottom: 20px">
                <el-button type="primary" @click="handleAddArea">新增区域</el-button>
              </div>
              <el-table :data="areaList" border>
                <el-table-column prop="name" label="区域名称" />
                <el-table-column label="操作" width="150">
                  <template #default="scope">
                    <el-button link type="primary" @click="handleEditArea(scope.row)">编辑</el-button>
                    <el-button link type="danger" @click="handleDeleteArea(scope.row)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-card>
          </el-tab-pane>

          <!-- 书架管理 -->
          <el-tab-pane label="书架管理" name="bookshelf">
            <el-card>
              <div style="margin-bottom: 20px">
                <el-button type="primary" @click="handleAddBookshelf">新增书架</el-button>
              </div>
              <el-table :data="bookshelfList" border>
                <el-table-column prop="area.name" label="所属区域" />
                <el-table-column prop="name" label="书架名称" />
                <el-table-column label="操作" width="150">
                  <template #default="scope">
                    <el-button link type="primary" @click="handleEditBookshelf(scope.row)">编辑</el-button>
                    <el-button link type="danger" @click="handleDeleteBookshelf(scope.row)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-card>
          </el-tab-pane>

          <!-- 层数管理 -->
          <el-tab-pane label="层数管理" name="shelfLayer">
            <el-card>
              <div style="margin-bottom: 20px">
                <el-button type="primary" @click="handleAddShelfLayer">新增层数</el-button>
              </div>
              <el-table :data="shelfLayerList" border>
                <el-table-column prop="bookshelf.area.name" label="所属区域" />
                <el-table-column prop="bookshelf.name" label="所属书架" />
                <el-table-column prop="name" label="层数名称" />
                <el-table-column label="操作" width="150">
                  <template #default="scope">
                    <el-button link type="primary" @click="handleEditShelfLayer(scope.row)">编辑</el-button>
                    <el-button link type="danger" @click="handleDeleteShelfLayer(scope.row)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-card>
          </el-tab-pane>
        </el-tabs>

        <!-- 区域对话框 -->
        <el-dialog v-model="areaDialogVisible" :title="areaDialogTitle" width="400px">
          <el-form :model="areaForm" :rules="areaRules" ref="areaFormRef" label-width="100px">
            <el-form-item label="区域名称" prop="name">
              <el-input v-model="areaForm.name" />
            </el-form-item>
          </el-form>
          <template #footer>
            <el-button @click="areaDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="handleSaveArea">保存</el-button>
          </template>
        </el-dialog>

        <!-- 书架对话框 -->
        <el-dialog v-model="bookshelfDialogVisible" :title="bookshelfDialogTitle" width="400px">
          <el-form :model="bookshelfForm" :rules="bookshelfRules" ref="bookshelfFormRef" label-width="100px">
            <el-form-item label="所属区域" prop="area_id">
              <el-select v-model="bookshelfForm.area_id" placeholder="请选择区域" style="width: 100%">
                <el-option v-for="area in areas" :key="area.id" :label="area.name" :value="area.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="书架名称" prop="name">
              <el-input v-model="bookshelfForm.name" />
            </el-form-item>
          </el-form>
          <template #footer>
            <el-button @click="bookshelfDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="handleSaveBookshelf">保存</el-button>
          </template>
        </el-dialog>

        <!-- 层数对话框 -->
        <el-dialog v-model="shelfLayerDialogVisible" :title="shelfLayerDialogTitle" width="400px">
          <el-form :model="shelfLayerForm" :rules="shelfLayerRules" ref="shelfLayerFormRef" label-width="100px">
            <el-form-item label="所属书架" prop="bookshelf_id">
              <el-select v-model="shelfLayerForm.bookshelf_id" placeholder="请选择书架" style="width: 100%" @change="handleBookshelfSelectChange">
                <el-option v-for="bookshelf in bookshelves" :key="bookshelf.id" :label="bookshelf.name" :value="bookshelf.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="层数名称" prop="name">
              <el-input v-model="shelfLayerForm.name" />
            </el-form-item>
          </el-form>
          <template #footer>
            <el-button @click="shelfLayerDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="handleSaveShelfLayer">保存</el-button>
          </template>
        </el-dialog>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { locationApi } from '../api/location'

const activeTab = ref('area')

const areaList = ref([])
const bookshelfList = ref([])
const shelfLayerList = ref([])
const areas = ref([])
const bookshelves = ref([])

const areaDialogVisible = ref(false)
const areaDialogTitle = ref('新增区域')
const areaFormRef = ref(null)
const areaForm = reactive({ name: '' })
const areaRules = {
  name: [{ required: true, message: '请输入区域名称', trigger: 'blur' }]
}
let currentAreaId = null

const bookshelfDialogVisible = ref(false)
const bookshelfDialogTitle = ref('新增书架')
const bookshelfFormRef = ref(null)
const bookshelfForm = reactive({ area_id: null, name: '' })
const bookshelfRules = {
  area_id: [{ required: true, message: '请选择区域', trigger: 'change' }],
  name: [{ required: true, message: '请输入书架名称', trigger: 'blur' }]
}
let currentBookshelfId = null

const shelfLayerDialogVisible = ref(false)
const shelfLayerDialogTitle = ref('新增层数')
const shelfLayerFormRef = ref(null)
const shelfLayerForm = reactive({ bookshelf_id: null, name: '' })
const shelfLayerRules = {
  bookshelf_id: [{ required: true, message: '请选择书架', trigger: 'change' }],
  name: [{ required: true, message: '请输入层数名称', trigger: 'blur' }]
}
let currentShelfLayerId = null

const loadAreas = async () => {
  try {
    areas.value = await locationApi.area.list()
    if (activeTab.value === 'area') {
      areaList.value = areas.value
    }
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const loadBookshelves = async () => {
  try {
    const allBookshelves = []
    for (const area of areas.value) {
      const list = await locationApi.bookshelf.list(area.id)
      allBookshelves.push(...list.map(b => ({ ...b, area })))
    }
    bookshelfList.value = allBookshelves
    bookshelves.value = allBookshelves
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const loadShelfLayers = async () => {
  try {
    const allLayers = []
    for (const bookshelf of bookshelves.value) {
      const list = await locationApi.shelfLayer.list(bookshelf.id)
      allLayers.push(...list.map(l => ({ ...l, bookshelf })))
    }
    shelfLayerList.value = allLayers
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const handleAddArea = () => {
  currentAreaId = null
  areaDialogTitle.value = '新增区域'
  areaForm.name = ''
  areaDialogVisible.value = true
}

const handleEditArea = (row) => {
  currentAreaId = row.id
  areaDialogTitle.value = '编辑区域'
  areaForm.name = row.name
  areaDialogVisible.value = true
}

const handleSaveArea = async () => {
  try {
    await areaFormRef.value.validate()
    if (currentAreaId) {
      await locationApi.area.update(currentAreaId, { name: areaForm.name })
      ElMessage.success('更新成功')
    } else {
      await locationApi.area.create({ name: areaForm.name })
      ElMessage.success('创建成功')
    }
    areaDialogVisible.value = false
    loadAreas()
  } catch (error) {
    if (error !== false) {
      ElMessage.error(error.message || '保存失败')
    }
  }
}

const handleDeleteArea = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该区域吗？', '提示', { type: 'warning' })
    await locationApi.area.delete(row.id)
    ElMessage.success('删除成功')
    loadAreas()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const handleAddBookshelf = () => {
  currentBookshelfId = null
  bookshelfDialogTitle.value = '新增书架'
  bookshelfForm.area_id = null
  bookshelfForm.name = ''
  bookshelfDialogVisible.value = true
}

const handleEditBookshelf = (row) => {
  currentBookshelfId = row.id
  bookshelfDialogTitle.value = '编辑书架'
  bookshelfForm.area_id = row.area_id
  bookshelfForm.name = row.name
  bookshelfDialogVisible.value = true
}

const handleSaveBookshelf = async () => {
  try {
    await bookshelfFormRef.value.validate()
    if (currentBookshelfId) {
      await locationApi.bookshelf.update(currentBookshelfId, bookshelfForm)
      ElMessage.success('更新成功')
    } else {
      await locationApi.bookshelf.create(bookshelfForm)
      ElMessage.success('创建成功')
    }
    bookshelfDialogVisible.value = false
    loadBookshelves()
  } catch (error) {
    if (error !== false) {
      ElMessage.error(error.message || '保存失败')
    }
  }
}

const handleDeleteBookshelf = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该书架吗？', '提示', { type: 'warning' })
    await locationApi.bookshelf.delete(row.id)
    ElMessage.success('删除成功')
    loadBookshelves()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const handleAddShelfLayer = () => {
  currentShelfLayerId = null
  shelfLayerDialogTitle.value = '新增层数'
  shelfLayerForm.bookshelf_id = null
  shelfLayerForm.name = ''
  shelfLayerDialogVisible.value = true
}

const handleEditShelfLayer = (row) => {
  currentShelfLayerId = row.id
  shelfLayerDialogTitle.value = '编辑层数'
  shelfLayerForm.bookshelf_id = row.bookshelf_id
  shelfLayerForm.name = row.name
  shelfLayerDialogVisible.value = true
}

const handleBookshelfSelectChange = () => {
  // 当选择书架变化时，可以加载对应的层数列表
}

const handleSaveShelfLayer = async () => {
  try {
    await shelfLayerFormRef.value.validate()
    if (currentShelfLayerId) {
      await locationApi.shelfLayer.update(currentShelfLayerId, shelfLayerForm)
      ElMessage.success('更新成功')
    } else {
      await locationApi.shelfLayer.create(shelfLayerForm)
      ElMessage.success('创建成功')
    }
    shelfLayerDialogVisible.value = false
    loadShelfLayers()
  } catch (error) {
    if (error !== false) {
      ElMessage.error(error.message || '保存失败')
    }
  }
}

const handleDeleteShelfLayer = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该层数吗？', '提示', { type: 'warning' })
    await locationApi.shelfLayer.delete(row.id)
    ElMessage.success('删除成功')
    loadShelfLayers()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

onMounted(() => {
  loadAreas()
  loadBookshelves()
  loadShelfLayers()
})
</script>

<style scoped>
.location-manage {
  min-height: 100vh;
  background: #f5f5f5;
}

.el-header h1 {
  margin: 0;
  line-height: 60px;
  font-size: 24px;
}

.el-main {
  padding: 20px;
}
</style>





