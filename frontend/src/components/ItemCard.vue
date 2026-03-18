<template>
  <div class="item-card" :class="{ 'is-processing': isProcessing }">
    <div class="item-main">
      <span v-if="type === 'formula' && item.status !== 'none_tool'"
            class="status-indicator"
            :class="item.status === 'started' ? 'online' : 'offline'">
      </span>
      <div v-if="type === 'cask'" class="app-icon-container">
          <img v-if="item.iconBase64"
               :src="'data:image/png;base64,' + item.iconBase64"
               class="app-icon" />
          <span v-else class="app-placeholder">📦</span>
        </div>
      <div class="info-meta">
        <span class="name">{{ item.name }}</span>
        <span class="version">{{ item.version }}</span>
      </div>
    </div>

    <div class="item-actions" v-if="type === 'formula' && item.status !== 'none_tool'">
      <button
          class="mac-btn btn-restart"
          @click.stop="$emit('restart', item.name)"
          :disabled="isProcessing"
          title="重启服务"
      >
        <svg class="restart-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 12a9 9 0 11-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/>
          <path d="M21 3v5h-5"/>
        </svg>
      </button>

      <button @click="$emit('action', item)"
              class="mac-btn"
              :class="item.status === 'started' ? 'btn-stop' : 'btn-start'"
              :disabled="isProcessing">
        <span v-if="isProcessing" class="mini-loader"></span>
        <span>{{ item.status === 'started' ? '停止' : '启动' }}</span>
      </button>
    </div>
  </div>
</template>

<script setup>
/**
 * 定义 Props
 * @param item - 软件对象数据
 * @param type - 类型：'formula' 或 'cask'
 * @param isProcessing - 当前项是否正在处理中（由外部 processingMap 决定）
 */
defineProps({
  item: { type: Object, required: true },
  type: { type: String, default: 'formula' },
  isProcessing: { type: Boolean, default: false }
})

/**
 * 定义 Emit
 * 当点击按钮时，通知父组件执行具体的 Start/Stop 逻辑
 */
defineEmits(['action', 'restart'])
</script>

<style scoped>
/* 这里放置 item-card 相关的专用 CSS */
/* 动画和基础样式已经在 main.css 中，这里可以放一些细微的间距调整 */
.info-meta {
  display: flex;
  flex-direction: column; /* 强制垂直排列 */
  align-items: flex-start; /* 左对齐 */
  gap: 2px; /* 名字和版本号之间留一点小缝隙 */
}

/* 6. 状态指示灯 */
.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}
.online { background: #34C759; box-shadow: 0 0 8px rgba(52, 199, 89, 0.4); }
.offline { background: #FF3B30; opacity: 0.5; }
.name { font-size: 14px; font-weight: 500; color: #fff; }
.version { font-size: 11px; color: #888; font-family: 'SF Mono', Menlo, monospace; }
.item-main {
  display: flex;
  align-items: center;
  gap: 12px;
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
.item-actions {
  display: flex;
  flex-direction: row; /* 强制横向 */
  align-items: center;
  gap: 6px; /* 按钮之间的间距 */
  flex-shrink: 0; /* 防止操作区被名字挤压 */
}
.btn-restart {
  width: 28px; /* 方形按钮 */
  padding: 0;
  color: #0A84FF;
  background: rgba(10, 132, 255, 0.15);
  border-color: rgba(10, 132, 255, 0.2);
}
.btn-restart:hover {
  background: #0A84FF;
  color: white;
}
.restart-icon {
  display: inline-block;
  width: 14px;
  height: 14px;
}

/* 4. [优化] 针对不同的服务状态，你可以微调背景色 */
.btn-start {
  min-width: 50px;
  color: #30D158;
  background: rgba(48, 209, 88, 0.15);
  border-color: rgba(48, 209, 88, 0.2);
}
.btn-start:hover {
  background: #238A4A;
  color: white;
}
/* 停止按钮：深红色背景 */
.btn-stop {
  min-width: 50px;
  color: #FF453A;
  background: rgba(255, 69, 58, 0.15);
  border-color: rgba(255, 69, 58, 0.2);
}
.btn-stop:hover {
  background: #A62A24;
  color: white;
}
/* 旋转动画 */
.restart-icon {
  width: 14px;
  height: 14px;
}
/* 处理中状态 */
.is-processing {
  opacity: 0.6; /* 只让卡片变淡，不影响按钮区域的占位 */
}
.is-processing .restart-icon {
  animation: spin 1s linear infinite;
  pointer-events: none; /* 防止重复点击 */
}
.app-icon-container {
  width: 32px;
  height: 32px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.app-icon {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>