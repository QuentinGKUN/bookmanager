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

          <!-- 确认归还人信息 -->
          <div v-else-if="!borrowerConfirmed" class="borrower-section">
            <el-form :model="borrowerForm" :rules="borrowerRules" ref="borrowerFormRef" label-width="100px">
              <el-form-item label="电话" prop="phone">
                <el-input
                  v-model="borrowerForm.phone"
                  placeholder="请使用扫码枪扫描归还人电话或手动输入"
                  @keyup.enter="handleConfirmBorrower"
                  ref="phoneInput"
                />
              </el-form-item>
            </el-form>
            <el-button type="primary" size="large" @click="handleConfirmBorrower">确认归还人</el-button>

            <el-alert
              v-if="borrowerError"
              :title="borrowerError"
              type="error"
              :closable="false"
              style="margin-top: 20px"
            />
          </div>

          <!-- 归还图书 -->
          <div v-else>
            <!-- 归还人信息展示 -->
            <div class="borrower-info-section">
              <h3>归还人信息</h3>
              <el-descriptions :column="2" border>
                <el-descriptions-item label="姓名">{{ borrowerInfo.borrower_name }}</el-descriptions-item>
                <el-descriptions-item label="电话">{{ borrowerInfo.borrower_phone }}</el-descriptions-item>
                <el-descriptions-item label="借阅时间" :span="2">{{ formatTime(borrowerInfo.borrow_time) }}</el-descriptions-item>
              </el-descriptions>
              <el-button type="text" @click="handleResetBorrower" style="margin-top: 10px">更换归还人</el-button>
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
                <el-table-column prop="barcode" label="一维码" />
                <el-table-column prop="name" label="书名">
                  <template #default="scope">
                    {{ scope.row.name || scope.row.barcode }}
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
          </div>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import QRCode from 'qrcode'
import { borrowApi } from '../api/borrow'

const showQRCode = ref(true)
const qrcodeRef = ref(null)
const borrowerConfirmed = ref(false)
const borrowerFormRef = ref(null)
const phoneInput = ref(null)
const barcodeInput = ref(null)
const borrowerError = ref('')

const borrowerForm = reactive({
  phone: ''
})

const borrowerRules = {
  phone: [{ required: true, message: '请输入电话', trigger: 'blur' }]
}

const borrowerInfo = ref(null)
const barcode = ref('')
const returnedBooks = ref([])
const returnResult = ref(null)

const generateQRCode = async () => {
  if (qrcodeRef.value) {
    const url = window.location.href
    try {
      // 清空之前的内容
      qrcodeRef.value.innerHTML = ''
      await QRCode.toCanvas(qrcodeRef.value, url, {
        width: 200,
        margin: 2
      })
    } catch (error) {
      console.error('生成二维码失败', error)
      // 如果canvas失败，尝试使用img
      try {
        const dataUrl = await QRCode.toDataURL(url, { width: 200, margin: 2 })
        qrcodeRef.value.innerHTML = `<img src="${dataUrl}" alt="二维码" />`
      } catch (err) {
        console.error('生成二维码图片失败', err)
      }
    }
  }
}

const startReturn = () => {
  showQRCode.value = false
  nextTick(() => {
    phoneInput.value?.focus()
  })
}

const handleConfirmBorrower = async () => {
  try {
    await borrowerFormRef.value.validate()
    
    if (!borrowerForm.phone.trim()) {
      ElMessage.warning('请输入电话')
      return
    }

    borrowerError.value = ''
    const result = await borrowApi.getBorrowerByPhone({ phone: borrowerForm.phone.trim() })
    
    borrowerInfo.value = result
    borrowerConfirmed.value = true
    
    ElMessage.success('归还人信息确认成功')
    nextTick(() => {
      barcodeInput.value?.focus()
    })
  } catch (error) {
    borrowerError.value = error.message || '查询归还人信息失败'
    ElMessage.error(error.message || '查询归还人信息失败')
  }
}

const handleResetBorrower = () => {
  borrowerConfirmed.value = false
  borrowerInfo.value = null
  borrowerForm.phone = ''
  returnedBooks.value = []
  returnResult.value = null
  borrowerError.value = ''
  showQRCode.value = true
  nextTick(() => {
    generateQRCode()
  })
}

const handleReturn = async () => {
  if (!barcode.value.trim()) {
    ElMessage.warning('请输入一维码')
    return
  }

  if (!borrowerInfo.value) {
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
    await borrowApi.returnBook({
      barcode: barcode.value.trim(),
      borrower_phone: borrowerInfo.value.borrower_phone
    })
    
    // 添加到已归还列表
    const bookInfo = borrowerInfo.value.books.find(b => b.barcode === barcode.value.trim())
    returnedBooks.value.push({
      barcode: barcode.value.trim(),
      name: bookInfo?.name || null
    })
    
    returnResult.value = {
      message: '归还成功',
      type: 'success'
    }
    ElMessage.success('归还成功')
    barcode.value = ''
    nextTick(() => {
      barcodeInput.value?.focus()
    })
  } catch (error) {
    returnResult.value = {
      message: error.message || '归还失败',
      type: 'error'
    }
    ElMessage.error(error.message || '归还失败')
  }
}

const formatTime = (time) => {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleString('zh-CN')
}

onMounted(() => {
  generateQRCode()
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

.borrower-section {
  text-align: center;
  padding: 40px 0;
}

.borrower-section h3 {
  margin-bottom: 30px;
  font-size: 18px;
}

.borrower-info-section {
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
</style>


