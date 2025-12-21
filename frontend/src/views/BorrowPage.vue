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

          <!-- 借阅表单 -->
          <div v-else>
            <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
              <el-form-item label="姓名" prop="borrower_name">
                <el-input v-model="form.borrower_name" placeholder="请输入姓名" />
              </el-form-item>
              <el-form-item label="电话" prop="borrower_phone">
                <el-input v-model="form.borrower_phone" placeholder="请输入电话" />
              </el-form-item>
            </el-form>

            <!-- 扫码区域 -->
            <div class="scan-section">
              <h3>扫码借阅</h3>
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
                <el-table-column prop="barcode" label="一维码" />
                <el-table-column prop="name" label="书名">
                  <template #default="scope">
                    {{ scope.row.name || scope.row.barcode }}
                  </template>
                </el-table-column>
              </el-table>
            </div>

            <!-- 操作按钮 -->
            <div class="action-section">
              <el-button type="success" size="large" @click="handleSubmit" :disabled="borrowedBooks.length === 0">完成借阅</el-button>
              <el-button size="large" @click="handleReset">重置</el-button>
            </div>
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
import { onBeforeUnmount } from 'vue'
import { borrowApi } from '../api/borrow'

const showQRCode = ref(true)
const qrcodeRef = ref(null)
const formRef = ref(null)
const scanInput = ref(null)

const form = reactive({
  borrower_name: '',
  borrower_phone: ''
})

const rules = {
  borrower_name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  borrower_phone: [{ required: true, message: '请输入电话', trigger: 'blur' }]
}

const scanBarcode = ref('')
const borrowedBooks = ref([])

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

const startBorrow = () => {
  showQRCode.value = false
  nextTick(() => {
    scanInput.value?.focus()
  })
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
    // 调用扫码接口（不验证图书是否存在）
    const result = await borrowApi.scan({ barcode })
    
    borrowedBooks.value.push({
      barcode: result.barcode,
      name: result.name || null
    })
    
    ElMessage.success('添加成功')
    scanBarcode.value = ''
    scanInput.value?.focus()
  } catch (error) {
    // 即使图书不存在也添加到列表
    borrowedBooks.value.push({
      barcode: barcode,
      name: null
    })
    ElMessage.success('已添加')
    scanBarcode.value = ''
    scanInput.value?.focus()
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    
    if (borrowedBooks.value.length === 0) {
      ElMessage.warning('请至少添加一本图书')
      return
    }

    const barcodes = borrowedBooks.value.map(b => b.barcode)
    await borrowApi.create({
      borrower_name: form.borrower_name,
      borrower_phone: form.borrower_phone,
      barcodes: barcodes
    })

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
  form.borrower_name = ''
  form.borrower_phone = ''
  scanBarcode.value = ''
  borrowedBooks.value = []
  formRef.value?.resetFields()
}

onMounted(() => {
  generateQRCode()
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

