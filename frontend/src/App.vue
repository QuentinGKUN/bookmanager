<template>
  <div id="app">
    <NavMenu v-if="showNav" />
    <router-view />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import NavMenu from './components/NavMenu.vue'

const route = useRoute()
// 移动端借阅页面不显示导航
const showNav = computed(() => {
  return route.path !== '/borrow'
})
</script>

<script setup>
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

#app {
  min-height: 100vh;
}
</style>

