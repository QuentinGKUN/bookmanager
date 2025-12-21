<template>
  <div class="borrow-query">
    <el-container>
      <el-header>
        <h1>借阅记录查询</h1>
      </el-header>
      <el-main>
        <el-card>
          <el-form :inline="true" :model="searchForm" class="search-form">
            <el-form-item label="姓名">
              <el-input v-model="searchForm.borrower_name" placeholder="请输入姓名" clearable />
            </el-form-item>
            <el-form-item label="电话">
              <el-input v-model="searchForm.borrower_phone" placeholder="请输入电话" clearable />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSearch">查询</el-button>
              <el-button @click="handleReset">重置</el-button>
            </el-form-item>
          </el-form>

          <el-table :data="tableData" border style="width: 100%">
            <el-table-column prop="borrower_name" label="借阅人姓名" />
            <el-table-column prop="borrower_phone" label="借阅人电话" />
            <el-table-column prop="borrow_time" label="借阅时间" width="180">
              <template #default="scope">
                {{ formatTime(scope.row.borrow_time) }}
              </template>
            </el-table-column>
            <el-table-column label="借阅图书" min-width="300">
              <template #default="scope">
                <div v-for="book in scope.row.books" :key="book.barcode" style="margin-bottom: 5px">
                  {{ book.name || book.barcode }}
                </div>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="pagination.total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handlePageChange"
            style="margin-top: 20px; justify-content: flex-end"
          />
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { borrowApi } from '../api/borrow'

const searchForm = ref({
  borrower_name: '',
  borrower_phone: ''
})

const tableData = ref([])
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0
})

const formatTime = (time) => {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleString('zh-CN')
}

const loadData = async () => {
  try {
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.pageSize,
      ...searchForm.value
    }
    const res = await borrowApi.list(params)
    tableData.value = res.list
    pagination.value.total = res.total
  } catch (error) {
    ElMessage.error(error.message || '加载失败')
  }
}

const handleSearch = () => {
  pagination.value.page = 1
  loadData()
}

const handleReset = () => {
  searchForm.value = {
    borrower_name: '',
    borrower_phone: ''
  }
  handleSearch()
}

const handleSizeChange = () => {
  loadData()
}

const handlePageChange = () => {
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.borrow-query {
  min-height: 100vh;
  background: #f5f5f5;
}

.el-header h1 {
  margin: 0;
  line-height: 60px;
  font-size: 24px;
}

.search-form {
  margin-bottom: 20px;
}

.el-main {
  padding: 20px;
}
</style>



