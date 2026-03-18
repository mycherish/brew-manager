<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { BrowserOpenURL } from '../../wailsjs/runtime'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  },
  installedFormulae: {
    type: Array,
    default: () => []
  },
  installedCasks: {
    type: Array,
    default: () => []
  },
  installedTaps: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['close', 'action', 'restart'])

const searchQuery = ref('')
const hasSearched = ref(false)

// 本地搜索：只搜索已安装的包
const filteredFormulae = computed(() => {
  if (!searchQuery.value.trim()) return []
  const query = searchQuery.value.toLowerCase()
  return props.installedFormulae.filter(f => 
    f.name.toLowerCase().includes(query)
  )
})

const filteredCasks = computed(() => {
  if (!searchQuery.value.trim()) return []
  const query = searchQuery.value.toLowerCase()
  return props.installedCasks.filter(c => 
    c.name.toLowerCase().includes(query)
  )
})

const filteredTaps = computed(() => {
  if (!searchQuery.value.trim()) return []
  const query = searchQuery.value.toLowerCase()
  return props.installedTaps.filter(t => 
    t.name.toLowerCase().includes(query)
  )
})

const totalResults = computed(() => 
  filteredFormulae.value.length + 
  filteredCasks.value.length + 
  filteredTaps.value.length
)

function handleSearch() {
  if (!searchQuery.value.trim()) return
  hasSearched.value = true
}

function handleAction(item) {
  emit('action', item)
}

function handleRestart(name) {
  emit('restart', name)
}

function handleClose() {
  emit('close')
  searchQuery.value = ''
  hasSearched.value = false
}

function handleKeydown(e) {
  if (e.key === 'Escape') {
    handleClose()
  }
  if (e.key === 'Enter' && searchQuery.value.trim()) {
    handleSearch()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <Transition name="modal">
    <div v-if="isOpen" class="modal-overlay" @click="handleClose">
      <div class="modal-content" @click.stop>
        <!-- 搜索头部 -->
        <div class="search-header">
          <span class="search-icon">🔍</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索已安装的软件包..."
            class="search-input"
            ref="searchInput"
            autofocus
          />
          <button v-if="searchQuery" @click="searchQuery = ''" class="clear-btn">×</button>
        </div>

        <!-- 搜索结果 -->
        <div v-if="hasSearched" class="search-results">
          <!-- 空结果 -->
          <div v-if="totalResults === 0" class="empty-result">
            <div class="empty-icon">🔍</div>
            <p>未找到匹配的软件包</p>
            <p class="hint">尝试使用其他关键词搜索</p>
          </div>

          <!-- 有结果 -->
          <div v-else class="results-content">
            <div class="results-summary">
              找到 {{ totalResults }} 个已安装的软件
            </div>

            <!-- Formulae -->
            <div v-if="filteredFormulae.length > 0" class="result-section">
              <h4 class="section-title">
                <span class="badge badge-formula">Formula</span>
                命令行工具 ({{ filteredFormulae.length }})
              </h4>
              <div class="package-list">
                <div 
                  v-for="pkg in filteredFormulae" 
                  :key="pkg.name"
                  class="package-item installed"
                >
                  <div class="package-info">
                    <span class="package-name">{{ pkg.name }}</span>
                    <span class="package-version">{{ pkg.version }}</span>
                  </div>
                  <div class="package-actions">
                    <button 
                      v-if="pkg.status === 'started'"
                      @click="handleAction(pkg)"
                      class="btn-action"
                    >
                      停止
                    </button>
                    <button 
                      v-else-if="pkg.status === 'stopped'"
                      @click="handleAction(pkg)"
                      class="btn-action"
                    >
                      启动
                    </button>
                    <span v-else class="status-badge">{{ pkg.status === 'none' ? '非服务' : pkg.status }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Casks -->
            <div v-if="filteredCasks.length > 0" class="result-section">
              <h4 class="section-title">
                <span class="badge badge-cask">Cask</span>
                GUI 应用 ({{ filteredCasks.length }})
              </h4>
              <div class="package-list">
                <div 
                  v-for="pkg in filteredCasks" 
                  :key="pkg.name"
                  class="package-item installed"
                >
                  <div class="package-info">
                    <span class="package-name">{{ pkg.name }}</span>
                    <span class="package-version">{{ pkg.version }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Taps -->
            <div v-if="filteredTaps.length > 0" class="result-section">
              <h4 class="section-title">
                <span class="badge badge-tap">Tap</span>
                Taps ({{ filteredTaps.length }})
              </h4>
              <div class="package-list">
                <div 
                  v-for="tap in filteredTaps" 
                  :key="tap.name"
                  class="package-item tap-item"
                >
                  <div class="package-info tap-info">
                    <div class="tap-name-row">
                      <span class="package-name">{{ tap.name }}</span>
                      <span v-if="tap.official" class="badge badge-official-mini">官方</span>
                    </div>
                    <button 
                      v-if="tap.url"
                      class="tap-link"
                      @click.stop="BrowserOpenURL(tap.url)"
                    >
                      🔗 {{ tap.url.replace('https://github.com/', '') }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 初始状态提示 -->
        <div v-else class="search-hint">
          <div class="hint-icon">⌨️</div>
          <p>输入关键词搜索已安装的软件包</p>
          <div class="shortcuts">
            <span class="shortcut"><kbd>Esc</kbd> 关闭</span>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 80px;
  z-index: 1000;
}

.modal-content {
  background: #1c1c1e;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  width: 90%;
  max-width: 600px;
  max-height: 70vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}

.search-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.search-icon {
  font-size: 1.25rem;
  opacity: 0.6;
}

.search-input {
  flex: 1;
  background: transparent;
  border: none;
  color: #fff;
  font-size: 1rem;
  outline: none;
}

.search-input::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.clear-btn {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 50%;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  font-size: 1rem;
  line-height: 1;
}

.clear-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
}

.search-btn {
  padding: 8px 16px;
  background: #007AFF;
  border: none;
  border-radius: 6px;
  color: #fff;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 60px;
}

.search-btn:hover:not(:disabled) {
  background: #0056CC;
}

.search-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.search-results {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
}

.search-results::-webkit-scrollbar {
  width: 6px;
}

.search-results::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
}

.empty-result {
  text-align: center;
  padding: 40px 20px;
  color: rgba(255, 255, 255, 0.5);
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 16px;
  opacity: 0.5;
}

.hint {
  font-size: 0.875rem;
  opacity: 0.7;
  margin: 0;
}

.results-summary {
  margin-bottom: 16px;
  padding: 10px 14px;
  background: rgba(0, 122, 255, 0.1);
  border-radius: 8px;
  color: #007AFF;
  font-weight: 500;
  font-size: 0.875rem;
}

.result-section {
  margin-bottom: 20px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0 0 10px 0;
  font-size: 0.875rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.7);
}

.section-title.installed {
  color: #2ECC71;
}

.badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.625rem;
  font-weight: 600;
  text-transform: uppercase;
}

