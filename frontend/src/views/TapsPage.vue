<script setup>
import { ref, computed } from 'vue'
import { BrowserOpenURL } from '../../wailsjs/runtime'
import { SearchPackages } from '../../wailsjs/go/main/App'

const props = defineProps({
  taps: {
    type: Array,
    default: () => []
  },
  isProcessing: {
    type: Boolean,
    default: false
  },
  processingMap: {
    type: Map,
    default: () => new Map()
  }
})

const emit = defineEmits(['add', 'remove', 'update', 'update-all'])

const localSearch = ref('')
const newTapName = ref('')
const showAddForm = ref(false)
const searchResults = ref([])
const searching = ref(false)

const officialTaps = computed(() => props.taps.filter(tap => tap.official))
const thirdPartyTaps = computed(() => props.taps.filter(tap => !tap.official))

const filteredOfficialTaps = computed(() => {
  if (!localSearch.value) return officialTaps.value
  const query = localSearch.value.toLowerCase()
  return officialTaps.value.filter(tap => 
    tap.name.toLowerCase().includes(query)
  )
})

const filteredThirdPartyTaps = computed(() => {
  if (!localSearch.value) return thirdPartyTaps.value
  const query = localSearch.value.toLowerCase()
  return thirdPartyTaps.value.filter(tap => 
    tap.name.toLowerCase().includes(query)
  )
})

// 检查 tap 是否已安装
function isTapInstalled(tapName) {
  if (!tapName || !props.taps || !props.taps.length) return false
  
  return props.taps.some(tap => {
    // 精确匹配
    if (tap.name === tapName) return true
    
    // 处理可能的格式差异
    const installedName = tap.name.toLowerCase().trim()
    const searchName = tapName.toLowerCase().trim()
    
    if (installedName === searchName) return true
    
    // 处理不同格式：
    // 已安装: "kong/deck", 搜索结果: "deck"
    // 已安装: "homebrew/core", 搜索结果: "core"
    const installedParts = installedName.split('/')
    const searchParts = searchName.split('/')
    
    // 如果搜索结果只有一部分（没有 /），比较最后一部分
    if (searchParts.length === 1 && installedParts.length === 2) {
      return installedParts[1] === searchParts[0]
    }
    
    return false
  })
}

async function searchTaps() {
  if (!newTapName.value.trim()) {
    searchResults.value = []
    return
  }
  
  searching.value = true
  try {
    const results = await SearchPackages(newTapName.value.trim())
    searchResults.value = results.taps || []
    
  } catch (err) {
    console.error('搜索失败:', err)
    searchResults.value = []
  } finally {
    searching.value = false
  }
}

function handleAddTap() {
  if (!newTapName.value.trim()) return
  emit('add', newTapName.value.trim())
  newTapName.value = ''
  searchResults.value = []
  showAddForm.value = false
}

function handleAddTapFromSearch(tapName) {
  emit('add', tapName)
  newTapName.value = ''
  searchResults.value = []
  showAddForm.value = false
}

const showConfirmDialog = ref(false)
const confirmTapName = ref('')

function handleRemoveTap(tapName) {
  confirmTapName.value = tapName
  showConfirmDialog.value = true
}

function confirmRemove() {
  emit('remove', confirmTapName.value)
  showConfirmDialog.value = false
  confirmTapName.value = ''
}

function cancelRemove() {
  showConfirmDialog.value = false
  confirmTapName.value = ''
}

function handleUpdateTap(tapName) {
  emit('update', tapName)
}

function handleUpdateAll() {
  emit('update-all')
}

function formatDescription(desc) {
  if (!desc) return '暂无描述'
  return desc.length > 80 ? desc.substring(0, 80) + '...' : desc
}

