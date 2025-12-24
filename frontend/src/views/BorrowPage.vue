<template>
  <div class="borrow-page">
    <el-container>
      <el-header>
        <h1>图书借阅</h1>
      </el-header>
      <el-main>
        <el-card>
          <!-- 二维码展示 -->
          <div class="qrcode-section" v-if="showQRCode">
            <h2>请扫描二维码进行借阅</h2>
            <div ref="qrcodeRef" style="display: flex; justify-content: center; margin: 20px 0;"></div>
            <el-button type="primary" @click="startBorrow">开始借阅</el-button>
          </div>

          <!-- 用户信息输入 -->
          <div v-else-if="!userConfirmed" class="user-section">
            <h3>请输入您的信息</h3>
            <el-form :model="userForm" :rules="userRules" ref="userFormRef" label-width="100px">
              <el-form-item label="姓名" prop="name">
                <el-input v-model="userForm.name" placeholder="请输入姓名" />
              </el-form-item>
              <el-form-item label="电话" prop="phone">
                <el-input v-model="userForm.phone" placeholder="请输入电话" />
              </el-form-item>
            </el-form>
            <el-button type="primary" size="large" @click="handleConfirmUser">确定</el-button>
          </div>

          <!-- 借阅操作 -->
          <div v-else>
            <!-- 用户信息展示 -->
            <div class="user-info-section">
              <h3>借阅人信息</h3>
              <el-descriptions :column="2" border>
                <el-descriptions-item label="姓名">{{ currentUser.name }}</el-descriptions-item>
                <el-descriptions-item label="电话">{{ currentUser.phone }}</el-descriptions-item>
              </el-descriptions>
              <el-button type="text" @click="handleResetUser" style="margin-top: 10px">更换借阅人</el-button>
            </div>

            <!-- 扫码区域 -->
            <div class="scan-section">
              <h3>请扫描图书一维码</h3>
              <el-input
                v-model="scanBarcode"
                placeholder="请使用扫码枪扫描图书一维码"
                @keyup.enter="handleScan"
                ref="scanInput"
                style="margin-bottom: 20px"
              />
              <el-button type="primary" @click="handleScan">手动添加</el-button>
            </div>

            <!-- 已借阅图书列表 -->
            <div class="book-list-section" v-if="borrowedBooks.length > 0">
              <h3>已借阅图书（{{ borrowedBooks.length }}本）</h3>
              <el-table :data="borrowedBooks" border>
                <el-table-column type="index" label="序号" width="60" />
                <el-table-column prop="barcode" label="一维码" />
                <el-table-column prop="name" label="书名">
                  <template #default="scope">
                    {{ scope.row.name || scope.row.barcode }}
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="100">
                  <template #default="scope">
                    <el-button type="danger" size="small" @click="handleRemoveBook(scope.$index)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>

            <!-- 操作按钮 -->
            <div class="action-section">
              <el-button type="success" size="large" @click="handleComplete" :disabled="borrowedBooks.length === 0">完成借阅</el-button>
              <el-button size="large" @click="handleReset">重置</el-button>
            </div>
          </div>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import QRCode from 'qrcode'
import { borrowApi } from '../api/borrow'

const showQRCode = ref(true)
const qrcodeRef = ref(null)
const userFormRef = ref(null)
const scanInput = ref(null)
const userConfirmed = ref(false)
const currentUser = ref(null)

const userForm = reactive({
  name: '',
  phone: ''
})

const userRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入电话', trigger: 'blur' }]
}

const scanBarcode = ref('')
const borrowedBooks = ref([])
let refreshTimer = null

const generateQRCode = async () => {
  if (qrcodeRef.value) {
    const url = window.location.href
    try {
      qrcodeRef.value.innerHTML = ''
      await QRCode.toCanvas(qrcodeRef.value, url, {
        width: 200,
        margin: 2
      })
    } catch (error) {
      console.error('生成二维码失败', error)
      try {
        const dataUrl = await QRCode.toDataURL(url, { width: 200, margin: 2 })
        qrcodeRef.value.innerHTML = `<img src="${dataUrl}" alt="二维码" />`
      } catch (err) {
        console.error('生成二维码图片失败', err)
      }
    }
  }
}

const startBorrow = async () => {
  showQRCode.value = false
  // 检查是否有未完成的借阅
  await loadBorrowData()
  if (currentUser.value) {
    userConfirmed.value = true
    nextTick(() => {
      scanInput.value?.focus()
    })
  } else {
    nextTick(() => {
      userFormRef.value?.clearValidate()
    })
  }
}

