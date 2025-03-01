<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import {ElHeader, ElRow, ElCol} from "element-plus";
</script>

  <template>
    <el-row>
      <el-col :span="2">
      <RouterLink to="/">
        <el-button>Home</el-button></RouterLink>
      </el-col>
      <el-col :span="2">
        <RouterLink to="/chart"><el-button>Chart</el-button></RouterLink>
        </el-col>
    </el-row>
      <el-container  style="width: 80%;height: 100%; display:block; margin: auto">
        <el-header style="width: 80%;"><h1>Dependency for the specified Package</h1></el-header>
        <el-container width="300px">

          <RouterView/>
        </el-container>
      </el-container>
  </template>

<style scoped>
</style>