// 生成 GitHub 地址
function getGitHubUrl(tap) {
  if (tap.url && tap.url.includes('github.com')) {
    // 如果后端提供了 url，直接使用
    return tap.url
  }
  // 否则根据 tap 名称生成
  // homebrew/core -> https://github.com/homebrew/homebrew-core
  // user/repo -> https://github.com/user/homebrew-repo 或 https://github.com/user/repo
  const parts = tap.name.split('/')
  if (parts.length === 2) {
    return `https://github.com/${parts[0]}/homebrew-${parts[1]}`
  }
  return `https://github.com/${tap.name}`
}

// 提取简洁的 GitHub 路径显示
function getGitHubDisplay(tap) {
  const url = getGitHubUrl(tap)
  return url.replace('https://github.com/', '')
}

// 打开 GitHub 链接
function openGitHubUrl(tap) {
  const url = getGitHubUrl(tap)
  BrowserOpenURL(url)
}
</script>

<template>
  <div class="taps-page">
    <div class="page-header">
      <h2 class="page-title">Taps 管理</h2>
      <div class="header-actions">
        <button type="button" @click="showAddForm = !showAddForm" class="btn btn-outline" :disabled="isProcessing">
          {{ showAddForm ? '取消' : '+ 添加 Tap' }}
        </button>
        <button type="button" @click="handleUpdateAll" class="btn btn-primary" :disabled="isProcessing">
          <svg class="btn-icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12a9 9 0 11-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/>
            <path d="M21 3v5h-5"/>
          </svg>
          更新所有
        </button>
      </div>
    </div>
    
    <!-- 本地搜索 -->
    <div class="local-search">
      <span class="search-icon">🔍</span>
      <input 
        v-model="localSearch" 
        type="text" 
        placeholder="搜索已添加的 taps..." 
        class="search-input"
      />
    </div>

    <!-- 添加 Tap 表单 -->
    <div v-if="showAddForm" class="add-tap-form">
      <input
        v-model="newTapName"
        type="text"
        placeholder="搜索或输入 tap 名称..."
        class="form-input"
        @keyup.enter="searchTaps"
      />
      <button type="button" @click="searchTaps" class="btn btn-primary" :disabled="!newTapName.trim() || searching">
        <span v-if="searching" class="spinner"></span>
        <span v-else>搜索</span>
      </button>
    </div>

    <!-- 搜索结果 -->
    <div v-if="showAddForm && searchResults.length > 0" class="search-results">
      <h4 class="results-title">搜索结果 ({{ searchResults.length }})</h4>
      <div class="results-list">
        <div v-for="tap in searchResults" :key="tap.name" class="result-item">
          <div class="result-info">
            <div class="result-name">
              {{ tap.name }}
              <span v-if="isTapInstalled(tap.name)" class="installed-badge">已安装</span>
            </div>
            <div class="result-desc">{{ tap.description || '暂无描述' }}</div>
            <div class="result-stars">⭐ {{ tap.stars || 0 }}</div>
            <button 
              class="result-url"
              @click="BrowserOpenURL(tap.url)"
            >
              <span class="url-icon">🔗</span>
              {{ tap.url.replace('https://github.com/', '') }}
            </button>
          </div>
          <button 
            v-if="!isTapInstalled(tap.name)"
            type="button"
            @click="handleAddTapFromSearch(tap.name)" 
            class="btn btn-primary btn-small"
            :disabled="processingMap.get(`tap-add-${tap.name}`)"
          >
            <span v-if="processingMap.get(`tap-add-${tap.name}`)" class="spinner"></span>
            <span v-else>添加</span>
          </button>
          <span v-else class="installed-text">✓ 已添加</span>
        </div>
      </div>
    </div>

    <!-- 官方 Taps -->
    <div v-if="filteredOfficialTaps.length > 0" class="tap-section">
      <h3 class="section-title">
        <span class="badge badge-official">官方</span>
        官方 Taps ({{ filteredOfficialTaps.length }})
      </h3>
      <div class="tap-list">
        <div v-for="tap in filteredOfficialTaps" :key="tap.name" class="tap-card official">
          <div class="tap-info">
            <div class="tap-name">{{ tap.name }}</div>
            <div class="tap-desc">{{ formatDescription(tap.description) }}</div>
            <button 
              class="tap-url"
              @click.stop="openGitHubUrl(tap)"
            >
              <span class="url-icon">🔗</span>
              {{ getGitHubDisplay(tap) }}
            </button>
          </div>
          <button 
            type="button"
            @click="handleUpdateTap(tap.name)" 
            class="btn-icon"
            :disabled="processingMap.get(`tap-update-${tap.name}`)"
            title="更新"
          >
            <span v-if="processingMap.get(`tap-update-${tap.name}`)" class="spinner"></span>
            <span v-else><svg class="icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12a9 9 0 11-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/></svg></span>
          </button>
        </div>
      </div>
    </div>

    <!-- 第三方 Taps -->
    <div v-if="filteredThirdPartyTaps.length > 0" class="tap-section">
      <h3 class="section-title">
        <span class="badge badge-third">第三方</span>
        第三方 Taps ({{ filteredThirdPartyTaps.length }})
      </h3>
      <div class="tap-list">
        <div v-for="tap in filteredThirdPartyTaps" :key="tap.name" class="tap-card third-party">
          <div class="tap-info">
            <div class="tap-name">{{ tap.name }}</div>
            <div class="tap-desc">{{ formatDescription(tap.description) }}</div>
            <button 
              class="tap-url"
              @click.stop="openGitHubUrl(tap)"
            >
              <span class="url-icon">🔗</span>
              {{ getGitHubDisplay(tap) }}
            </button>
          </div>
          <div class="tap-actions">
            <button 
              type="button"
              @click="handleUpdateTap(tap.name)" 
              class="btn-icon"
              :disabled="processingMap.get(`tap-update-${tap.name}`)"
              title="更新"
            >
              <span v-if="processingMap.get(`tap-update-${tap.name}`)" class="spinner"></span>
              <span v-else><svg class="icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12a9 9 0 11-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/></svg></span>
            </button>
            <button 
              type="button"
              @click="handleRemoveTap(tap.name)" 
              class="btn-icon btn-danger"
              :disabled="processingMap.get(`tap-remove-${tap.name}`)"
              title="移除"
            >
              <span v-if="processingMap.get(`tap-remove-${tap.name}`)" class="spinner"></span>
              <span v-else><svg class="icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg></span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 确认对话框 -->
    <div v-if="showConfirmDialog" class="confirm-overlay" @click="cancelRemove">
      <div class="confirm-dialog" @click.stop>
        <h3>确认移除</h3>
        <p>确定要移除 tap "{{ confirmTapName }}" 吗？</p>
        <div class="confirm-actions">
          <button type="button" @click="cancelRemove" class="btn btn-outline">取消</button>
          <button type="button" @click="confirmRemove" class="btn btn-danger">确认移除</button>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="props.taps.length === 0" class="empty-state">
      <div class="empty-icon">🔗</div>
      <p>暂无添加的 Taps</p>
      <p class="hint">点击"添加 Tap"按钮开始添加第三方软件源</p>
    </div>
    <div v-else-if="filteredOfficialTaps.length === 0 && filteredThirdPartyTaps.length === 0" class="empty-state">
      <div class="empty-icon">🔍</div>
      <p>未找到匹配的 Taps</p>
    </div>
  </div>
