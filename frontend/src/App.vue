<script setup>
import { reactive, ref, computed, onMounted, onUnmounted } from 'vue'
import { GetBrewData, StartService, StopService } from '../wailsjs/go/main/App'

const data = reactive({
  formulae: [],
  casks: [],
  loading: false
})

const searchQuery = ref('')
const processingMap = reactive(new Map())
const toast = reactive({ show: false, msg: '', type: 'success' })

// 统一的提示函数
function showToast(msg, type = 'success') {
  toast.msg = msg
  toast.type = type
  toast.show = true
  setTimeout(() => { toast.show = false }, 3000)
}

// 搜索过滤逻辑
const filteredFormulae = computed(() => {
  return data.formulae.filter(item => item.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
})

const filteredCasks = computed(() => {
  return data.casks.filter(item => item.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
})

// 数据刷新逻辑
async function updateList() {
  try {
    const res = await GetBrewData()
    data.formulae = res.formulae
    data.casks = res.casks
  } catch (err) {
    console.error("刷新失败:", err)
  }
}

// 手动刷新按钮逻辑
async function manualRefresh() {
  data.loading = true
  await updateList()
  data.loading = false
  showToast("数据已同步")
}

// 定时器管理
let timer = null
onMounted(() => {
  updateList()
  timer = setInterval(() => {
    updateList()
  }, 10000)
})

onUnmounted(() => { if (timer) clearInterval(timer) })

// 服务启动/停止逻辑
async function handleService(item) {
  processingMap.set(item.name, true)
  try {
    let result = item.status === 'started' ? await StopService(item.name) : await StartService(item.name)
    if (result.success) {
      showToast(result.message, 'success')
      await updateList()
    } else {
      showToast(result.message, 'error')
    }
  } catch (err) {
    showToast("系统错误: " + err, 'error')
  } finally {
    processingMap.delete(item.name)
  }
}
</script>

<template>
  <div class="container">
    <header class="drag-region">
      <div class="header-content">
        <div class="title-group">
          <h2>Brew Manager</h2>
          <span class="sync-tag">Auto Sync ON</span>
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

    <div class="main-content">
      <div class="lists-wrapper">
        <section class="list-column">
          <div class="column-header">
            <h3>TERMINAL TOOLS</h3>
            <span class="count-badge">{{ filteredFormulae.length }}</span>
          </div>
          <div class="scroll-area">
            <div v-for="item in filteredFormulae" :key="item.name" 
                 class="item-card" :class="{ 'is-processing': processingMap.has(item.name) }">
              <div class="item-main">
                <span v-if="item.status !== 'none_tool'" 
                      class="status-indicator" 
                      :class="item.status === 'started' ? 'online' : 'offline'">
                </span>
                <div class="info-meta">
                  <span class="name">{{ item.name }}</span>
                  <span class="version">{{ item.version }}</span>
                </div>
              </div>
              
              <div class="item-actions" v-if="item.status !== 'none_tool'">
                <button @click="handleService(item)" 
                        class="mac-btn"
                        :class="item.status === 'started' ? 'btn-stop' : 'btn-start'"
                        :disabled="processingMap.has(item.name)">
                  {{ processingMap.has(item.name) ? '...' : (item.status === 'started' ? '停止' : '启动') }}
                </button>
              </div>
            </div>
          </div>
        </section>

        <section class="list-column">
          <div class="column-header">
            <h3>GUI APPLICATIONS</h3>
            <span class="count-badge">{{ filteredCasks.length }}</span>
          </div>
          <div class="scroll-area">
            <div v-for="item in filteredCasks" :key="item.name" class="item-card">
              <div class="item-main">
                <div class="info-meta">
                  <span class="name">{{ item.name }}</span>
                  <span class="version text-dim">{{ item.version }}</span>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>

    <transition name="toast-slide">
      <div v-if="toast.show" class="toast-notification" :class="toast.type">
        <div class="toast-content">
          <span class="toast-icon">{{ toast.type === 'success' ? '✅' : '❌' }}</span>
          <span class="toast-msg">{{ toast.msg }}</span>
        </div>
      </div>
    </transition>
  </div>
</template>

<style scoped>
/* 1. 基础容器：毛玻璃与全屏控制 */
.container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: rgba(28, 28, 30, 0.7); /* macOS 暗色面板色 */
  backdrop-filter: blur(40px) saturate(180%);
  color: #fff;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue", sans-serif;
  overflow: hidden;
}

/* 2. 顶部拖拽区与导航 */
.drag-region {
  --wails-draggable: drag;
  padding: 45px 24px 15px; /* 顶部留出红绿灯空间 */
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.title-group h2 {
  margin: 0;
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.5px;
}

.sync-tag {
  font-size: 10px;
  color: #34C759;
  text-transform: uppercase;
  font-weight: 700;
  margin-top: 4px;
  display: block;
}

/* 3. 工具栏：搜索与按钮 */
.toolbar {
  display: flex;
  gap: 12px;
  --wails-draggable: no-drag;
}

.search-box {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  padding: 6px 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.search-input {
  background: transparent;
  border: none;
  color: #fff;
  font-size: 13px;
  outline: none;
  width: 160px;
}

.btn-refresh {
  background: #007AFF; /* macOS 蓝色按钮 */
  border: none;
  color: white;
  padding: 6px 16px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
}

/* 4. 列表区域布局 */
.main-content {
  flex: 1;
  overflow: hidden;
  padding: 10px 20px 20px;
}

.lists-wrapper {
  display: flex;
  gap: 20px;
  height: 100%;
}

.list-column {
  flex: 1;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.05);
  display: flex;
  flex-direction: column;
}

.column-header {
  padding: 15px 15px 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.column-header h3 {
  font-size: 11px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.4);
  margin: 0;
  letter-spacing: 0.5px;
}

.count-badge {
  background: rgba(255, 255, 255, 0.1);
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 10px;
}

.scroll-area {
  flex: 1;
  overflow-y: auto;
  padding: 0 8px 15px;
}

/* 5. 列表项卡片设计 */
.item-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  margin-bottom: 2px;
  border-radius: 8px;
  transition: background 0.2s ease;
}

.item-card:hover {
  background: rgba(255, 255, 255, 0.06);
}

.item-main {
  display: flex;
  align-items: center;
  gap: 12px;
}

.info-meta {
  display: flex;
  flex-direction: column;
}

.name { font-size: 14px; font-weight: 500; color: #fff; }
.version { font-size: 11px; color: #888; font-family: 'SF Mono', Menlo, monospace; }

/* 6. 状态指示灯 */
.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}
.online { background: #34C759; box-shadow: 0 0 8px rgba(52, 199, 89, 0.4); }
.offline { background: #FF3B30; opacity: 0.5; }

/* 7. 操作按钮 (macOS 风格) */
.mac-btn {
  padding: 5px 12px;
  border-radius: 5px;
  font-size: 12px;
  font-weight: 500;
  border: 0.5px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.1);
  color: #eee;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-start:hover { background: #34C759; border-color: #34C759; color: white; }
.btn-stop:hover { background: #FF3B30; border-color: #FF3B30; color: white; }

/* 8. 右上角 Toast 通知 */
.toast-notification {
  position: fixed;
  top: 40px;
  right: 20px;
  min-width: 240px;
  background: rgba(40, 40, 40, 0.8);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.3);
  padding: 14px;
  z-index: 9999;
  --wails-draggable: no-drag !important;
}

.toast-content { display: flex; align-items: center; gap: 10px; }
.toast-icon { font-size: 16px; }
.toast-msg { font-size: 13px; font-weight: 500; color: #fff; }

.toast-slide-enter-from, .toast-slide-leave-to { opacity: 0; transform: translateX(30px); }
.toast-slide-enter-active, .toast-slide-leave-active { transition: all 0.4s cubic-bezier(0.23, 1, 0.32, 1); }

/* 处理中状态 */
.is-processing { opacity: 0.5; pointer-events: none; }

/* 滚动条美化 */
.scroll-area::-webkit-scrollbar { width: 5px; }
.scroll-area::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.1); border-radius: 10px; }
</style>