<script setup>
/**
 * LaunchPad - 服务组面板
 * 支持配置服务组（Brew 服务 + Docker 容器），批量启动/停止/重启
 */
import { ref, computed, onMounted } from 'vue'
import { 
  GetLaunchGroups, SaveLaunchGroup, DeleteLaunchGroup, 
  LaunchGroup, StopGroup, RestartGroup, GetAvailableServices 
} from '../../wailsjs/go/main/App'

const emit = defineEmits(['show-toast'])

// 状态
const groups = ref([])
const availableServices = ref({ brew_services: [], docker_containers: [] })
const loading = ref(true)
const launching = ref('') // 正在启动的组名
const stopping = ref('')  // 正在停止的组名
const restarting = ref('') // 正在重启的组名
const launchResults = ref(null) // 操作结果列表

// Modal 状态
const showModal = ref(false)
const editingGroup = ref(null) // null = 新建，有值 = 编辑
const formName = ref('')
const selectedBrew = ref([])
const selectedDocker = ref([])
const saving = ref(false)
const confirmingDelete = ref('') // 正在确认删除的组名

// 筛选
const modalSearch = ref('')
const modalFilter = ref('all') // 'all' | 'brew' | 'docker'

// 过滤后的 Brew 服务
function filteredBrewServices() {
  if (modalFilter.value === 'docker') return []
  if (!modalSearch.value) return availableServices.value.brew_services || []
  const q = modalSearch.value.toLowerCase()
  return (availableServices.value.brew_services || []).filter(s =>
    s.name.toLowerCase().includes(q)
  )
}

// 过滤后的 Docker 容器
function filteredDockerContainers() {
  if (modalFilter.value === 'brew') return []
  if (!modalSearch.value) return availableServices.value.docker_containers || []
  const q = modalSearch.value.toLowerCase()
  return (availableServices.value.docker_containers || []).filter(c =>
    c.name.toLowerCase().includes(q) || c.image.toLowerCase().includes(q)
  )
}

// 服务状态标签
function brewStatusLabel(status) {
  const map = { started: '运行中', stopped: '已停止', none: '已停止' }
  return map[status] || status
}

function containerStateLabel(state) {
  const map = { running: '运行中', exited: '已停止', paused: '已暂停' }
  return map[state] || state
}

// 加载数据
async function loadData() {
  loading.value = true
  try {
    const [g, s] = await Promise.all([
      GetLaunchGroups(),
      GetAvailableServices()
    ])
    groups.value = g || []
    availableServices.value = s || { brew_services: [], docker_containers: [] }
  } catch (err) {
    console.error('加载服务组失败:', err)
  } finally {
    loading.value = false
  }
}

// 打开新建 Modal
function openCreateModal() {
  editingGroup.value = null
  formName.value = ''
  selectedBrew.value = []
  selectedDocker.value = []
  modalSearch.value = ''
  modalFilter.value = 'all'
  showModal.value = true
}

// 打开编辑 Modal
function openEditModal(group) {
  editingGroup.value = group
  formName.value = group.name
  selectedBrew.value = [...(group.brew_services || [])]
  selectedDocker.value = [...(group.docker_containers || [])]
  modalSearch.value = ''
  modalFilter.value = 'all'
  showModal.value = true
}

// 关闭 Modal
function closeModal() {
  showModal.value = false
  editingGroup.value = null
}

// 切换选择
function toggleBrew(name) {
  const idx = selectedBrew.value.indexOf(name)
  if (idx > -1) {
    selectedBrew.value.splice(idx, 1)
  } else {
    selectedBrew.value.push(name)
  }
}

function toggleDocker(name) {
  const idx = selectedDocker.value.indexOf(name)
  if (idx > -1) {
    selectedDocker.value.splice(idx, 1)
  } else {
    selectedDocker.value.push(name)
  }
}

function isBrewSelected(name) {
  return selectedBrew.value.includes(name)
}

function isDockerSelected(name) {
  return selectedDocker.value.includes(name)
}

