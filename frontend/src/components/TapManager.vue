<script setup>
import { ref, computed, reactive } from 'vue'

// Props
const props = defineProps({
  taps: {
    type: Array,
    default: () => []
  },
  isProcessing: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['add', 'remove', 'update', 'update-all'])

// Local state
const newTapName = ref('')
const showAddForm = ref(false)
const processingMap = reactive({})
const selectedTap = ref(null)
const showTapInfo = ref(false)
const packageCounts = ref({})

// Computed properties
const hasTaps = computed(() => props.taps.length > 0)
const officialTaps = computed(() => props.taps.filter(tap => tap.official))
const thirdPartyTaps = computed(() => props.taps.filter(tap => !tap.official))

// Methods
function handleAddTap() {
  if (!newTapName.value.trim()) {
    return
  }
  
  emit('add', newTapName.value.trim())
  newTapName.value = ''
  showAddForm.value = false
}

function handleRemoveTap(tapName) {
  if (window.confirm(`确定要移除 tap "${tapName}" 吗？`)) {
    emit('remove', tapName)
  }
}

function handleUpdateTap(tapName) {
  processingMap[tapName] = true
  emit('update', tapName)
  setTimeout(() => {
    delete processingMap[tapName]
  }, 2000)
}

function handleUpdateAll() {
  emit('update-all')
}

function handleShowInfo(tap) {
  selectedTap.value = tap
  showTapInfo.value = true
  // 模拟获取 package 数量（实际应由后端提供）
  packageCounts.value[tap.name] = {
    formulae: Math.floor(Math.random() * 1000),
    casks: Math.floor(Math.random() * 50)
  }
}

function formatDescription(desc) {
  if (!desc) return '暂无描述'
  return desc.length > 100 ? desc.substring(0, 100) + '...' : desc
}
</script>

<template>
  <div class="tap-manager">
    <!-- 头部操作栏 -->
    <div class="tap-header">
      <div class="header-left">
        <h3 class="tap-title">Taps 管理</h3>
        <span class="tap-count">{{ taps.length }} taps</span>
      </div>
      <div class="header-right">
        <button 
          @click="showAddForm = !showAddForm" 
          class="btn btn-outline"
          :disabled="isProcessing"
        >
          <span v-if="showAddForm">取消</span>
          <span v-else>+ 添加 Tap</span>
        </button>
        <button 
          @click="handleUpdateAll" 
          class="btn btn-primary"
          :disabled="isProcessing"
        >
          <span v-if="isProcessing" class="mini-loader"></span>
          <span v-else><svg class="btn-icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12a9 9 0 11-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/></svg>更新所有</span>
        </button>
      </div>
    </div>

    <!-- 添加 tap 表单 -->
    <Transition name="slide-down">
      <div v-if="showAddForm" class="add-tap-form">
        <div class="form-group">
          <label for="tap-name">Tap 名称</label>
          <input
            id="tap-name"
            v-model="newTapName"
            type="text"
            placeholder="例如: homebrew/cask-fonts 或 user/repo"
            class="form-input"
            @keyup.enter="handleAddTap"
          />
          <div class="form-hint">
            支持 GitHub 仓库格式：username/repository 或 homebrew/cask-fonts
          </div>
        </div>
        <div class="form-actions">
          <button @click="handleAddTap" class="btn btn-primary" :disabled="!newTapName.trim()">
            添加
          </button>
        </div>
      </div>
    </Transition>

    <!-- Tap 列表 -->
    <div class="tap-list">
      <!-- 官方 Taps -->
      <div v-if="officialTaps.length > 0" class="tap-section">
        <h4 class="section-title">
          <span class="badge badge-official">官方</span>
          官方 Taps
        </h4>
        <div class="tap-grid">
          <div 
            v-for="tap in officialTaps" 
            :key="tap.name"
            class="tap-card official"
            @click="handleShowInfo(tap)"
          >
            <div class="tap-card-header">
              <div class="tap-name-group">
                <span class="tap-name">{{ tap.name }}</span>
                <span class="tap-badge official">官方</span>
              </div>
              <div class="tap-actions">
                <button 
                  @click.stop="handleUpdateTap(tap.name)"
                  class="btn-icon"
                  :title="`更新 ${tap.name}`"
                  :disabled="processingMap[tap.name]"
                >
                  <span v-if="processingMap[tap.name]" class="mini-loader"></span>
                  <span v-else><svg class="icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12a9 9 0 11-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/></svg></span>
                </button>
              </div>
            </div>
            
            <div class="tap-card-body">
              <div class="tap-description">{{ formatDescription(tap.description) }}</div>
              <div class="tap-meta">
                <span class="tap-url">{{ tap.url || 'N/A' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 第三方 Taps -->
      <div v-if="thirdPartyTaps.length > 0" class="tap-section">
        <h4 class="section-title">
          <span class="badge badge-third">第三方</span>
          第三方 Taps
        </h4>
        <div class="tap-grid">
          <div 
            v-for="tap in thirdPartyTaps" 
            :key="tap.name"
            class="tap-card third-party"
            @click="handleShowInfo(tap)"
          >
            <div class="tap-card-header">
              <div class="tap-name-group">
                <span class="tap-name">{{ tap.name }}</span>
                <span class="tap-badge third">第三方</span>
              </div>
              <div class="tap-actions">
                <button 
                  @click.stop="handleUpdateTap(tap.name)"
                  class="btn-icon"
                  :title="`更新 ${tap.name}`"
                  :disabled="processingMap[tap.name]"
                >
                  <span v-if="processingMap[tap.name]" class="mini-loader"></span>
                  <span v-else><svg class="icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12a9 9 0 11-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/></svg></span>
                </button>
                <button 
                  @click.stop="handleRemoveTap(tap.name)"
                  class="btn-icon btn-danger"
                  title="移除 tap"
                >
                  <svg class="icon-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
                </button>
              </div>
            </div>
            
            <div class="tap-card-body">
              <div class="tap-description">{{ formatDescription(tap.description) }}</div>
              <div class="tap-meta">
                <span class="tap-url">{{ tap.url || 'Git 仓库' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="!hasTaps" class="empty-state">
        <div class="empty-icon">📦</div>
        <h4>暂无 Tap</h4>
        <p>Tap 是 Homebrew 的软件源仓库，点击上方"添加 Tap"按钮开始添加</p>
      </div>
    </div>

    <!-- Tap 详情模态框 -->
    <Transition name="modal">
      <div v-if="showTapInfo && selectedTap" class="modal-overlay" @click="showTapInfo = false">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h3>{{ selectedTap.name }}</h3>
            <button @click="showTapInfo = false" class="modal-close">×</button>
          </div>
          <div class="modal-body">
            <div class="tap-info">
              <div class="info-row">
                <span class="info-label">类型：</span>
                <span class="info-value">
                  <span :class="selectedTap.official ? 'badge official' : 'badge third'">
                    {{ selectedTap.official ? '官方' : '第三方' }}
                  </span>
                </span>
              </div>
              
              <div class="info-row">
                <span class="info-label">Git URL：</span>
                <span class="info-value">{{ selectedTap.url || '未指定' }}</span>
              </div>
              
              <div class="info-row" v-if="selectedTap.description">
                <span class="info-label">描述：</span>
                <span class="info-value">{{ selectedTap.description }}</span>
              </div>
              
              <div class="info-row" v-if="packageCounts[selectedTap.name]">
                <span class="info-label">包数量：</span>
                <span class="info-value">
                  {{ packageCounts[selectedTap.name].formulae }} formulae, 
                  {{ packageCounts[selectedTap.name].casks }} casks
                </span>
              </div>
              
              <div class="info-row">
                <span class="info-label">添加命令：</span>
                <code class="info-value">brew tap {{ selectedTap.name }}</code>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button @click="showTapInfo = false" class="btn">关闭</button>
            <button 
              v-if="!selectedTap.official"
              @click="handleRemoveTap(selectedTap.name); showTapInfo = false"
              class="btn btn-danger"
            >
              移除
            </button>
            <button 
              @click="handleUpdateTap(selectedTap.name)"
              class="btn btn-primary"
              :disabled="processingMap[selectedTap.name]"
            >
              <span v-if="processingMap[selectedTap.name]" class="mini-loader"></span>
              <span v-else>更新</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.tap-manager {
  background: rgba(255, 255, 255, 0.7);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.tap-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.tap-title {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: #333;
}

.tap-count {
  background: rgba(52, 152, 219, 0.1);
  color: #3498db;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 500;
}

.header-right {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 8px 16px;
  border-radius: 8px;
  border: none;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.btn:hover {
  transform: translateY(-1px);
}

.btn:active {
  transform: translateY(0);
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.btn-outline {
  background: transparent;
  border: 1px solid #3498db;
  color: #3498db;
}

.btn-outline:hover {
  background: rgba(52, 152, 219, 0.1);
}

.btn-primary {
  background: linear-gradient(135deg, #3498db, #2980b9);
  color: white;
  box-shadow: 0 2px 6px rgba(52, 152, 219, 0.3);
}

.btn-primary:hover {
  box-shadow: 0 4px 12px rgba(52, 152, 219, 0.4);
}

.btn-danger {
  background: linear-gradient(135deg, #e74c3c, #c0392b);
  color: white;
  box-shadow: 0 2px 6px rgba(231, 76, 60, 0.3);
}

.btn-danger:hover {
  box-shadow: 0 4px 12px rgba(231, 76, 60, 0.4);
}

.btn-icon {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 1.125rem;
  padding: 4px 8px;
  border-radius: 6px;
  opacity: 0.7;
  transition: all 0.2s ease;
}

.btn-icon:hover {
  opacity: 1;
  background: rgba(0, 0, 0, 0.05);
}

.btn-icon.btn-danger:hover {
  background: rgba(231, 76, 60, 0.1);
}

/* 添加 tap 表单 */
.add-tap-form {
  background: rgba(249, 250, 251, 0.8);
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  border: 1px solid rgba(226, 232, 240, 0.8);
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #4a5568;
}

.form-input {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 0.875rem;
  transition: border-color 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.1);
}

.form-hint {
  margin-top: 6px;
  font-size: 0.75rem;
  color: #718096;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* Tap 列表 */
.tap-list {
  margin-top: 20px;
}

.tap-section {
  margin-bottom: 30px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0 0 16px 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: #333;
}

.badge {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.badge-official {
  background: rgba(52, 152, 219, 0.1);
  color: #3498db;
}

.badge-third {
  background: rgba(46, 204, 113, 0.1);
  color: #2ecc71;
}

.tap-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.tap-card {
  background: white;
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.tap-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: #3498db;
}

.tap-card.official:hover {
  border-color: #3498db;
}

.tap-card.third-party:hover {
  border-color: #2ecc71;
}

.tap-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.tap-name-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tap-name {
  font-weight: 600;
  color: #333;
  font-size: 1rem;
}

.tap-badge {
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.625rem;
  font-weight: 600;
  text-transform: uppercase;
}

.tap-badge.official {
  background: rgba(52, 152, 219, 0.1);
  color: #3498db;
}

.tap-badge.third {
  background: rgba(46, 204, 113, 0.1);
  color: #2ecc71;
}

.tap-actions {
  display: flex;
  gap: 4px;
}

.tap-card-body {
  color: #666;
}

.tap-description {
  font-size: 0.875rem;
  line-height: 1.4;
  margin-bottom: 8px;
  color: #4a5568;
}

.tap-meta {
  font-size: 0.75rem;
  color: #718096;
}

.tap-url {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  display: block;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: #718096;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state h4 {
  margin: 0 0 8px 0;
  color: #4a5568;
}

.empty-state p {
  margin: 0;
  font-size: 0.875rem;
  max-width: 400px;
  margin: 0 auto;
}

/* 模态框 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e2e8f0;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #333;
}

.modal-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #718096;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-close:hover {
  background: rgba(0, 0, 0, 0.05);
  color: #333;
}

.modal-body {
  padding: 20px;
}

.tap-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-row {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.info-label {
  font-weight: 600;
  color: #4a5568;
  min-width: 80px;
  font-size: 0.875rem;
}

.info-value {
  flex: 1;
  color: #333;
  font-size: 0.875rem;
  word-break: break-all;
}

.info-value code {
  background: #f7fafc;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.75rem;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 20px;
  border-top: 1px solid #e2e8f0;
}

/* 动画 */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
  max-height: 200px;
  overflow: hidden;
}

.slide-down-enter-from,
.slide-down-leave-to {
  max-height: 0;
  opacity: 0;
  transform: translateY(-10px);
}

.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

/* 加载动画 */
.mini-loader {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 0.8s linear infinite;
}

/* SVG 图标样式 */
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

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>