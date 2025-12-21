<template>
  <div class="book-list">
    <el-container>
      <el-header>
        <div class="header-content">
          <h1>图书管理</h1>
          <el-button type="primary" @click="$router.push('/books/add')">新增图书</el-button>
        </div>
      </el-header>
      <el-main>
        <el-card>
          <el-form :inline="true" :model="searchForm" class="search-form">
            <el-form-item label="书名">
              <el-input v-model="searchForm.name" placeholder="请输入书名" clearable />
            </el-form-item>
            <el-form-item label="一维码">
              <el-input v-model="searchForm.barcode" placeholder="请输入一维码" clearable />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSearch">查询</el-button>
              <el-button @click="handleReset">重置</el-button>
            </el-form-item>
          </el-form>

          <el-table :data="tableData" border style="width: 100%">
            <el-table-column prop="barcode" label="一维码" width="150" />
            <el-table-column prop="name" label="书名" />
            <el-table-column prop="quantity" label="总数量" width="100" />
            <el-table-column prop="in_stock" label="在库数量" width="100" />
            <el-table-column prop="shelf_layer_name" label="位置" />
            <el-table-column prop="price" label="价格" width="100">
              <template #default="scope">
                {{ scope.row.price ? '¥' + scope.row.price : '-' }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="scope">
                <el-button link type="primary" @click="handleEdit(scope.row)">编辑</el-button>
                <el-button link type="danger" @click="handleDelete(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { bookApi } from '../api/book'
import { useRouter } from 'vue-router'

const router = useRouter()

const searchForm = ref({
  name: '',
  barcode: ''
})

const tableData = ref([])
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0
})

const loadData = async () => {
  try {
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.pageSize,
      ...searchForm.value
    }
    const res = await bookApi.list(params)
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
    name: '',
    barcode: ''
  }
  handleSearch()
}

const handleEdit = (row) => {
  router.push(`/books/edit/${row.id}`)
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该图书吗？', '提示', {
      type: 'warning'
    })
    await bookApi.delete(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
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
.book-list {
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

.search-form {
  margin-bottom: 20px;
}

.el-main {
  padding: 20px;
}
</style>