// 保存服务组
async function handleSave() {
  if (!formName.value.trim()) {
    emit('show-toast', { msg: '请输入服务组名称', type: 'error' })
    return
  }
  if (selectedBrew.value.length === 0 && selectedDocker.value.length === 0) {
    emit('show-toast', { msg: '请至少选择一个服务或容器', type: 'error' })
    return
  }

  saving.value = true
  try {
    const result = await SaveLaunchGroup(formName.value.trim(), selectedBrew.value, selectedDocker.value)
    emit('show-toast', { msg: result.message, type: result.success ? 'success' : 'error' })
    if (result.success) {
      closeModal()
      await loadData()
    }
  } catch (err) {
    emit('show-toast', { msg: '保存失败: ' + err, type: 'error' })
  } finally {
    saving.value = false
  }
}

// 删除服务组 - 点击后显示确认，再次点击确认删除
async function handleDelete(group) {
  if (confirmingDelete.value !== group.name) {
    confirmingDelete.value = group.name
    return
  }
  try {
    const result = await DeleteLaunchGroup(group.name)
    emit('show-toast', { msg: result.message, type: result.success ? 'success' : 'error' })
    if (result.success) {
      confirmingDelete.value = ''
      await loadData()
    }
  } catch (err) {
    emit('show-toast', { msg: '删除失败: ' + err, type: 'error' })
  } finally {
    confirmingDelete.value = ''
  }
}

// 点击外部取消确认
function handleCardClick(group) {
  if (confirmingDelete.value && confirmingDelete.value !== group.name) {
    confirmingDelete.value = ''
  }
}

// 全部启动
async function handleLaunch(group) {
  launching.value = group.name
  launchResults.value = null
  try {
    const result = await LaunchGroup(group.name)
    if (result.success) {
      emit('show-toast', { msg: result.message, type: 'success' })
      launchResults.value = result.data || []
      // 刷新服务状态
      await loadData()
    } else {
      emit('show-toast', { msg: result.message, type: 'error' })
    }
  } catch (err) {
    emit('show-toast', { msg: '启动失败: ' + err, type: 'error' })
  } finally {
    launching.value = ''
  }
}

// 统计项数
function itemCount(group) {
  return (group.brew_services?.length || 0) + (group.docker_containers?.length || 0)
}

// 计算服务组的运行状态
// 'running' = 全部运行, 'partial' = 部分运行, 'stopped' = 全部停止, 'unknown' = 无法判断
function groupStatus(group) {
  const brewSvcs = availableServices.value.brew_services || []
  const dockerCntrs = availableServices.value.docker_containers || []

  // 构建 brew 服务名->状态 映射
  const brewMap = {}
  for (const s of brewSvcs) {
    brewMap[s.name] = s.status
  }
  // 构建 docker 容器名->状态 映射
  const dockerMap = {}
  for (const c of dockerCntrs) {
    dockerMap[c.name] = c.state
  }

  let runningCount = 0
  let totalCount = 0

  // 检查 brew 服务
  for (const name of (group.brew_services || [])) {
    totalCount++
    if (brewMap[name] === 'started') runningCount++
  }
  // 检查 docker 容器
  for (const name of (group.docker_containers || [])) {
    totalCount++
    if (dockerMap[name] === 'running') runningCount++
  }

  if (totalCount === 0) return 'unknown'
  if (runningCount === totalCount) return 'running'
  if (runningCount === 0) return 'stopped'
  return 'partial'
}

// 状态标签文案和样式
function statusLabel(status) {
  const map = { running: '运行中', partial: '部分运行', stopped: '已停止', unknown: '未知' }
  return map[status] || status
}

// 全部停止
async function handleStop(group) {
  stopping.value = group.name
  launchResults.value = null
  try {
    const result = await StopGroup(group.name)
    if (result.success) {
      emit('show-toast', { msg: result.message, type: 'success' })
      launchResults.value = result.data || []
      // 刷新服务状态
      await loadData()
    } else {
      emit('show-toast', { msg: result.message, type: 'error' })
    }
  } catch (err) {
    emit('show-toast', { msg: '停止失败: ' + err, type: 'error' })
  } finally {
    stopping.value = ''
  }
}

