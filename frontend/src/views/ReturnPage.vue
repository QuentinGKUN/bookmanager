<template>
  <div class="return-page">
    <el-container>
      <el-header>
        <h1>图书归还</h1>
      </el-header>
      <el-main>
        <el-card>
          <!-- 二维码展示 -->
          <div class="qrcode-section" v-if="showQRCode">
            <h2>请扫描二维码进行归还</h2>
            <div ref="qrcodeRef" style="display: flex; justify-content: center; margin: 20px 0;"></div>
            <el-button type="primary" @click="startReturn">开始归还</el-button>
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

          <!-- 归还操作 -->
          <div v-else>
            <!-- 用户信息展示 -->
            <div class="user-info-section">
              <h3>归还人信息</h3>
              <el-descriptions :column="2" border>
                <el-descriptions-item label="姓名">{{ currentUser.name }}</el-descriptions-item>
                <el-descriptions-item label="电话">{{ currentUser.phone }}</el-descriptions-item>
              </el-descriptions>
              <el-button type="text" @click="handleResetUser" style="margin-top: 10px">更换归还人</el-button>
            </div>

            <!-- 扫码归还区域 -->
            <div class="scan-section">
              <h3>请扫描图书一维码进行归还</h3>
              <el-input
                v-model="barcode"
                placeholder="请使用扫码枪扫描图书一维码"
                @keyup.enter="handleReturn"
                ref="barcodeInput"
                style="margin-bottom: 20px"
              />
              <el-button type="primary" size="large" @click="handleReturn">归还</el-button>
            </div>

            <!-- 已归还图书列表 -->
            <div class="book-list-section" v-if="returnedBooks.length > 0">
              <h3>已归还图书（{{ returnedBooks.length }}本）</h3>
              <el-table :data="returnedBooks" border>
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

            <el-alert
              v-if="returnResult"
              :title="returnResult.message"
              :type="returnResult.type"
              :closable="false"
              style="margin-top: 20px"
            />

            <!-- 操作按钮 -->
            <div class="action-section">
              <el-button type="success" size="large" @click="handleComplete" :disabled="returnedBooks.length === 0">完成归还</el-button>
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
const barcodeInput = ref(null)
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

const barcode = ref('')
const returnedBooks = ref([])
const returnResult = ref(null)
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

const startReturn = async () => {
  showQRCode.value = false
  // 检查是否有未完成的归还
  await loadReturnData()
  if (currentUser.value) {
    userConfirmed.value = true
    nextTick(() => {
      barcodeInput.value?.focus()
    })
  } else {
    nextTick(() => {
      userFormRef.value?.clearValidate()
    })
  }
}

const loadReturnData = async () => {
  try {
    const result = await borrowApi.getReturnUser()
    if (result.user) {
      currentUser.value = result.user
      returnedBooks.value = result.books || []
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
      await borrowApi.setReturnUser({
        name: userForm.name.trim(),
        phone: userForm.phone.trim()
      })
      
      currentUser.value = {
        name: userForm.name.trim(),
        phone: userForm.phone.trim()
      }
      userConfirmed.value = true
      returnedBooks.value = []
      returnResult.value = null
      
      ElMessage.success('用户信息已保存')
      nextTick(() => {
        barcodeInput.value?.focus()
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

const handleReturn = async () => {
  if (!barcode.value.trim()) {
    ElMessage.warning('请输入一维码')
    return
  }

  if (!currentUser.value) {
    ElMessage.warning('请先确认归还人信息')
    return
  }

  // 检查是否已归还
  if (returnedBooks.value.some(b => b.barcode === barcode.value.trim())) {
    ElMessage.warning('该图书已归还')
    barcode.value = ''
    barcodeInput.value?.focus()
    return
  }

  try {
    const result = await borrowApi.addReturnBook({ barcode: barcode.value.trim() })
    await loadReturnData() // 重新加载列表
    ElMessage.success('添加成功')
    barcode.value = ''
    nextTick(() => {
      barcodeInput.value?.focus()
    })
  } catch (error) {
    ElMessage.error(error.message || '添加失败')
    barcode.value = ''
    barcodeInput.value?.focus()
  }
}

const handleRemoveBook = async (index) => {
  try {
    await borrowApi.removeReturnBook({ index })
    await loadReturnData() // 重新加载列表
    ElMessage.success('删除成功')
  } catch (error) {
    ElMessage.error(error.message || '删除失败')
  }
}

const handleComplete = async () => {
  if (returnedBooks.value.length === 0) {
    ElMessage.warning('请至少添加一本图书')
    return
  }

  try {
    await borrowApi.completeReturn({ use_redis: true })
    returnResult.value = {
      message: '归还成功',
      type: 'success'
    }
    ElMessage.success('归还成功')
    handleReset()
    showQRCode.value = true
    nextTick(() => {
      generateQRCode()
    })
  } catch (error) {
    returnResult.value = {
      message: error.message || '归还失败',
      type: 'error'
    }
    ElMessage.error(error.message || '归还失败')
  }
}

const handleReset = () => {
  userForm.name = ''
  userForm.phone = ''
  barcode.value = ''
  returnedBooks.value = []
  currentUser.value = null
  userConfirmed.value = false
  returnResult.value = null
  userFormRef.value?.resetFields()
  stopRefreshTimer()
}

const handleResetUser = () => {
  userConfirmed.value = false
  currentUser.value = null
  returnedBooks.value = []
  returnResult.value = null
  stopRefreshTimer()
}

const startRefreshTimer = () => {
  stopRefreshTimer()
  refreshTimer = setInterval(async () => {
    await loadReturnData()
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
.return-page {
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
  text-align: center;
}

.scan-section h3 {
  margin-bottom: 30px;
  font-size: 18px;
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
