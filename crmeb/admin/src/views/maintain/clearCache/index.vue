<template>
  <div class="divBox">
    <el-card class="box-card">
      <el-button type="primary" v-hasPermi="['admin:system:config:clear:cache']" v-debounceClick="handleClear"
        >清除缓存</el-button
      >
    </el-card>
  </div>
</template>

<script>
// +---------------------------------------------------------------------
// +---------------------------------------------------------------------
// +---------------------------------------------------------------------
// +---------------------------------------------------------------------
// +---------------------------------------------------------------------
import { clearCacheApi } from '@/api/systemConfig';
export default {
  name: 'clearCache',
  methods: {
    handleClear() {
      clearCacheApi().then((response) => {
        this.$message.success('清除成功');
      });
    },
  },
};
</script>

<style scoped></style>