const loadBorrowData = async () => {
  try {
    const result = await borrowApi.getBorrowUser()
    if (result.user) {
      currentUser.value = result.user
      borrowedBooks.value = result.books || []
      userConfirmed.value = true
    }
  } catch (error) {
    // 忽略错误，可能是Redis中没有数据
  }
}

const handleConfirmUser = async () => {
  try {
    await userFormRef.value.validate()
    
    try {
      await borrowApi.setBorrowUser({
        name: userForm.name.trim(),
        phone: userForm.phone.trim()
      })
      
      currentUser.value = {
        name: userForm.name.trim(),
        phone: userForm.phone.trim()
      }
      userConfirmed.value = true
      borrowedBooks.value = []
      
      ElMessage.success('用户信息已保存')
      nextTick(() => {
        scanInput.value?.focus()
      })
      
      // 开始定时刷新
      startRefreshTimer()
    } catch (error) {
      ElMessage.error(error.message || '保存用户信息失败')
    }
  } catch (error) {
    // 表单验证失败
  }
}

const handleScan = async () => {
  if (!scanBarcode.value.trim()) {
    ElMessage.warning('请输入一维码')
    return
  }

  const barcode = scanBarcode.value.trim()
  
  // 检查是否已添加
  if (borrowedBooks.value.some(b => b.barcode === barcode)) {
    ElMessage.warning('该图书已添加')
    scanBarcode.value = ''
    scanInput.value?.focus()
    return
  }

  try {
    const result = await borrowApi.addBorrowBook({ barcode })
    await loadBorrowData() // 重新加载列表
    ElMessage.success('添加成功')
    scanBarcode.value = ''
    scanInput.value?.focus()
  } catch (error) {
    ElMessage.error(error.message || '添加失败')
    scanBarcode.value = ''
    scanInput.value?.focus()
  }
}

const handleRemoveBook = async (index) => {
  try {
    await borrowApi.removeBorrowBook({ index })
    await loadBorrowData() // 重新加载列表
    ElMessage.success('删除成功')
  } catch (error) {
    ElMessage.error(error.message || '删除失败')
  }
}

const handleComplete = async () => {
  if (borrowedBooks.value.length === 0) {
    ElMessage.warning('请至少添加一本图书')
    return
  }

  try {
    await borrowApi.completeBorrow({ use_redis: true })
    ElMessage.success('借阅成功')
    handleReset()
    showQRCode.value = true
    nextTick(() => {
      generateQRCode()
    })
  } catch (error) {
    ElMessage.error(error.message || '借阅失败')
  }
}

const handleReset = () => {
  userForm.name = ''
  userForm.phone = ''
  scanBarcode.value = ''
  borrowedBooks.value = []
  currentUser.value = null
  userConfirmed.value = false
  userFormRef.value?.resetFields()
  stopRefreshTimer()
}

const handleResetUser = () => {
  userConfirmed.value = false
  currentUser.value = null
  borrowedBooks.value = []
  stopRefreshTimer()
}

const startRefreshTimer = () => {
  stopRefreshTimer()
  refreshTimer = setInterval(async () => {
    await loadBorrowData()
  }, 2000) // 每2秒刷新一次
}

const stopRefreshTimer = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

onMounted(() => {
  generateQRCode()
})

onUnmounted(() => {
  stopRefreshTimer()
})
</script>

<style scoped>
.borrow-page {
  min-height: 100vh;
  background: #f5f5f5;
}

.el-header h1 {
  margin: 0;
  line-height: 60px;
  font-size: 24px;
  text-align: center;
}

.el-main {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.qrcode-section {
  text-align: center;
  padding: 40px 0;
}

.user-section {
  padding: 40px 0;
  text-align: center;
}

.user-section h3 {
  margin-bottom: 30px;
  font-size: 18px;
}

.user-info-section {
  margin-bottom: 30px;
  padding: 20px;
  background: #f9f9f9;
  border-radius: 4px;
}

.scan-section {
  margin: 30px 0;
  padding: 20px;
  background: #f9f9f9;
  border-radius: 4px;
}

.book-list-section {
  margin: 30px 0;
}

.action-section {
  margin-top: 30px;
  text-align: center;
}

.action-section .el-button {
  margin: 0 10px;
  min-width: 120px;
}
</style>
