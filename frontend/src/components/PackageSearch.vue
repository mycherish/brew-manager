<script setup>
import { ref, reactive } from 'vue'
import { SearchPackages, InstallPackage } from '../../wailsjs/go/main/App'

// Emits
const emit = defineEmits(['install-success', 'add-tap'])

// Local state
const searchQuery = ref('')
const searchResults = reactive({
  formulae: [],
  casks: [],
  taps: [],
  total: 0
})
const isSearching = ref(false)
const hasSearched = ref(false)
const installing = reactive({})

// Methods
async function handleSearch() {
  if (!searchQuery.value.trim()) {
    return
  }
  
  isSearching.value = true
  hasSearched.value = false
  
  try {
    const results = await SearchPackages(searchQuery.value.trim())
    searchResults.formulae = results.formulae || []
    searchResults.casks = results.casks || []
    searchResults.taps = results.taps || []
    searchResults.total = results.total || 0
    hasSearched.value = true
  } catch (err) {
    console.error('搜索失败:', err)
  } finally {
    isSearching.value = false
  }
}

async function handleInstall(name, isCask) {
  if (installing[name]) return
  
  installing[name] = true
  
  try {
    const result = await InstallPackage(name, isCask)
    if (result.success) {
      emit('install-success', name)
      alert(result.message)
    } else {
      alert(result.message)
    }
  } catch (err) {
    alert('安装失败: ' + err)
  } finally {
    delete installing[name]
  }
}

function handleKeyup(e) {
  if (e.key === 'Enter') {
    handleSearch()
  }
}
</script>