.badge-installed {
  background: rgba(46, 204, 113, 0.2);
  color: #2ECC71;
}

.badge-formula {
  background: rgba(52, 152, 219, 0.2);
  color: #3498DB;
}

.badge-cask {
  background: rgba(155, 89, 182, 0.2);
  color: #9B59B6;
}

.badge-tap {
  background: rgba(241, 196, 15, 0.2);
  color: #F1C40F;
}

.package-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.package-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  transition: all 0.2s ease;
}

.package-item:hover {
  background: rgba(255, 255, 255, 0.06);
}

.package-item.installed {
  border-left: 3px solid #2ECC71;
}

.package-info {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.tap-info {
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}

.package-name {
  font-weight: 500;
  font-family: 'Monaco', 'Menlo', monospace;
  color: #fff;
}

.package-type {
  font-size: 0.625rem;
  color: rgba(255, 255, 255, 0.4);
  text-transform: uppercase;
  padding: 2px 6px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
}

.package-status {
  font-size: 0.75rem;
  color: #2ECC71;
}

.package-version {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.5);
}

.status-badge {
  font-size: 0.625rem;
  color: rgba(255, 255, 255, 0.4);
  padding: 2px 8px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
}

.badge-official-mini {
  background: rgba(0, 122, 255, 0.2);
  color: #007AFF;
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 0.625rem;
  font-weight: 600;
}

.package-actions {
  display: flex;
  gap: 6px;
}

.btn-action {
  padding: 6px 12px;
  background: rgba(0, 122, 255, 0.2);
  border: 1px solid rgba(0, 122, 255, 0.3);
  border-radius: 6px;
  color: #007AFF;
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-action:hover {
  background: rgba(0, 122, 255, 0.3);
}

.btn-install {
  padding: 6px 14px;
  background: #2ECC71;
  border: none;
  border-radius: 6px;
  color: #fff;
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 60px;
}

.btn-install:hover:not(:disabled) {
  background: #27AE60;
}

.btn-install:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-installed {
  padding: 6px 14px;
  background: rgba(46, 204, 113, 0.1);
  border: 1px solid rgba(46, 204, 113, 0.2);
  border-radius: 6px;
  color: #2ECC71;
  font-size: 0.75rem;
  font-weight: 500;
  cursor: default;
}

.btn-add-tap {
  padding: 6px 14px;
  background: #F39C12;
  border: none;
  border-radius: 6px;
  color: #fff;
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-add-tap:hover {
  background: #E67E22;
}

.tap-stars {
  font-size: 0.75rem;
  color: #F1C40F;
}

.tap-description {
  margin: 0 0 4px 0;
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.5);
  line-height: 1.4;
}

.tap-name-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.tap-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.7rem;
  color: #007AFF;
  background: transparent;
  border: none;
  padding: 0;
  cursor: pointer;
  opacity: 0.8;
  transition: opacity 0.2s ease;
}

.tap-link:hover {
  opacity: 1;
  text-decoration: underline;
}

.more-hint {
  text-align: center;
  padding: 10px;
  color: rgba(255, 255, 255, 0.4);
  font-size: 0.75rem;
  font-style: italic;
}

.search-hint {
  text-align: center;
  padding: 40px 20px;
  color: rgba(255, 255, 255, 0.5);
}

.hint-icon {
  font-size: 2rem;
  margin-bottom: 12px;
}

.search-hint p {
  margin: 0 0 16px 0;
}

.shortcuts {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.shortcut {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.4);
}

kbd {
  padding: 2px 8px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 4px;
  font-family: inherit;
  font-size: 0.75rem;
}

.spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 动画 */
.modal-enter-active,
.modal-leave-active {
  transition: all 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