// 重启
async function handleRestart(group) {
  restarting.value = group.name
  launchResults.value = null
  try {
    const result = await RestartGroup(group.name)
    if (result.success) {
      emit('show-toast', { msg: result.message, type: 'success' })
      launchResults.value = result.data || []
      // 刷新服务状态
      await loadData()
    } else {
      emit('show-toast', { msg: result.message, type: 'error' })
    }
  } catch (err) {
    emit('show-toast', { msg: '重启失败: ' + err, type: 'error' })
  } finally {
    restarting.value = ''
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="launchpad">
    <div class="launchpad-header">
      <h3 class="launchpad-title">🚀 服务组</h3>
      <button class="add-group-btn" @click="openCreateModal">+ 新建服务组</button>
    </div>

    <!-- 加载中 -->
    <div v-if="loading" class="lp-loading">加载中...</div>

    <!-- 空状态 -->
    <div v-else-if="groups.length === 0" class="lp-empty">
      <p>还没有服务组，点击"新建服务组"来配置你的开发环境</p>
    </div>

    <!-- 启动组列表 -->
    <div v-else class="lp-groups">
      <div v-for="group in groups" :key="group.name" class="lp-group-card" :class="groupStatus(group)" @click="handleCardClick(group)">
        <div class="lp-group-info">
          <div class="lp-group-name-row">
            <span class="lp-group-name">{{ group.name }}</span>
            <span class="lp-status-badge" :class="groupStatus(group)">{{ statusLabel(groupStatus(group)) }}</span>
          </div>
          <div class="lp-group-meta">
            <span v-if="group.brew_services?.length" class="lp-meta-tag brew">
              🍺 {{ group.brew_services.length }} 个 Brew 服务
            </span>
            <span v-if="group.docker_containers?.length" class="lp-meta-tag docker">
              🐳 {{ group.docker_containers.length }} 个容器
            </span>
          </div>
        </div>

        <div class="lp-group-actions">
          <!-- 运行中：显示重启 + 全部停止 -->
          <template v-if="groupStatus(group) === 'running'">
            <button
              class="lp-restart-btn"
              :disabled="!!restarting || !!stopping || !!launching"
              @click="handleRestart(group)"
            >
              <span v-if="restarting === group.name" class="btn-spinner"></span>
              <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="23 4 23 10 17 10"/>
                <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
              </svg>
              {{ restarting === group.name ? '重启中...' : '重启' }}
            </button>
            <button
              class="lp-stop-btn"
              :disabled="!!stopping || !!launching || !!restarting"
              @click="handleStop(group)"
            >
              <span v-if="stopping === group.name" class="btn-spinner"></span>
              <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/>
              </svg>
              {{ stopping === group.name ? '停止中...' : '全部停止' }}
            </button>
          </template>
          <!-- 部分运行：显示全部启动 + 停止 -->
          <template v-if="groupStatus(group) === 'partial'">
            <button
              class="lp-launch-btn"
              :disabled="!!launching || !!stopping || !!restarting"
              @click="handleLaunch(group)"
            >
              <span v-if="launching === group.name" class="btn-spinner"></span>
              <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <polygon points="5 3 19 12 5 21 5 3"/>
              </svg>
              {{ launching === group.name ? '启动中...' : '全部启动' }}
            </button>
            <button
              class="lp-stop-btn"
              :disabled="!!stopping || !!launching || !!restarting"
              @click="handleStop(group)"
            >
              <span v-if="stopping === group.name" class="btn-spinner"></span>
              <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/>
              </svg>
              {{ stopping === group.name ? '停止中...' : '停止' }}
            </button>
          </template>
          <!-- 全部停止：显示全部启动 -->
          <template v-if="groupStatus(group) === 'stopped' || groupStatus(group) === 'unknown'">
            <button
              class="lp-launch-btn"
              :disabled="!!launching || !!stopping || !!restarting"
              @click="handleLaunch(group)"
            >
              <span v-if="launching === group.name" class="btn-spinner"></span>
              <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <polygon points="5 3 19 12 5 21 5 3"/>
              </svg>
              {{ launching === group.name ? '启动中...' : '全部启动' }}
            </button>
          </template>
          <button class="lp-edit-btn" @click="openEditModal(group)" title="编辑">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
          </button>
          <button
            class="lp-delete-btn"
            :class="{ confirming: confirmingDelete === group.name }"
            @click.stop="handleDelete(group)"
            :title="confirmingDelete === group.name ? '再次点击确认删除' : '删除'"
          >
            <svg v-if="confirmingDelete !== group.name" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="3 6 5 6 21 6"/>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
            </svg>
            <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
            <span v-if="confirmingDelete === group.name" class="lp-delete-confirm-text">确认</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 操作结果 -->
    <div v-if="launchResults" class="lp-results">
      <div class="lp-results-header">
        <span>操作结果</span>
        <button class="lp-results-close" @click="launchResults = null">✕</button>
      </div>
      <div v-for="r in launchResults" :key="r.name" class="lp-result-item" :class="{ success: r.success, fail: !r.success }">
        <span class="lp-result-icon">{{ r.success ? '✓' : '✗' }}</span>
        <span class="lp-result-type">{{ r.type === 'brew' ? '🍺' : '🐳' }}</span>
        <span class="lp-result-name">{{ r.name }}</span>
        <span class="lp-result-msg">{{ r.message }}</span>
      </div>
    </div>

    <!-- 配置 Modal -->
    <div v-if="showModal" class="lp-modal-overlay" @click.self="closeModal">
      <div class="lp-modal">
        <div class="lp-modal-header">
          <h4>{{ editingGroup ? '编辑服务组' : '新建服务组' }}</h4>
          <button class="lp-modal-close" @click="closeModal">✕</button>
        </div>

        <div class="lp-modal-body">
          <!-- 组名称 -->
          <div class="lp-form-group">
            <label class="lp-label">服务组名称</label>
            <input
              v-model="formName"
              type="text"
              class="lp-input"
              placeholder="例如：PHP 开发环境"
              :disabled="!!editingGroup"
            />
          </div>

          <!-- 搜索和筛选 -->
          <div class="lp-modal-search">
            <input
              v-model="modalSearch"
              type="text"
              class="lp-input"
              placeholder="搜索服务或容器..."
            />
            <div class="lp-filter-tabs">
              <button 
                :class="{ active: modalFilter === 'all' }" 
                @click="modalFilter = 'all'"
              >全部</button>
              <button 
                :class="{ active: modalFilter === 'brew' }" 
                @click="modalFilter = 'brew'"
              >🍺 Brew 服务</button>
              <button 
                :class="{ active: modalFilter === 'docker' }" 
                @click="modalFilter = 'docker'"
              >🐳 Docker 容器</button>
            </div>
          </div>

          <!-- Brew 服务列表 -->
          <div v-if="modalFilter !== 'docker'" class="lp-check-section">
            <div class="lp-section-title">
              Brew 服务
              <span class="lp-selected-count">{{ selectedBrew.length }} 已选</span>
            </div>
            <div v-if="filteredBrewServices().length === 0" class="lp-no-results">
              无匹配服务
            </div>
            <label 
              v-for="s in filteredBrewServices()" 
              :key="'brew-' + s.name"
              class="lp-check-item"
              :class="{ checked: isBrewSelected(s.name) }"
            >
              <input 
                type="checkbox" 
                :checked="isBrewSelected(s.name)"
                @change="toggleBrew(s.name)"
              />
              <span class="lp-check-name">{{ s.name }}</span>
              <span class="lp-check-status" :class="s.status">{{ brewStatusLabel(s.status) }}</span>
            </label>
          </div>

          <!-- Docker 容器列表 -->
          <div v-if="modalFilter !== 'brew'" class="lp-check-section">
            <div class="lp-section-title">
              Docker 容器
              <span class="lp-selected-count">{{ selectedDocker.length }} 已选</span>
            </div>
            <div v-if="filteredDockerContainers().length === 0" class="lp-no-results">
              {{ availableServices.docker_containers?.length === 0 ? '无容器，请确保 Docker 已启动' : '无匹配容器' }}
            </div>
            <label
              v-for="c in filteredDockerContainers()"
              :key="'docker-' + c.id"
              class="lp-check-item"
              :class="{ checked: isDockerSelected(c.name) }"
            >
              <input
                type="checkbox"
                :checked="isDockerSelected(c.name)"
                @change="toggleDocker(c.name)"
              />
              <span class="lp-check-name">{{ c.name }}</span>
              <span class="lp-check-meta">{{ c.image }}</span>
              <span class="lp-check-status" :class="c.state">{{ containerStateLabel(c.state) }}</span>
            </label>
          </div>
        </div>

        <div class="lp-modal-footer">
          <button class="lp-cancel-btn" @click="closeModal">取消</button>
          <button class="lp-save-btn" :disabled="saving" @click="handleSave">
            <span v-if="saving" class="btn-spinner"></span>
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.launchpad {
  margin-top: 32px;
}

.launchpad-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.launchpad-title {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
}

.add-group-btn {
  padding: 8px 16px;
  background: rgba(0, 122, 255, 0.2);
  border: 1px solid rgba(0, 122, 255, 0.3);
  border-radius: 8px;
  color: #007AFF;
  font-size: 0.8rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.add-group-btn:hover {
  background: rgba(0, 122, 255, 0.3);
}

.lp-loading,
.lp-empty {
  padding: 24px;
  text-align: center;
  color: rgba(255, 255, 255, 0.4);
  font-size: 0.875rem;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 8px;
}

.lp-groups {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.lp-group-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  padding: 14px 18px;
  transition: all 0.2s;
}

.lp-group-card:hover {
  background: rgba(255, 255, 255, 0.08);
}

.lp-group-card.running {
  border-left: 3px solid #34C759;
}

.lp-group-card.partial {
  border-left: 3px solid #FF9500;
}

.lp-group-card.stopped {
  border-left: 3px solid rgba(142, 142, 147, 0.3);
}

.lp-group-info {
  flex: 1;
  min-width: 0;
}

.lp-group-name-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.lp-group-name {
  font-size: 0.95rem;
  font-weight: 600;
}

.lp-status-badge {
  font-size: 0.7rem;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

.lp-status-badge.running {
  background: rgba(52, 199, 89, 0.15);
  color: #34C759;
}

.lp-status-badge.partial {
  background: rgba(255, 149, 0, 0.15);
  color: #FF9500;
}

.lp-status-badge.stopped {
  background: rgba(142, 142, 147, 0.2);
  color: #8E8E93;
}

.lp-status-badge.unknown {
  background: rgba(255, 255, 255, 0.06);
  color: rgba(255, 255, 255, 0.4);
}

.lp-group-meta {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.lp-meta-tag {
  font-size: 0.75rem;
  padding: 2px 8px;
  border-radius: 6px;
}

.lp-meta-tag.brew {
  background: rgba(255, 149, 0, 0.15);
  color: #FF9500;
}

.lp-meta-tag.docker {
  background: rgba(0, 122, 255, 0.15);
  color: #007AFF;
}

.lp-group-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.lp-launch-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: linear-gradient(135deg, #34C759, #30B350);
  border: none;
  border-radius: 8px;
  color: #fff;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.lp-launch-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(52, 199, 89, 0.3);
}

.lp-launch-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.lp-stop-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(255, 59, 48, 0.15);
  border: 1px solid rgba(255, 59, 48, 0.25);
  border-radius: 8px;
  color: #FF3B30;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.lp-stop-btn:hover:not(:disabled) {
  background: rgba(255, 59, 48, 0.25);
  transform: translateY(-1px);
}

.lp-stop-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.lp-restart-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(255, 149, 0, 0.15);
  border: 1px solid rgba(255, 149, 0, 0.25);
  border-radius: 8px;
  color: #FF9500;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.lp-restart-btn:hover:not(:disabled) {
  background: rgba(255, 149, 0, 0.25);
  transform: translateY(-1px);
}

.lp-restart-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.lp-edit-btn,
.lp-delete-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.35);
  transition: all 0.15s ease;
}

.lp-edit-btn:hover {
  color: #007AFF;
  background: rgba(0, 122, 255, 0.1);
}

.lp-delete-btn:hover {
  color: #FF3B30;
  background: rgba(255, 59, 48, 0.1);
}

.lp-delete-btn.confirming {
  width: auto;
  padding: 0 8px;
  gap: 3px;
  color: #FF3B30;
  background: rgba(255, 59, 48, 0.12);
  border: 1px solid rgba(255, 59, 48, 0.2);
}

.lp-delete-confirm-text {
  font-size: 0.7rem;
  font-weight: 500;
}

/* 启动结果 */
.lp-results {
  margin-top: 16px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  overflow: hidden;
}

.lp-results-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  font-size: 0.85rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.7);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.lp-results-close {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.4);
  cursor: pointer;
  font-size: 0.9rem;
}

