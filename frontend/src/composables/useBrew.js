import { reactive, ref, computed, onMounted, onUnmounted } from 'vue'
import { GetBrewData, StartService, StopService, GetAppIcon, RestartService, AddTap, RemoveTap, UpdateTap, UpdateAllTaps, SearchPackages, InstallPackage } from '../../wailsjs/go/main/App'

// 自动刷新间隔（2分钟）
const AUTO_REFRESH_INTERVAL = 120000

export function useBrew() {
  const data = reactive({ formulae: [], casks: [], taps: [], loading: false })
  const currentTab = ref('overview') // 当前 tab
  const searchQuery = ref('')
  const processingMap = reactive(new Map())
  const toast = reactive({ show: false, msg: '', type: 'success' })
  
  // 自动刷新进度条状态
  const refreshProgress = reactive({
    isRefreshing: false,    // 是否正在刷新
    progress: 0,           // 进度百分比 (0-100)
    lastRefreshTime: null   // 上次刷新时间
  })
  
  let refreshTimer = null

  // 统一的提示函数
  function showToast(msg, type = 'success') {
    toast.msg = msg
    toast.type = type
    toast.show = true
    setTimeout(() => { toast.show = false }, 3000)
  }

  // 设置当前 tab
  function setCurrentTab(tab) {
    currentTab.value = tab
  }

  // 搜索过滤逻辑
  const filteredFormulae = computed(() => 
    data.formulae.filter(item => item.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
  )

  const filteredCasks = computed(() => 
    data.casks.filter(item => item.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
  )

  // 全量刷新（启动时用）
  async function updateList() {
    try {
      data.loading = true
      refreshProgress.isRefreshing = true
      
      const res = await GetBrewData()
      // 遍历 Casks 为每个应用请求图标
      const caskWithIcons = await Promise.all(res.casks.map(async (item) => {
        const icon = await GetAppIcon(item.name)
        return { ...item, iconBase64: icon }
      }))
      data.casks = caskWithIcons
      data.formulae = res.formulae
      data.taps = res.taps || []
      
      // 更新刷新时间
      refreshProgress.lastRefreshTime = Date.now()
      refreshProgress.progress = 0
    } catch (err) {
      console.error("刷新失败:", err)
    } finally {
      data.loading = false
      refreshProgress.isRefreshing = false
    }
  }

  // 只刷新当前 tab 的数据
  async function refreshCurrentTab() {
    try {
      refreshProgress.isRefreshing = true
      
      if (currentTab.value === 'gui') {
        await refreshCasks()
      } else if (currentTab.value === 'tui') {
        await refreshFormulae()
      } else if (currentTab.value === 'taps') {
        await refreshTaps()
      } else {
        // overview 需要全部数据
        await updateList()
      }
      
      // 更新刷新时间
      refreshProgress.lastRefreshTime = Date.now()
      refreshProgress.progress = 0
    } finally {
      refreshProgress.isRefreshing = false
    }
  }

  // 只刷新 casks (GUI 应用)
  async function refreshCasks() {
    try {
      const res = await GetBrewData()
      const caskWithIcons = await Promise.all(res.casks.map(async (item) => {
        const icon = await GetAppIcon(item.name)
        return { ...item, iconBase64: icon }
      }))
      data.casks = caskWithIcons
    } catch (err) {
      console.error("刷新 Casks 失败:", err)
    }
  }

  // 只刷新 formulae (命令行工具)
  async function refreshFormulae() {
    try {
      const res = await GetBrewData()
      data.formulae = res.formulae
    } catch (err) {
      console.error("刷新 Formulae 失败:", err)
    }
  }

  // 只刷新 taps
  async function refreshTaps() {
    try {
      const res = await GetBrewData()
      data.taps = res.taps || []
    } catch (err) {
      console.error("刷新 Taps 失败:", err)
    }
  }

  // 启动定时刷新（每 2 分钟刷新当前 tab）
  function startAutoRefresh() {
    if (refreshTimer) return
    
    // 初始化刷新时间
    refreshProgress.lastRefreshTime = Date.now()
    
    // 启动进度条更新定时器（每秒更新一次）
    const progressTimer = setInterval(() => {
      if (refreshProgress.lastRefreshTime && !refreshProgress.isRefreshing) {
        const elapsed = Date.now() - refreshProgress.lastRefreshTime
        const progress = Math.min((elapsed / AUTO_REFRESH_INTERVAL) * 100, 100)
        refreshProgress.progress = progress
      }
    }, 1000)
    
    // 启动自动刷新定时器
    refreshTimer = setInterval(() => {
      refreshCurrentTab()
    }, AUTO_REFRESH_INTERVAL)
    
    // 将进度条定时器存储到 refreshTimer 对象中以便清理
    refreshTimer.progressTimer = progressTimer
  }

  // 停止定时刷新
  function stopAutoRefresh() {
    if (refreshTimer) {
      // 清除进度条定时器
      if (refreshTimer.progressTimer) {
        clearInterval(refreshTimer.progressTimer)
      }
      clearInterval(refreshTimer)
      refreshTimer = null
    }
    // 重置进度状态
    refreshProgress.progress = 0
    refreshProgress.lastRefreshTime = null
  }

  // 服务启动/停止逻辑
  async function handleService(item) {
    processingMap.set(item.name, true)
    try {
      let result = item.status === 'started' ? await StopService(item.name) : await StartService(item.name)
      if (result.success) {
        showToast(result.message, 'success')
        await refreshCurrentTab() // 操作后刷新当前 tab
      } else {
        showToast(result.message, 'error')
      }
    } catch (err) {
      showToast("系统错误: " + err, 'error')
    } finally {
      processingMap.delete(item.name)
    }
  }

  async function handleRestart(name) {
    if (processingMap.has(name)) return
    processingMap.set(name, true)
    try {
      const result = await RestartService(name)
      showToast(result.success ? result.message : "系统错误: " + result.message, result.success ? 'success' : 'error')
      if (result.success) {
        await refreshCurrentTab()
      }
    } catch (err) {
      showToast("系统异常: " + err, 'error')
    } finally {
      processingMap.delete(name)
    }
  }

  // Tap 管理方法
  async function handleAddTap(tapName) {
    const key = `tap-add-${tapName}`
    if (processingMap.has(key)) return
    processingMap.set(key, true)
    
    try {
      const result = await AddTap(tapName)
      showToast(result.message, result.success ? 'success' : 'error')
      if (result.success) {
        await refreshTaps()
      }
    } catch (err) {
      showToast("添加 tap 失败: " + err, 'error')
    } finally {
      processingMap.delete(key)
    }
  }

  async function handleRemoveTap(tapName, force = false) {
    const key = `tap-remove-${tapName}`
    if (processingMap.has(key)) return
    processingMap.set(key, true)
    
    try {
      const result = await RemoveTap(tapName, force)
      showToast(result.message, result.success ? 'success' : 'error')
      if (result.success) {
        await refreshTaps()
      }
    } catch (err) {
      showToast("移除 tap 失败: " + err, 'error')
    } finally {
      processingMap.delete(key)
    }
  }

  async function handleUpdateTap(tapName) {
    const key = `tap-update-${tapName}`
    if (processingMap.has(key)) return
    processingMap.set(key, true)
    
    try {
      const result = await UpdateTap(tapName)
      showToast(result.message, result.success ? 'success' : 'error')
    } catch (err) {
      showToast("更新 tap 失败: " + err, 'error')
    } finally {
      processingMap.delete(key)
    }
  }

  async function handleUpdateAllTaps() {
    const key = 'tap-update-all'
    if (processingMap.has(key)) return
    processingMap.set(key, true)
    
    try {
      const result = await UpdateAllTaps()
      showToast(result.message, result.success ? 'success' : 'error')
      if (result.success) {
        await refreshTaps()
      }
    } catch (err) {
      showToast("更新所有 taps 失败: " + err, 'error')
    } finally {
      processingMap.delete(key)
    }
  }

  // 统计信息
  const stats = computed(() => ({
    gui: data.casks.length,
    tui: data.formulae.length,
    taps: data.taps.length
  }))

  return {
    data, 
    currentTab, 
    setCurrentTab, 
    searchQuery, 
    processingMap, 
    toast, 
    stats,
    refreshProgress, // 自动刷新进度状态
    filteredFormulae, 
    filteredCasks,
    updateList, 
    refreshCurrentTab, 
    refreshCasks, 
    refreshFormulae, 
    refreshTaps,
    startAutoRefresh, 
    stopAutoRefresh,
    handleService, 
    handleRestart,
    showToast,
    handleAddTap, 
    handleRemoveTap, 
    handleUpdateTap, 
    handleUpdateAllTaps
  }
}