<template>
  <div class="package-search">
    <div class="search-header">
      <h3 class="search-title">🔍 搜索软件包</h3>
      <p class="search-desc">搜索 Homebrew 官方仓库中的软件包</p>
    </div>
    
    <div class="search-box">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="输入软件包名称，如：nginx、firefox、node..."
        class="search-input"
        @keyup.enter="handleSearch"
      />
      <button 
        @click="handleSearch" 
        class="btn btn-primary"
        :disabled="isSearching || !searchQuery.trim()"
      >
        <span v-if="isSearching" class="mini-loader"></span>
        <span v-else>搜索</span>
      </button>
    </div>
    
    <!-- 搜索结果 -->
    <div v-if="hasSearched" class="search-results">
      <div v-if="searchResults.total === 0" class="empty-result">
        <div class="empty-icon">🔍</div>
        <p>未找到匹配的软件包</p>
        <p class="hint">尝试使用其他关键词搜索</p>
      </div>
      
      <div v-else class="results-content">
        <div class="results-summary">
          找到 {{ searchResults.total }} 个结果
          <span v-if="searchResults.formulae.length > 0">
            ({{ searchResults.formulae.length }} 个 formulae
          </span>
          <span v-if="searchResults.casks.length > 0">
            {{ searchResults.formulae.length > 0 ? ',' : '(' }}
            {{ searchResults.casks.length }} 个 casks
          </span>
          <span v-if="searchResults.taps.length > 0">
            {{ searchResults.formulae.length > 0 || searchResults.casks.length > 0 ? ',' : '(' }}
            {{ searchResults.taps.length }} 个 taps)
          </span>
          <span v-else-if="searchResults.formulae.length > 0 || searchResults.casks.length > 0">
            )
          </span>
        </div>
        
        <!-- Formulae 结果 -->
        <div v-if="searchResults.formulae.length > 0" class="result-section">
          <h4 class="section-title">
            <span class="badge badge-formula">Formula</span>
            命令行工具 ({{ searchResults.formulae.length }})
          </h4>
          <div class="package-list">
            <div 
              v-for="pkg in searchResults.formulae.slice(0, 10)" 
              :key="pkg"
              class="package-item"
            >
              <div class="package-info">
                <span class="package-name">{{ pkg }}</span>
                <span class="package-type">formula</span>
              </div>
              <button 
                @click="handleInstall(pkg, false)"
                class="btn-install"
                :disabled="installing[pkg]"
              >
                <span v-if="installing[pkg]" class="mini-loader"></span>
                <span v-else>安装</span>
              </button>
            </div>
            <div v-if="searchResults.formulae.length > 10" class="more-hint">
              还有 {{ searchResults.formulae.length - 10 }} 个结果...
            </div>
          </div>
        </div>
        
        <!-- Casks 结果 -->
        <div v-if="searchResults.casks.length > 0" class="result-section">
          <h4 class="section-title">
            <span class="badge badge-cask">Cask</span>
            桌面应用 ({{ searchResults.casks.length }})
          </h4>
          <div class="package-list">
            <div 
              v-for="pkg in searchResults.casks.slice(0, 10)" 
              :key="pkg"
              class="package-item"
            >
              <div class="package-info">
                <span class="package-name">{{ pkg }}</span>
                <span class="package-type">cask</span>
              </div>
              <button 
                @click="handleInstall(pkg, true)"
                class="btn-install"
                :disabled="installing[pkg]"
              >
                <span v-if="installing[pkg]" class="mini-loader"></span>
                <span v-else>安装</span>
              </button>
            </div>
            <div v-if="searchResults.casks.length > 10" class="more-hint">
              还有 {{ searchResults.casks.length - 10 }} 个结果...
            </div>
          </div>
        </div>
        
        <!-- Taps 结果 -->
        <div v-if="searchResults.taps.length > 0" class="result-section">
          <h4 class="section-title">
            <span class="badge badge-tap">Tap</span>
            第三方仓库 ({{ searchResults.taps.length }})
          </h4>
          <div class="package-list">
            <div 
              v-for="tap in searchResults.taps" 
              :key="tap.name"
              class="package-item tap-item"
            >
              <div class="package-info tap-info">
                <span class="package-name">{{ tap.name }}</span>
                <span class="tap-stars">⭐ {{ tap.stars }}</span>
                <p v-if="tap.description" class="tap-description">{{ tap.description }}</p>
              </div>
              <button 
                @click="$emit('add-tap', tap.name)"
                class="btn-install btn-add-tap"
              >
                添加
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.package-search {
  background: rgba(255, 255, 255, 0.7);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.search-header {
  margin-bottom: 16px;
}

.search-title {
  margin: 0 0 4px 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #333;
}

.search-desc {
  margin: 0;
  font-size: 0.875rem;
  color: #718096;
}

.search-box {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.search-input {
  flex: 1;
  padding: 10px 14px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 0.875rem;
  transition: border-color 0.2s ease;
}

.search-input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.1);
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

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.btn-primary {
  background: linear-gradient(135deg, #3498db, #2980b9);
  color: white;
  box-shadow: 0 2px 6px rgba(52, 152, 219, 0.3);
}

.btn-primary:hover {
  box-shadow: 0 4px 12px rgba(52, 152, 219, 0.4);
}

.search-results {
  margin-top: 20px;
  max-height: 400px;
  overflow-y: auto;
  padding-right: 8px;
}

.search-results::-webkit-scrollbar {
  width: 5px;
}

.search-results::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 10px;
}

.empty-result {
  text-align: center;
  padding: 40px 20px;
  color: #718096;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-result p {
  margin: 0 0 8px 0;
}

.empty-result .hint {
  font-size: 0.875rem;
  opacity: 0.7;
}

.results-summary {
  margin-bottom: 16px;
  padding: 10px 14px;
  background: rgba(52, 152, 219, 0.1);
  border-radius: 8px;
  color: #3498db;
  font-weight: 500;
}

.result-section {
  margin-bottom: 20px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0 0 12px 0;
  font-size: 1rem;
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

.badge-formula {
  background: rgba(46, 204, 113, 0.1);
  color: #2ecc71;
}

.badge-cask {
  background: rgba(155, 89, 182, 0.1);
  color: #9b59b6;
}

.badge-tap {
  background: rgba(241, 196, 15, 0.1);
  color: #f39c12;
}

.package-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.package-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: white;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  transition: all 0.2s ease;
}

.package-item:hover {
  border-color: #3498db;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.package-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.package-name {
  font-weight: 600;
  color: #333;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.package-type {
  font-size: 0.75rem;
  color: #718096;
  text-transform: uppercase;
}

.btn-install {
  padding: 6px 14px;
  border-radius: 6px;
  border: none;
  background: linear-gradient(135deg, #2ecc71, #27ae60);
  color: white;
  font-size: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 4px rgba(46, 204, 113, 0.3);
}

.btn-install:hover {
  box-shadow: 0 4px 8px rgba(46, 204, 113, 0.4);
  transform: translateY(-1px);
}

.btn-install:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.more-hint {
  text-align: center;
  padding: 10px;
  color: #718096;
  font-size: 0.875rem;
  font-style: italic;
}

.tap-item {
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
}

.tap-info {
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}

.tap-stars {
  font-size: 0.75rem;
  color: #f39c12;
}

.tap-description {
  margin: 0;
  font-size: 0.75rem;
  color: #718096;
  line-height: 1.4;
}

.btn-add-tap {
  background: linear-gradient(135deg, #f39c12, #e67e22);
  box-shadow: 0 2px 4px rgba(243, 156, 18, 0.3);
}

.btn-add-tap:hover {
  box-shadow: 0 4px 8px rgba(243, 156, 18, 0.4);
}

.mini-loader {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>