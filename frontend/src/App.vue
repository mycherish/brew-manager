<script setup>
import { reactive, ref, computed, onMounted, onUnmounted } from 'vue' // 增加生命周期钩子
import { GetBrewData, StartService, StopService } from '../wailsjs/go/main/App'

const data = reactive({
  formulae: [],
  casks: [],
  loading: false
})

const searchQuery = ref('') // 搜索关键词 🔍

// 计算属性：过滤后的 Formulae
const filteredFormulae = computed(() => {
  return data.formulae.filter(item => 
    item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// 计算属性：过滤后的 Casks
const filteredCasks = computed(() => {
  return data.casks.filter(item =>
    item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})
// 为了管理每行的 loading，我们用一个 Map 来存储
const processingMap = reactive(new Map())

// 成功失败提示
const toast = reactive({
  show: false,
  msg: '',
  type: 'success'
})

function showToast(msg, type = 'success') {
  toast.msg = msg
  toast.type = type
  toast.show = true
  setTimeout(() => { toast.show = false }, 3000) // 3秒后消失
}

// 正在同步系统数据
async function updateList() {
  const res = await GetBrewData()
  data.formulae = res.formulae
  data.casks = res.casks
}
// 定时器
let timer = null
onMounted(() => {
  // 1. 进来先执行一次
  updateList()

  // 2. 开启每 10 秒一次的自动刷新
  timer = setInterval(() => {
    console.log("自动刷新中...")
    updateList()
  }, 10000)
})

onUnmounted(() => {
  // 3. 别忘了在组件卸载时清理定时器，防止内存泄漏
  if (timer) {
    clearInterval(timer)
  }
})

// 手动刷新
async function manualRefresh() {
  data.loading = true
  await updateList()
  data.loading = false
}

// 处理服务启动/停止
async function handleService(item) {
  // item.loading = true // 给单个项目加个加载状态，防止重复点击
  // 设置当前行正在处理中
  processingMap.set(item.name, true)
  
  try {
    let result;
    if (item.status === 'started') {
      result = await StopService(item.name)
    } else {
      result = await StartService(item.name)
    }
    // 根据结果给予反馈 (这里使用简单的 alert，或者你可以自定义一个 Toast 组件)
    if (result.success) {
      showToast("操作成功: " + result.message)
      // 成功后立即刷新列表
      await updateList()
    } else {
      // 失败弹出
      showToast(result.message, 'error')
    }
    
  } catch (err) {
    alert("系统错误: " + err)
  } finally {
    // 结束处理状态
    processingMap.delete(item.name)
  }
  
}

</script>

<template>
  <div class="container">
    <header class="drag-region">
      <h2>Brew-Manager</h2>
      <div class="toolbar">
        <button @click="manualRefresh" :disabled="data.loading">
          {{ data.loading ? '正在刷新...' : '手动刷新' }}
        </button>
        <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="输入软件名搜索..." 
            class="search-input"
          />
          <span class="refresh-tip">每 10s 自动同步状态</span>
      </div>
    </header>
    

    <div class="lists">
      <section>
        <h3>终端工具 ({{ filteredFormulae.length }})</h3>
        <div class="scroll-box">
          <div v-for="item in filteredFormulae" :key="item.name" class="item" :class="{ 'is-running': item.status === 'started', 'is-processing': processingMap.has(item.name) }">
              <div class="name-box">
                <span v-if="item.status !== 'none'" 
                      class="dot" 
                      :class="item.status === 'started' ? 'on' : 'off'">
                </span>
                <span class="name">{{ item.name }}</span>
            </div>
            <div class="action-box">
              <span class="version">{{ item.version }}</span>
              <button 
                v-if="item.status !== 'none_tool'" 
                @click="handleService(item)"
                class="btn-mini"
                :disabled="processingMap.has(item.name)"
                :class="item.status === 'started' ? 'btn-stop' : 'btn-start'"
              >
                <span v-if="processingMap.has(item.name)" class="loader"></span>
                <span v-else>{{ item.status === 'started' ? '停止' : '启动' }}</span>
                
              </button>
            </div>
          </div>
        </div>
      </section>

      <section>
        <h3>桌面应用 ({{ filteredCasks.length }})</h3>
        <div class="scroll-box">
          <div v-for="item in filteredCasks" :key="item.name" class="item" :class="{ 'is-running': item.status === 'started' }">
            <div class="name-box">
              <span v-if="item.status !== 'none'" 
                    class="dot" 
                    :class="item.status === 'started' ? 'on' : 'off'">
              </span>
              <span class="name">{{ item.name }}</span>
            </div>
            <span class="version">{{ item.version }}</span>
          </div>
        </div>
      </section>
    </div>
  </div>
  <transition name="fade">
    <div v-if="toast.show" class="toast" :class="toast.type">
      {{ toast.msg }}
    </div>
  </transition>
</template>

<style scoped>
body {
  margin: 0;
  /* 使用透明背景，这样 main.go 里的 WindowIsTranslucent 才会生效 */
  background-color: rgba(0, 0, 0, 0); 
}
/* 高亮运行中的行 */
.item.is-running {
  background-color: rgba(66, 185, 131, 0.1); /* 浅绿色背景 */
  border-left: 3px solid #42b983; /* 左侧绿色竖线 */
}
.name-box { display: flex; align-items: center; gap: 8px; }
.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}
.on { background-color: #42b983; box-shadow: 0 0 5px #42b983; }
.off { background-color: #ff4d4f; }
.name { color: #646cff; font-weight: bold; }
.container { 
  /* 加上一点半透明背景，配合毛玻璃效果更好看 */
  background-color: rgba(30, 30, 30, 0.6); 
  backdrop-filter: blur(20px); /* 额外的毛玻璃加成 */
  height: 100vh;
  border-radius: 10px; /* 让窗口圆角更明显 */
  overflow: hidden;
  display: flex;
  flex-direction: column;
  padding: 0 20px 20px 20px;
}
.lists { display: flex; gap: 20px; margin-top: 20px; justify-content: center; }
.scroll-box { 
  height: 500px; 
  overflow-y: auto; 
  width: 320px;
  border: 1px solid #444; 
  background: rgba(26, 26, 26, 0.4);; 
  padding: 5px;

  /* 必须：排除掉不需要拖拽的交互元素 */
  /* 如果不加 no-drag，你的搜索框将无法选中，按钮也将无法点击 */
  /* --wails-draggable: no-drag !important; */
}
.item { 
  display: flex; justify-content: space-between; 
  padding: 8px; border-bottom: 1px solid #333; font-size: 14px;
  transition: all 0.3s ease; /* 增加一点平滑过渡 */
}
.name { color: #646cff; font-weight: bold; }
.version { color: #888; font-family: monospace; }
section h3 { color: #fff; border-bottom: 2px solid #646cff; padding-bottom: 5px; }
button { 
  cursor: pointer; 
  padding: 10px 20px; 
  font-weight: bold; 
  /* --wails-draggable: no-drag !important; */
}
/* 2. 定义拖拽区 */
.toolbar {
  display: flex;
  gap: 15px;
  justify-content: center;
  align-items: center;
  margin-bottom: 20px;
  padding: 10px 0 20px 0;
  background: rgba(255, 255, 255, 0.05);
  /* --wails-draggable: drag; */
}
.search-input {
  padding: 8px 15px;
  border-radius: 20px;
  border: 1px solid #444;
  background: #2a2a2a;
  color: white;
  width: 250px;
  outline: none;
  /* --wails-draggable: no-drag !important; */
}

.search-input:focus {
  border-color: #646cff;
}
.action-box { display: flex; align-items: center; gap: 10px; }
.btn-mini {
  padding: 4px 8px;
  font-size: 12px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  /* --wails-draggable: no-drag !important; */
}
.btn-start { background: #42b983; color: white; }
.btn-stop { background: #ff4d4f; color: white; }
.btn-mini:hover { opacity: 0.8; }
.refresh-tip {
  font-size: 12px;
  color: #666;
  margin-left: 10px;
}
.loader {
  width: 12px;
  height: 12px;
  border: 2px solid #FFF;
  border-bottom-color: transparent;
  border-radius: 50%;
  display: inline-block;
  animation: rotation 1s linear infinite;
}
@keyframes rotation {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 处理中的行变淡 */
.is-processing {
  opacity: 0.6;
  pointer-events: none; /* 防止过程中再次点击 */
}

/* 按钮点击反馈 */
.btn-mini:active {
  transform: scale(0.95);
}

.btn-mini {
  transition: all 0.2s ease;
  min-width: 60px; /* 固定宽度防止文字变化时按钮抖动 */
}
/* 1. 全局背景透明，让系统的毛玻璃效果透出来 */
#app {
  background-color: transparent !important;
  height: 100vh;
}
/* 1. 定义顶部的拖拽大区 */
.drag-region {
  --wails-draggable: drag; /* 开启拖拽 */
  user-select: none;       /* 禁止选中文字，防止干扰拖拽 */
  padding-top: 32px;       /* 为左上角红绿灯留出空间 */
  width: 100%;
  -webkit-user-select: none;
  cursor: default;
}
/* 2. 标题和文字通常不需要交互，保持默认 */
h2 {
  margin: 0;
  padding: 10px 0;
  text-align: center;
}
/* 5. 精确排除：只有这些真正要点的地方才不让拖拽 */
.search-input, 
button, 
.btn-mini,
.scroll-box,
.refresh-tip {
  --wails-draggable: no-drag !important;
}
.toast {
  position: fixed;
  top: 50px;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 20px;
  border-radius: 8px;
  color: white;
  z-index: 9999;
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 15px rgba(0,0,0,0.3);
}
.toast.success { background: rgba(66, 185, 131, 0.9); }
.toast.error { background: rgba(255, 77, 79, 0.9); }

.fade-enter-active, .fade-leave-active { transition: opacity 0.5s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>