<script setup>
/**
 * Overview 组件 - 概览页面
 * 显示系统统计数据和最近安装的软件
 */
import { computed } from 'vue'

const props = defineProps({
  // Homebrew 数据
  data: {
    type: Object,
    required: true
  },
  // 是否正在处理中
  isProcessing: {
    type: Boolean,
    default: false
  }
})

// 计算统计数据
const stats = computed(() => ({
  formulae: props.data.formulae?.length || 0,
  casks: props.data.casks?.length || 0,
  taps: props.data.taps?.length || 0,
  total: (props.data.formulae?.length || 0) + (props.data.casks?.length || 0)
}))

// 最近安装的软件（取前5个）
const recentItems = computed(() => {
  const all = [
    ...(props.data.formulae || []).map(f => ({ ...f, type: 'formula' })),
    ...(props.data.casks || []).map(c => ({ ...c, type: 'cask' }))
  ]
  return all.slice(0, 5)
})
</script>

<template>
  <div class="overview-page">
    <h2 class="page-title">概览</h2>
    
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">📦</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.total }}</div>
          <div class="stat-label">已安装软件</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">💻</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.formulae }}</div>
          <div class="stat-label">命令行工具</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">🖥️</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.casks }}</div>
          <div class="stat-label">GUI 应用</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">🔗</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.taps }}</div>
          <div class="stat-label">Taps 源</div>
        </div>
      </div>
    </div>

    <!-- 最近安装 -->
    <div class="recent-section" v-if="recentItems.length > 0">
      <h3>已安装软件</h3>
      <div class="recent-list">
        <div v-for="item in recentItems" :key="item.name" class="recent-item">
          <span class="item-type" :class="item.type">{{ item.type === 'formula' ? '💻' : '🖥️' }}</span>
          <span class="item-name">{{ item.name }}</span>
          <span class="item-version">{{ item.version }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overview-page {
  padding: 20px;
  color: #fff;
}

.page-title {
  margin: 0 0 24px 0;
  font-size: 1.5rem;
  font-weight: 600;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 16px;
  margin-bottom: 32px;
}

.stat-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.2s ease;
}

.stat-card:hover {
  background: rgba(255, 255, 255, 0.08);
  transform: translateY(-2px);
}

.stat-icon {
  font-size: 2rem;
  opacity: 0.9;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.5);
}

.quick-actions {
  margin-bottom: 32px;
}

.quick-actions h3 {
  margin: 0 0 16px 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
}

.action-buttons {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: rgba(0, 122, 255, 0.2);
  border: 1px solid rgba(0, 122, 255, 0.3);
  border-radius: 8px;
  color: #fff;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn:hover:not(:disabled) {
  background: rgba(0, 122, 255, 0.3);
  transform: translateY(-1px);
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-icon {
  font-size: 1rem;
}

.recent-section h3 {
  margin: 0 0 16px 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
}

.recent-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.recent-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  transition: all 0.2s ease;
}

.recent-item:hover {
  background: rgba(255, 255, 255, 0.06);
}

.item-type {
  font-size: 1.125rem;
}

.item-name {
  flex: 1;
  font-weight: 500;
  font-family: 'Monaco', 'Menlo', monospace;
}

.item-version {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.5);
}
</style>
