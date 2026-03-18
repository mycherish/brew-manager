<script setup>
import { ref, computed } from 'vue'
import ItemCard from '../components/ItemCard.vue'

const props = defineProps({
  formulae: {
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

const filteredFormulae = computed(() => {
  if (!localSearch.value) return props.formulae
  const query = localSearch.value.toLowerCase()
  return props.formulae.filter(formula => 
    formula.name.toLowerCase().includes(query)
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
  <div class="tui-tools-page">
    <div class="page-header">
      <h2 class="page-title">命令行工具</h2>
      <span class="count-badge">{{ filteredFormulae.length }} 个工具</span>
    </div>
    
    <!-- 本地搜索 -->
    <div class="local-search">
      <span class="search-icon">🔍</span>
      <input 
        v-model="localSearch" 
        type="text" 
        placeholder="搜索已安装的工具..." 
        class="search-input"
      />
    </div>

    <!-- 工具列表 -->
    <div class="tools-list">
      <ItemCard 
        v-for="item in filteredFormulae" 
        :key="item.name"
        type="formula"
        :item="item"
        :is-processing="processingMap.has(item.name)"
        @action="handleAction"
        @restart="handleRestart"
      />
      
      <div v-if="filteredFormulae.length === 0" class="empty-state">
        <div class="empty-icon">💻</div>
        <p v-if="formulae.length === 0">暂无安装的命令行工具</p>
        <p v-else>未找到匹配的工具</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tui-tools-page {
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

.tools-list {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tools-list::-webkit-scrollbar {
  width: 5px;
}

.tools-list::-webkit-scrollbar-thumb {
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
