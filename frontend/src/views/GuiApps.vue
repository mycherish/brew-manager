<script setup>
import { ref, computed } from 'vue'
import ItemCard from '../components/ItemCard.vue'

const props = defineProps({
  casks: {
    type: Array,
    default: () => []
  },
  processingMap: {
    type: Map,
    default: () => new Map()
  }
})

const emit = defineEmits(['action', 'restart'])

const localSearch = ref('')

const filteredCasks = computed(() => {
  if (!localSearch.value) return props.casks
  const query = localSearch.value.toLowerCase()
  return props.casks.filter(cask => 
    cask.name.toLowerCase().includes(query)
  )
})

function handleAction(item) {
  emit('action', item)
}

function handleRestart(name) {
  emit('restart', name)
}
</script>

<template>
  <div class="gui-apps-page">
    <div class="page-header">
      <h2 class="page-title">GUI 应用</h2>
      <span class="count-badge">{{ filteredCasks.length }} 个应用</span>
    </div>
    
    <!-- 本地搜索 -->
    <div class="local-search">
      <span class="search-icon">🔍</span>
      <input 
        v-model="localSearch" 
        type="text" 
        placeholder="搜索已安装的应用..." 
        class="search-input"
      />
    </div>

    <!-- 应用列表 -->
    <div class="apps-list">
      <ItemCard 
        v-for="item in filteredCasks" 
        :key="item.name"
        type="cask"
        :item="item"
        :is-processing="processingMap.has(item.name)"
        @action="handleAction"
        @restart="handleRestart"
      />
      
      <div v-if="filteredCasks.length === 0" class="empty-state">
        <div class="empty-icon">🖥️</div>
        <p v-if="casks.length === 0">暂无安装的 GUI 应用</p>
        <p v-else>未找到匹配的应用</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.gui-apps-page {
  padding: 20px;
  color: #fff;
  height: 100%;
  display: flex;
  flex-direction: column;
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

.count-badge {
  background: rgba(255, 255, 255, 0.1);
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.7);
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

.apps-list {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.apps-list::-webkit-scrollbar {
  width: 5px;
}

.apps-list::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
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
  margin: 0;
  font-size: 0.875rem;
}
</style>
