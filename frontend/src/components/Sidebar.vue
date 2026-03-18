<script setup>
/**
 * Sidebar 组件
 * 左侧导航栏，包含 tab 切换、统计信息和刷新按钮
 */
const props = defineProps({
  // 当前选中的 tab
  currentTab: {
    type: String,
    required: true
  },
  // 统计数据
  stats: {
    type: Object,
    default: () => ({ gui: 0, tui: 0, taps: 0 })
  },
  // 是否正在加载
  isLoading: {
    type: Boolean,
    default: false
  },
  // 自动刷新进度状态
  refreshProgress: {
    type: Object,
    default: () => ({ isRefreshing: false, progress: 0, lastRefreshTime: null })
  }
})

const emit = defineEmits(['change-tab', 'refresh'])

// Tab 配置
const tabs = [
  { id: 'overview', label: '概览', icon: '🏠' },
  { id: 'gui', label: 'GUI 应用', icon: '📦' },
  { id: 'tui', label: '命令行工具', icon: '💻' },
  { id: 'taps', label: 'Taps', icon: '🔗' }
]

/**
 * 处理 tab 切换
 * @param {string} tabId - tab ID
 */
function handleTabClick(tabId) {
  emit('change-tab', tabId)
}

/**
 * 处理手动刷新
 */
function handleRefresh() {
  emit('refresh')
}

/**
 * 格式化剩余时间
 * @param {number} progress - 进度百分比
 * @returns {string} 格式化后的时间字符串
 */
function formatTimeLeft(progress) {
  if (!props.refreshProgress.lastRefreshTime || progress >= 100) {
    return '即将刷新...'
  }
  const totalSeconds = 120 // 2分钟
  const elapsed = Math.floor((progress / 100) * totalSeconds)
  const remaining = totalSeconds - elapsed
  const minutes = Math.floor(remaining / 60)
  const seconds = remaining % 60
  return `${minutes}:${seconds.toString().padStart(2, '0')} 后刷新`
}
</script>

<template>
  <aside class="sidebar">
    <!-- 顶部标题 -->
    <div class="sidebar-header">
      <h1 class="app-title">Brew Manager</h1>
    </div>
    
    <!-- 导航列表 -->
    <nav class="sidebar-nav">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        class="nav-item"
        :class="{ active: currentTab === tab.id }"
        @click="handleTabClick(tab.id)"
      >
        <span class="nav-icon">{{ tab.icon }}</span>
        <span class="nav-label">{{ tab.label }}</span>
        <span v-if="tab.id === 'gui'" class="nav-badge">{{ stats.gui }}</span>
        <span v-if="tab.id === 'tui'" class="nav-badge">{{ stats.tui }}</span>
        <span v-if="tab.id === 'taps'" class="nav-badge">{{ stats.taps }}</span>
      </button>
    </nav>
    
    <!-- 底部区域 -->
    <div class="sidebar-footer">
      <!-- 刷新按钮和进度条 -->
      <div class="refresh-section">
        <button 
          class="refresh-btn" 
          @click="handleRefresh" 
          :disabled="isLoading || refreshProgress.isRefreshing"
        >
          <svg v-if="refreshProgress.isRefreshing" class="refresh-icon spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12a9 9 0 11-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/>
            <path d="M21 3v5h-5"/>
          </svg>
          <svg v-else class="refresh-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12a9 9 0 11-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/>
            <path d="M21 3v5h-5"/>
          </svg>
          <span class="refresh-label">{{ refreshProgress.isRefreshing ? '刷新中...' : '刷新' }}</span>
        </button>
        
        <!-- 进度条 -->
        <div class="progress-bar">
          <div 
            class="progress-fill" 
            :style="{ width: refreshProgress.progress + '%' }"
          ></div>
        </div>
        <div class="progress-text">{{ formatTimeLeft(refreshProgress.progress) }}</div>
      </div>
      
      <!-- 版本信息 -->
      <div class="app-version">v1.3.0</div>
    </div>
  </aside>
</template>

<style scoped>
/**
 * 侧边栏容器
 */
.sidebar {
  width: 200px;
  background: rgba(0, 0, 0, 0.3);
  border-right: 1px solid rgba(255, 255, 255, 0.05);
  display: flex;
  flex-direction: column;
  padding: 20px 0;
}

/**
 * 顶部标题区域
 */
.sidebar-header {
  padding: 0 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  margin-bottom: 16px;
}

.app-title {
  margin: 0;
  font-size: 1rem;
  font-weight: 700;
  color: #fff;
  letter-spacing: -0.3px;
}

/**
 * 导航列表
 */
.sidebar-nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 0 12px;
}

/**
 * 导航项
 */
.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  background: transparent;
  border: none;
  border-radius: 8px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.9);
}

.nav-item.active {
  background: rgba(0, 122, 255, 0.2);
  color: #007AFF;
}

/**
 * 导航图标
 */
.nav-icon {
  font-size: 1.125rem;
  width: 24px;
  text-align: center;
}

/**
 * 导航标签
 */
.nav-label {
  flex: 1;
}

/**
 * 导航徽章（统计数量）
 */
.nav-badge {
  background: rgba(255, 255, 255, 0.1);
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.6);
}

.nav-item.active .nav-badge {
  background: rgba(0, 122, 255, 0.3);
  color: #007AFF;
}

/**
 * 底部区域
 */
.sidebar-footer {
  padding: 16px 12px 0;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  margin-top: 16px;
}

/**
 * 刷新区域
 */
.refresh-section {
  margin-bottom: 12px;
}

/**
 * 刷新按钮
 */
.refresh-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 10px 12px;
  background: rgba(0, 122, 255, 0.15);
  border: 1px solid rgba(0, 122, 255, 0.3);
  border-radius: 8px;
  color: #007AFF;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.refresh-btn:hover:not(:disabled) {
  background: rgba(0, 122, 255, 0.25);
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.refresh-label {
  flex: 1;
  text-align: center;
}

/**
 * 刷新图标
 */
.refresh-icon {
  width: 16px;
  height: 16px;
}

.refresh-icon.spin {
  animation: spin 1s linear infinite;
}

/**
 * 进度条容器
 */
.progress-bar {
  height: 4px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
  margin-top: 10px;
  overflow: hidden;
}

/**
 * 进度条填充
 */
.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #007AFF, #00C7BE);
  border-radius: 2px;
  transition: width 1s linear;
}

/**
 * 进度文字
 */
.progress-text {
  font-size: 0.7rem;
  color: rgba(255, 255, 255, 0.4);
  text-align: center;
  margin-top: 6px;
}

/**
 * 版本信息
 */
.app-version {
  text-align: center;
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.4);
  padding: 10px 0;
}

/**
 * 旋转动画
 */
@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
