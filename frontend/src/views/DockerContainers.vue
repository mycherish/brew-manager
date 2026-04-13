<script setup>
import { ref, computed, onMounted } from 'vue'
import { GetDockerContainers, StartDockerContainer, StopDockerContainer } from '../../wailsjs/go/main/App'

const containers = ref([])
const loading = ref(true)
const error = ref('')
const processingMap = ref(new Map())
const localSearch = ref('')

const filteredContainers = computed(() => {
  if (!localSearch.value) return containers.value
  const query = localSearch.value.toLowerCase()
  return containers.value.filter(c => 
    c.name.toLowerCase().includes(query) || 
    c.image.toLowerCase().includes(query)
  )
})

const runningCount = computed(() => {
  return containers.value.filter(c => c.state === 'running').length
})

async function loadContainers() {
  loading.value = true
  error.value = ''
  try {
    const result = await GetDockerContainers()
    if (!result.success) {
      error.value = result.message
      containers.value = []
    } else {
      containers.value = result.data || []
    }
  } catch (e) {
    error.value = '无法连接 Docker，请确认 Docker Desktop 已启动'
  } finally {
    loading.value = false
  }
}

async function startContainer(container) {
  const key = container.id
  processingMap.value.set(key, true)
  try {
    const result = await StartDockerContainer(container.id)
    if (result.success) {
      await loadContainers()
    } else {
      alert(result.message)
    }
  } finally {
    processingMap.value.delete(key)
  }
}

async function stopContainer(container) {
  const key = container.id
  processingMap.value.set(key, true)
  try {
    const result = await StopDockerContainer(container.id)
    if (result.success) {
      await loadContainers()
    } else {
      alert(result.message)
    }
  } finally {
    processingMap.value.delete(key)
  }
}

function isProcessing(id) {
  return processingMap.value.has(id)
}

onMounted(() => {
  loadContainers()
})
</script>

<template>
  <div class="docker-page">
    <div class="page-header">
      <h2 class="page-title">Docker 容器</h2>
      <div class="header-stats">
        <span class="stat-badge running">{{ runningCount }} 运行中</span>
        <span class="stat-badge total">{{ containers.length }} 总计</span>
      </div>
    </div>
    
    <!-- 本地搜索 -->
    <div class="local-search">
      <span class="search-icon">🔍</span>
      <input 
        v-model="localSearch" 
        type="text" 
        placeholder="搜索容器名称或镜像..." 
        class="search-input"
      />
      <button class="refresh-btn" @click="loadContainers" :disabled="loading">
        {{ loading ? '⟳' : '↻' }}
      </button>
    </div>

    <!-- 加载中 -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载容器列表...</p>
    </div>

    <!-- 错误提示 -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <p>{{ error }}</p>
      <button class="retry-btn" @click="loadContainers">重试</button>
    </div>

    <!-- 容器列表 -->
    <div v-else class="containers-list">
      <div 
        v-for="container in filteredContainers" 
        :key="container.id" 
        class="container-card"
        :class="{ running: container.state === 'running' }"
      >
        <div class="container-status-dot" :class="container.state"></div>
        
        <div class="container-info">
          <div class="container-name">{{ container.name }}</div>
          <div class="container-meta">
            <span class="meta-item">
              <span class="meta-label">镜像:</span>
              {{ container.image }}
            </span>
            <span class="meta-item" v-if="container.ports">
              <span class="meta-label">端口:</span>
              {{ container.ports }}
            </span>
          </div>
          <div class="container-status">{{ container.status }}</div>
        </div>
        
        <div class="container-actions">
          <button 
            v-if="container.state === 'running'"
            class="action-btn stop"
            :disabled="isProcessing(container.id)"
            @click="stopContainer(container)"
          >
            <span v-if="isProcessing(container.id)" class="btn-spinner"></span>
            <span v-else>⏹</span>
            停止
          </button>
          <button 
            v-else
            class="action-btn start"
            :disabled="isProcessing(container.id)"
            @click="startContainer(container)"
          >
            <span v-if="isProcessing(container.id)" class="btn-spinner"></span>
            <span v-else>▶</span>
            启动
          </button>
        </div>
      </div>
      
      <div v-if="filteredContainers.length === 0 && !loading" class="empty-state">
        <div class="empty-icon">🐳</div>
        <p v-if="containers.length === 0">暂无 Docker 容器</p>
        <p v-else>未找到匹配的容器</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.docker-page {
  padding: 20px;
  color: #fff;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.page-title {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
}

.header-stats {
  display: flex;
  gap: 8px;
}

.stat-badge {
  background: rgba(255, 255, 255, 0.1);
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.7);
}

.stat-badge.running {
  background: rgba(52, 199, 89, 0.2);
  color: #34C759;
}

.local-search {
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 10px 16px;
  margin-bottom: 20px;
}

.search-icon {
  font-size: 1rem;
  opacity: 0.6;
}

.search-input {
  flex: 1;
  background: transparent;
  border: none;
  color: #fff;
  font-size: 0.875rem;
  outline: none;
}

.search-input::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.refresh-btn {
  background: transparent;
  border: none;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  font-size: 1rem;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.2s;
}

.refresh-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.loading-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.5);
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(255, 255, 255, 0.1);
  border-top-color: #0A84FF;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.5);
}

.error-icon {
  font-size: 3rem;
  margin-bottom: 16px;
}

.retry-btn {
  margin-top: 16px;
  padding: 8px 20px;
  background: #0A84FF;
  border: none;
  border-radius: 6px;
  color: #fff;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s;
}

.retry-btn:hover {
  background: #0070E0;
}

.containers-list {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.containers-list::-webkit-scrollbar {
  width: 5px;
}

.containers-list::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
}

.container-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  padding: 14px 16px;
  transition: all 0.2s;
}

.container-card:hover {
  background: rgba(255, 255, 255, 0.08);
}

.container-card.running {
  border-left: 3px solid #34C759;
}

.container-status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.container-status-dot.running {
  background: #34C759;
  box-shadow: 0 0 8px rgba(52, 199, 89, 0.5);
}

.container-status-dot.exited {
  background: #8E8E93;
}

.container-status-dot.paused {
  background: #FF9500;
}

.container-info {
  flex: 1;
  min-width: 0;
}

.container-name {
  font-size: 0.95rem;
  font-weight: 500;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.container-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 0.8rem;
  color: rgba(255, 255, 255, 0.6);
  margin-bottom: 4px;
}

.meta-item {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.meta-label {
  color: rgba(255, 255, 255, 0.4);
  margin-right: 4px;
}

.container-status {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.5);
}

.container-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border: none;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn.start {
  background: rgba(52, 199, 89, 0.2);
  color: #34C759;
}

.action-btn.start:hover:not(:disabled) {
  background: rgba(52, 199, 89, 0.3);
}

.action-btn.stop {
  background: rgba(255, 59, 48, 0.2);
  color: #FF3B30;
}

.action-btn.stop:hover:not(:disabled) {
  background: rgba(255, 59, 48, 0.3);
}

.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-spinner {
  width: 12px;
  height: 12px;
  border: 2px solid transparent;
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: rgba(255, 255, 255, 0.5);
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state p {
  margin: 0;
  font-size: 0.875rem;
}
</style>
