<script setup>
import {useDependencies} from '../stores/dependenciesStore.ts'
import {
  ElRow,
  ElCol,
} from 'element-plus'

const store = useDependencies();
const getPercentage = (value) => {
  return Math.min(Math.max(value * 10, 0), 100);
};
const depOrdered = (store.dependencies)?store.dependencies.sort((a, b) => {return a.score < b.score}): []
const getColorObject = (value) => {
  if (value <= 3) return '#ff6b6b'; // rosso per valori bassi
  if (value <= 6) return '#ffd166'; // giallo per valori medi
  return '#06d6a0'; // verde per valori alti
};

</script>

<template>
  <div class="bar-chart-container">
    <h2>Dependecnies Chart</h2>
    <div class="chart">
      <el-row v-for="(item, index) in depOrdered" :key="index" class="chart-row">
        <el-col :span="8" class="item-name">{{ item.name }}</el-col>
        <el-col :span="16">
          <div class="chart-progress">
            <el-progress
              :percentage="getPercentage(item.score)"
              :color="getColorObject(item.score)"
              :format="() => item.score >= 0? item.score.toString(): '' "
              :stroke-width="20"
            />
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>



<style scoped>
.bar-chart-container {
  font-family: Arial, sans-serif;
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.chart {
  margin-bottom: 20px;
}

.chart-row {
  margin-bottom: 20px;
  align-items: center;
}

.item-name {
  font-weight: bold;
  text-align: right;
  padding-right: 15px;
}

.chart-progress {
  padding-right: 15px;
}

:deep(.el-progress-bar__outer) {
  border-radius: 4px;
}

:deep(.el-progress-bar__inner) {
  border-radius: 4px;
}

:deep(.el-progress__text) {
  font-weight: bold;
}
</style>