.lp-result-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 16px;
  font-size: 0.85rem;
}

.lp-result-item.success {
  color: rgba(255, 255, 255, 0.8);
}

.lp-result-item.fail {
  color: #FF3B30;
}

.lp-result-icon {
  width: 16px;
  text-align: center;
}

.lp-result-type {
  font-size: 0.9rem;
}

.lp-result-name {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 0.8rem;
  min-width: 100px;
}

.lp-result-msg {
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.8rem;
}

/* Modal */
.lp-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.lp-modal {
  background: #1C1C1E;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 14px;
  width: 600px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}

.lp-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.lp-modal-header h4 {
  margin: 0;
  font-size: 1.05rem;
  font-weight: 600;
}

.lp-modal-close {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.4);
  cursor: pointer;
  font-size: 1rem;
  padding: 4px 8px;
  border-radius: 4px;
}

.lp-modal-close:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
}

.lp-modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
}

.lp-modal-body::-webkit-scrollbar {
  width: 4px;
}

.lp-modal-body::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

.lp-modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 14px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.lp-form-group {
  margin-bottom: 14px;
}

.lp-label {
  display: block;
  font-size: 0.8rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.6);
  margin-bottom: 6px;
}

.lp-input {
  width: 100%;
  padding: 10px 14px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: #fff;
  font-size: 0.875rem;
  outline: none;
  transition: border-color 0.2s;
}