</template>

<style scoped>
.taps-page {
  padding: 20px;
  color: #fff;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.taps-page::-webkit-scrollbar {
  width: 5px;
}

.taps-page::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
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

.header-actions {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 8px 16px;
  border-radius: 6px;
  border: none;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-outline {
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: #fff;
}

.btn-outline:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.1);
}

.btn-primary {
  background: #007AFF;
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  background: #0056CC;
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

.add-tap-form {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
}

.form-input {
  flex: 1;
  padding: 10px 14px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #fff;
  font-size: 0.875rem;
  outline: none;
}

.form-input::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.form-input:focus {
  border-color: #007AFF;
}

.search-results {
  margin-bottom: 20px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
}

.results-title {
  margin: 0 0 12px 0;
  font-size: 0.875rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.7);
}

.results-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 300px;
  overflow-y: auto;
}

.results-list::-webkit-scrollbar {
  width: 5px;
}

.results-list::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
}

.result-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 6px;
}

.result-info {
  flex: 1;
  min-width: 0;
  margin-right: 10px;
}

.result-name {
  font-weight: 600;
  font-size: 0.875rem;
  margin-bottom: 4px;
  font-family: 'Monaco', 'Menlo', monospace;
}

.result-desc {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.6);
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.result-stars {
  font-size: 0.7rem;
  color: #FFD700;
}

