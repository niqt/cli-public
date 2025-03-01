<script setup lang="ts">
import {ref, onMounted, computed, reactive} from 'vue'
import {
  ElTable,
  ElTableColumn,
  ElSlider,
  ElFormItem,
  ElRow,
  ElCol, ElMessage, ElMessageBox,
} from 'element-plus'
import {useDependencies} from '../stores/dependenciesStore.ts'
import type {DependenciesForAppDto, DependencyDto} from "@/model/dependeciesDto.ts";
import type {PackageDto} from "@/model/packageDto.ts";

let tableData = ref<DependencyDto[]>([]);
const store = useDependencies();
const pkgName = ref('github.com/cli/cli/v2')
const pkgVersion = ref('v2.66.1')
const scoreRange = ref([0,10])
const dialogFormVisible = ref(false)
const formLabelWidth = '140px'
const currentPage = ref(1);
const pageSize = ref(10);
const packageId = ref(0);
const searchText = ref('');
const loading = ref(false)
const dialogVisibleDelete = ref(false)
const depSelected = ref( 0)

const form = ref({
  id: 0,
  name: '',
  version: '',
  score: '',
  packageId: 0,
})

onMounted(async() => {
  try {
    const result = await store.getDependencies("", "0", "10") as DependenciesForAppDto;
    if (result.data ) {
      tableData.value.push(...result.data);
    }
    if (result.error ) {
      ElMessage.error('Oops, there is an error.')
      console.log(result.error)
    }
  } catch (error) {
    ElMessage.error('Oops, there is an error.')
  }

  try {
    const result = await store.getPackage(pkgName.value, pkgVersion.value) as PackageDto;
    if (result.data ) {
      packageId.value = result.data.id
    }
    if (result.error ) {
      ElMessage.error('Oops, there is an error.')
      console.log(result.error)
    }
  } catch (error) {
    ElMessage.error('Oops, there is an error.')
  }

});

const handleSearch = (async() => {
  try {
    const result = await store.getDependencies(searchText.value, scoreRange.value[0].toString(), scoreRange.value[1].toString()) as DependenciesForAppDto;
    tableData.value = [];
    if (result.data ) {
      tableData.value = result.data;
    }
    if (result.error ) {
      ElMessage.error('Oops, there is an error.')
      console.log(result.error)
    }
  } catch (error) {
    ElMessage.error('Oops, there is an error.')
  }
});

const handleDownload = (async() => {
  try {
    loading.value = true;
    const result = await store.downloadDependency() as DependenciesForAppDto;
    tableData.value = [];
    if (result.data ) {
      tableData.value = result.data;
    }
    if (result.error ) {
      ElMessage.error('Oops, there is an error.')
      console.log(result.error)
    }
  } catch (error) {
    ElMessage.error('Oops, there is an error.')
  }
  loading.value = false;
});

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return tableData.value.slice(start, end);
});

const handlePageChange = (page: any) => {
  currentPage.value = page;
};

const handleEdit = (dep: DependencyDto) => {
  form.value.id = dep.id;
  form.value.name = dep.name;
  form.value.version = dep.version;
  form.value.score = (dep.score == -1)? '': dep.score.toString();
  form.value.packageId = dep.packageId;
  dialogFormVisible.value = true;
}

const handleAdd = () => {
  form.value.id = 0;
  form.value.name = "";
  form.value.version = "";
  form.value.score = '';
  form.value.packageId = packageId.value;
  dialogFormVisible.value = true;
}

const handleSave = (async() => {
  try {
    const dep = {
      id: form.value.id,
      name: form.value.name,
      version: form.value.version,
      score: (form.value.score.length > 0)?parseFloat(form.value.score): -1,
      packageId: form.value.packageId,
    } as DependencyDto;
    let result: DependenciesForAppDto;
    if (dep.id == 0) {
      result = await store.createDependency(dep) as DependenciesForAppDto;
      if (result.data && result.data?.length > 0) {
        tableData.value.push(result.data[0]);
      }
    } else {
      result = await store.updateDependency(dep) as DependenciesForAppDto;
      if (result.data && result.data?.length > 0) {
        tableData.value = tableData.value.map(currentDep => currentDep.id === dep.id ? dep: currentDep);
      }
    }

    dialogFormVisible.value = false;
    if (result.data ) {
      console.log("Saved")
    }
    if (result.error ) {
      console.log(result.error)
    }
  } catch (error) {
  }
});