.lp-input:focus {
  border-color: rgba(0, 122, 255, 0.4);
}

.lp-modal-search {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  align-items: center;
}

.lp-modal-search .lp-input {
  flex: 1;
  margin-bottom: 0;
}

.lp-filter-tabs {
  display: flex;
  gap: 4px;
}

.lp-filter-tabs button {
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.75rem;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.lp-filter-tabs button:hover {
  color: #fff;
}

.lp-filter-tabs button.active {
  background: rgba(0, 122, 255, 0.2);
  border-color: rgba(0, 122, 255, 0.3);
  color: #007AFF;
}

.lp-check-section {
  margin-bottom: 16px;
}

.lp-section-title {
  font-size: 0.85rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.7);
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.lp-selected-count {
  font-size: 0.75rem;
  font-weight: 400;
  color: rgba(255, 255, 255, 0.4);
}

.lp-no-results {
  padding: 16px;
  text-align: center;
  color: rgba(255, 255, 255, 0.3);
  font-size: 0.8rem;
}

.lp-check-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.15s;
}

.lp-check-item:hover {
  background: rgba(255, 255, 255, 0.04);
}

.lp-check-item.checked {
  background: rgba(0, 122, 255, 0.08);
}

.lp-check-item input[type="checkbox"] {
  accent-color: #007AFF;
  cursor: pointer;
  width: 16px;
  height: 16px;
}