.result-url {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.7rem;
  color: #007AFF;
  text-decoration: none;
  opacity: 0.8;
  transition: opacity 0.2s ease;
  background: transparent;
  border: none;
  padding: 0;
  cursor: pointer;
  font-family: inherit;
  margin-top: 4px;
}

.result-url:hover {
  opacity: 1;
  text-decoration: underline;
}

.installed-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  background: rgba(46, 204, 113, 0.2);
  color: #2ECC71;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  margin-left: 8px;
}

.installed-text {
  display: inline-flex;
  align-items: center;
  padding: 6px 12px;
  color: #2ECC71;
  font-size: 0.75rem;
  font-weight: 500;
}

.btn-small {
  padding: 6px 12px;
  font-size: 0.75rem;
}

.tap-section {
  margin-bottom: 24px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0 0 12px 0;
  font-size: 1rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
}

.badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
}

.badge-official {
  background: rgba(0, 122, 255, 0.2);
  color: #007AFF;
}

.badge-third {
  background: rgba(46, 204, 113, 0.2);
  color: #2ECC71;
}

.tap-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tap-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  transition: all 0.2s ease;
}

.tap-card:hover {
  background: rgba(255, 255, 255, 0.06);
}

.tap-card.official {
  border-left: 3px solid #007AFF;
}

.tap-card.third-party {
  border-left: 3px solid #2ECC71;
}

.tap-info {
  flex: 1;
  min-width: 0;
}

.tap-name {
  font-weight: 600;
  font-family: 'Monaco', 'Menlo', monospace;
  margin-bottom: 4px;
}

.tap-desc {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.5);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 6px;
}

.tap-url {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.7rem;
  color: #007AFF;
  text-decoration: none;
  opacity: 0.8;
  transition: opacity 0.2s ease;
  background: transparent;
  border: none;
  padding: 0;
  cursor: pointer;
  font-family: inherit;
}

.tap-url:hover {
  opacity: 1;
  text-decoration: underline;
}

.url-icon {
  font-size: 0.7rem;
}

.tap-actions {
  display: flex;
  gap: 6px;
}

.btn-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-icon:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.btn-icon.btn-danger:hover {
  background: rgba(231, 76, 60, 0.2);
  color: #E74C3C;
}

.btn-icon:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.icon-svg, .btn-icon-svg {
  width: 14px;
  height: 14px;
  vertical-align: middle;
}

.btn-icon-svg {
  width: 14px;
  height: 14px;
  margin-right: 6px;
}

.spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.2);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
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
  margin: 0 0 8px 0;
}

.confirm-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.confirm-dialog {
  background: #1c1c1e;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 24px;
  width: 90%;
  max-width: 400px;
  text-align: center;
}

.confirm-dialog h3 {
  margin: 0 0 16px 0;
  font-size: 1.25rem;
  color: #fff;
}

.confirm-dialog p {
  margin: 0 0 24px 0;
  color: rgba(255, 255, 255, 0.7);
}

.confirm-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.btn-danger {
  background: #E74C3C;
  color: #fff;
}

.btn-danger:hover:not(:disabled) {
  background: #C0392B;
}

.hint {
  font-size: 0.875rem;
  opacity: 0.7;
}
</style>
