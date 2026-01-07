<template>
  <div class="book-form">
    <el-container>
      <el-header>
        <div class="header-content">
          <h1>{{ isEdit ? '编辑图书' : '新增图书' }}</h1>
          <el-button @click="$router.back()">返回</el-button>
        </div>
      </el-header>
      <el-main>
        <el-card>
          <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
            <el-form-item label="一维码" prop="barcode">
              <el-input
                v-model="form.barcode"
                placeholder="请使用扫码枪扫描一维码"
                :disabled="isEdit"
                @keyup.enter="handleBarcodeEnter"
                ref="barcodeInput"
              />
            </el-form-item>
            <el-form-item label="书名" prop="name">
              <el-input v-model="form.name" placeholder="请输入书名" ref="nameInput" />
            </el-form-item>
            <el-form-item label="数量" prop="quantity">
              <el-input-number v-model="form.quantity" :min="0" style="width: 100%" />
            </el-form-item>
            <el-form-item label="在库数量">
              <el-input-number v-model="form.in_stock" :min="0" style="width: 100%" />
              <div class="form-tip">不填时默认与数量相同</div>
            </el-form-item>
            <el-form-item label="位置">
              <LocationSelect v-model="form.shelf_layer_id" />
            </el-form-item>
            <el-form-item label="价格">
              <el-input-number v-model="form.price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
            <el-form-item label="备注">
              <el-input v-model="form.remark" type="textarea" :rows="4" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSubmit">保存</el-button>
              <el-button @click="$router.back()">取消</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { bookApi } from '../api/book'
import LocationSelect from '../components/LocationSelect.vue'

const route = useRoute()
const router = useRouter()

const isEdit = ref(false)
const formRef = ref(null)
const barcodeInput = ref(null)
const nameInput = ref(null)

const form = reactive({
  barcode: '',
  name: '',
  quantity: 0,
  in_stock: null,
  shelf_layer_id: null,
  price: null,
  remark: ''
})

const rules = {
  barcode: [{ required: true, message: '请输入一维码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入书名', trigger: 'blur' }],
  quantity: [{ required: true, message: '请输入数量', trigger: 'blur' }]
}

const handleBarcodeEnter = () => {
  // 扫码后自动聚焦到书名输入框
  nextTick(() => {
    nameInput.value?.focus()
  })
}

const loadBook = async () => {
  if (route.params.id) {
    isEdit.value = true
    try {
      const res = await bookApi.list({ barcode: '' })
      const book = res.list.find(b => b.id === parseInt(route.params.id))
      if (book) {
        Object.assign(form, {
          barcode: book.barcode,
          name: book.name,
          quantity: book.quantity,
          in_stock: book.in_stock,
          shelf_layer_id: book.shelf_layer_id,
          price: book.price,
          remark: book.remark || ''
        })
      }
    } catch (error) {
      ElMessage.error('加载失败')
    }
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    const data = { ...form }
    if (data.in_stock === null) {
      data.in_stock = data.quantity
    }
    
    if (isEdit.value) {
      await bookApi.update(route.params.id, data)
      ElMessage.success('更新成功')
    } else {
      await bookApi.create(data)
      ElMessage.success('创建成功')
    }
    router.push('/books')
  } catch (error) {
    if (error !== false) {
      ElMessage.error(error.message || '保存失败')
    }
  }
}

onMounted(() => {
  loadBook()
  // 自动聚焦到一维码输入框
  nextTick(() => {
    barcodeInput.value?.focus()
  })
})
</script>

<style scoped>
.book-form {
  min-height: 100vh;
  background: #f5f5f5;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.header-content h1 {
  margin: 0;
  font-size: 24px;
}

.el-main {
  padding: 20px;
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}
</style>