.lp-check-name {
  flex: 1;
  font-size: 0.85rem;
  font-family: 'Monaco', 'Menlo', monospace;
  color: rgba(255, 255, 255, 0.85);
}

.lp-check-meta {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.4);
  max-width: 160px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.lp-check-status {
  font-size: 0.7rem;
  padding: 2px 8px;
  border-radius: 4px;
  min-width: 52px;
  text-align: center;
}

.lp-check-status.started,
.lp-check-status.running {
  background: rgba(52, 199, 89, 0.15);
  color: #34C759;
}

.lp-check-status.stopped,
.lp-check-status.exited {
  background: rgba(142, 142, 147, 0.2);
  color: #8E8E93;
}

.lp-check-status.paused {
  background: rgba(255, 149, 0, 0.15);
  color: #FF9500;
}

.lp-modal-footer button {
  padding: 9px 22px;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.lp-cancel-btn {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.7);
}

.lp-cancel-btn:hover {
  background: rgba(255, 255, 255, 0.12);
}

.lp-save-btn {
  background: #007AFF;
  border: none;
  color: #fff;
}

.lp-save-btn:hover:not(:disabled) {
  background: #0066D6;
}

.lp-save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Spinner */
.btn-spinner {
  display: inline-block;
  width: 12px;
  height: 12px;
  border: 2px solid transparent;
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
