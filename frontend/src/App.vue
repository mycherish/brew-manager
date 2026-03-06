<script setup>
import { onMounted } from 'vue'
import { useBrew } from './composables/useBrew'
import Toast from './components/Toast.vue'
import ItemCard from './components/ItemCard.vue'

const { 
  data, searchQuery, processingMap, toast,
  filteredFormulae, filteredCasks, 
  updateList, handleService,
  showToast
} = useBrew()

// 手动刷新按钮逻辑
async function manualRefresh() {
  data.loading = true
  await updateList()
  data.loading = false
  showToast("数据已同步")
}

// 启动时更新数据
onMounted(async () => { // 注意这里加了 async
  data.loading = true
  await updateList()     // 必须 await，否则 loading 效果一闪而过
  data.loading = false
})

</script>

<template>
  <div class="container">
    <header class="drag-region">
      <div class="header-content">
        <div class="title-group">
          <h2>Brew Manager</h2>
        </div>
        
        <div class="toolbar">
          <div class="search-box">
            <span class="search-icon">🔍</span>
            <input v-model="searchQuery" type="text" placeholder="搜索组件..." class="search-input" />
          </div>
          <button @click="manualRefresh" class="btn-refresh" :disabled="data.loading">
            <span v-if="data.loading" class="mini-loader"></span>
            <span v-else>刷新</span>
          </button>
          
        </div>
      </div>
    </header>

    <Transition name="fade">
      <div v-if="data.loading" class="loading-overlay">
        <div class="loader-content">
          <div class="spinner-container">
            <div class="spinner-modern"></div>
          </div>
          <p class="loading-text">正在同步 Homebrew 数据...</p>
          <div class="loading-bar-placeholder">
            <div class="loading-bar-inner"></div>
          </div>
        </div>
      </div>
    </Transition>
    

    <div class="main-content">
      <div class="lists-wrapper">
        <section class="list-column">
          <div class="column-header">
            <h3>Formulae</h3>
            <span class="count-badge">{{ filteredFormulae.length }}</span>
          </div>
          <div class="scroll-area">
            <ItemCard 
              v-for="item in filteredFormulae" 
              :key="item.name"
              type="formula"
              :item="item"
              :is-processing="processingMap.has(item.name)"
              @action="handleService" 
            />
          </div>
        </section>

        <section class="list-column">
          <div class="column-header">
            <h3>Casks</h3>
            <span class="count-badge">{{ filteredCasks.length }}</span>
          </div>
          <div class="scroll-area">
            <ItemCard 
              v-for="item in filteredCasks" 
              :key="item.name"
              type="cask"
              :item="item"
            />
          </div>
        </section>
      </div>
    </div>

    <Toast v-if="toast.show" :msg="toast.msg" :type="toast.type" />
  </div>
</template>

<style scoped>
@import "././assets/main.css"
</style>