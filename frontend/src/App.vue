<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useBrew } from './composables/useBrew'
import Sidebar from './components/Sidebar.vue'
import Toast from './components/Toast.vue'
import Overview from './views/Overview.vue'
import GuiApps from './views/GuiApps.vue'
import TuiTools from './views/TuiTools.vue'
import TapsPage from './views/TapsPage.vue'

const {
  data, currentTab, setCurrentTab, stats,
  processingMap, toast,
  refreshProgress, // 自动刷新进度
  updateList, refreshCurrentTab,
  handleService, handleRestart,
  showToast,
  handleAddTap, handleRemoveTap, handleUpdateTap, handleUpdateAllTaps,
  startAutoRefresh, stopAutoRefresh
} = useBrew()

// 切换 tab
function handleChangeTab(tabId) {
  setCurrentTab(tabId)
  refreshCurrentTab()
}

// 手动刷新
function handleRefresh() {
  refreshCurrentTab()
}

// 启动时更新数据 + 启动定时刷新
onMounted(async () => {
  data.loading = true
  await updateList()
  data.loading = false
  // 启动定时刷新
  startAutoRefresh()
})

// 组件卸载时停止定时器
onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<template>
  <div class="app-container">
    <!-- 左侧导航 -->
    <Sidebar 
      :current-tab="currentTab"
      :stats="stats"
      :is-loading="data.loading"
      :refresh-progress="refreshProgress"
      @change-tab="handleChangeTab"
      @refresh="handleRefresh"
    />
    
    <!-- 主内容区域 -->
    <main class="main-content">
      <!-- 加载状态 -->
      <div v-if="data.loading && !data.formulae.length" class="loading-overlay">
        <div class="loader-content">
          <div class="spinner"></div>
          <p>正在同步 Homebrew 数据...</p>
        </div>
      </div>
      
      <!-- 页面内容 -->
      <div v-else class="page-content">
        <!-- 概览页 -->
        <Overview 
          v-if="currentTab === 'overview'"
          :data="data"
          :is-processing="data.loading"
        />
        
        <!-- GUI 应用页 -->
        <GuiApps 
          v-else-if="currentTab === 'gui'"
          :casks="data.casks"
          :processing-map="processingMap"
          @action="handleService"
          @restart="handleRestart"
        />
        
        <!-- TUI 工具页 -->
        <TuiTools 
          v-else-if="currentTab === 'tui'"
          :formulae="data.formulae"
          :processing-map="processingMap"
          @action="handleService"
          @restart="handleRestart"
        />
        
        <!-- Taps 管理页 -->
        <TapsPage 
          v-else-if="currentTab === 'taps'"
          :taps="data.taps"
          :is-processing="data.loading"
          :processing-map="processingMap"
          @add="handleAddTap"
          @remove="handleRemoveTap"
          @update="handleUpdateTap"
          @update-all="handleUpdateAllTaps"
        />
      </div>
    </main>
    
    <!-- Toast 提示 -->
    <Toast v-if="toast.show" :msg="toast.msg" :type="toast.type" />
  </div>
</template>

<style>
@import "./assets/main.css";

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue", sans-serif;
}

.app-container {
  height: 100vh;
  display: flex;
  background: rgba(28, 28, 30, 0.95);
  color: #fff;
  overflow: hidden;
}

.main-content {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.page-content {
  height: 100%;
  overflow-y: auto;
}

.page-content::-webkit-scrollbar {
  width: 6px;
}

.page-content::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(28, 28, 30, 0.95);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.loader-content {
  text-align: center;
  color: rgba(255, 255, 255, 0.7);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(255, 255, 255, 0.1);
  border-top-color: #007AFF;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loader-content p {
  font-size: 0.875rem;
}
</style>
