<template>
  <el-menu
    :default-active="activeIndex"
    mode="horizontal"
    router
    style="border-bottom: 1px solid #e6e6e6"
  >
    <el-menu-item index="/books">图书管理</el-menu-item>
    <el-menu-item index="/locations">位置管理</el-menu-item>
    <el-menu-item index="/borrow">借阅</el-menu-item>
    <el-menu-item index="/borrow/query">借阅查询</el-menu-item>
    <el-menu-item index="/return">归还</el-menu-item>
  </el-menu>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const activeIndex = computed(() => route.path)
</script>

<style scoped>
</style>