const handleDelete = (dep: DependencyDto) => {
  depSelected.value = dep.id;
  dialogVisibleDelete.value = true;
}

const handleCloseDelete = (async() => {
  try {
    const result = await store.deleteDependency(depSelected.value) as any;

    if (result.error ) {
      ElMessage.error('Oops, there is an error.')
      console.log(result.error)
    } else {
      tableData.value = tableData.value.filter(currentDep => currentDep.id !== depSelected.value);
    }
  } catch (error) {
    ElMessage.error('Oops, there is an error.')
  }
  dialogVisibleDelete.value = false;
});

</script>

<template>
  <el-form label-position="left"
           label-width="auto" style="width: 100%">

    <el-form-item label="Package name">
      <el-input v-model="pkgName" disabled/>
    </el-form-item>
    <el-form-item label="Package version">
      <el-input v-model="pkgVersion" disabled/>
    </el-form-item>
    <el-row>
      <el-col :span="10">
        <el-divider/>
      </el-col>
      <el-col :span="4"
              style="text-align: center; display: flex; justify-content: center; align-items: center;">
        <el-text class="mx-1">Filter</el-text>
      </el-col>
      <el-col :span="10">
        <el-divider/>
      </el-col>
    </el-row>
    <el-form-item label="Dependency name">
      <el-input v-model="searchText"/>
    </el-form-item>
    <el-form-item label="Score Range">
      <el-slider v-model="scoreRange" range show-stops :max="10"/>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="handleSearch()" :disabled="loading">Search</el-button>
    </el-form-item>
    <el-divider/>
    <el-row :gutter="20" class="button-row">
      <el-col :span="12" class="button-container">
        <el-form-item>
          <el-button type="primary" @click="handleDownload" :disabled="loading">Load dependencies</el-button>
        </el-form-item>
      </el-col>
      <el-col :span="12" class="button-container">
        <el-form-item>
        <el-button type="primary" @click="handleAdd()" :disabled="loading">Add Dependency</el-button>
      </el-form-item>
      </el-col>
    </el-row>
    <div>
    <el-table v-loading="loading" :data="paginatedData" :border="true"
              :stripe="true" width="100%">
      <el-table-column prop="name" label="Name" min-width="350px"/>
      <el-table-column prop="version" label="Version" min-width="200px"/>
      <el-table-column label="score" min-width="80px">
        <template #default="scope">
          <span>{{ scope.row.score.toString() === '-1' ? '-' : scope.row.score }}</span>
        </template>
      </el-table-column>
      <el-table-column fixed="right" label="Operations" width="120px">
        <template #default="scope">
          <el-button link type="primary" size="small" @click="handleEdit(scope.row)">Edit
          </el-button>
          <el-button link type="danger" size="small" @click="handleDelete(scope.row)">Delete
          </el-button>
        </template>
      </el-table-column>
    </el-table>
      <br>
      <el-pagination
        background
        layout="prev, pager, next"
        :total="tableData.length"
        :page-size="pageSize"
        :current-page.sync="currentPage"
        @current-change="handlePageChange"
      />
    </div>
  </el-form>
  <el-dialog v-model="dialogFormVisible" title="Edit Dependency" width="500">
    <el-form :model="form">
      <el-form-item label="Name *" :label-width="formLabelWidth">
        <el-input v-model="form.name" autocomplete="off"/>
      </el-form-item>
      <el-form-item label="Version *" :label-width="formLabelWidth">
       <el-input v-model="form.version" autocomplete="off"/>
      </el-form-item>
      <el-form-item label="Score" :label-width="formLabelWidth">
        <el-input v-model="form.score" autocomplete="off"/>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogFormVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleSave">
          Save
        </el-button>
      </div>
    </template>
  </el-dialog>
  <el-dialog
    v-model="dialogVisibleDelete"
    title="Delete"
    width="500"
  >
    <span>Do you want delete?</span>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisibleDelete = false">Cancel</el-button>
        <el-button type="primary" @click="handleCloseDelete">
          Confirm
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<style scoped>
.button-row {
  display: flex;
  align-items: center;
}

.button-container {
  display: flex;
  justify-content: center;
}
</style>
