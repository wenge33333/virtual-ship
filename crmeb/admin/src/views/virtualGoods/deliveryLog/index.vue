<template>
  <div class="divBox relative">
    <el-card class="box-card" shadow="never">
      <div class="clearfix">
        <div class="filter-container flex-between">
          <el-form inline size="small">
            <el-form-item label="订单号：">
              <el-input v-model="listPram.order_id" placeholder="请输入订单号" style="width:220px" clearable />
            </el-form-item>
            <el-form-item label="状态：">
              <el-select v-model="listPram.status" clearable placeholder="全部" style="width:120px">
                <el-option :value="1" label="成功" />
                <el-option :value="2" label="失败" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="el-icon-search" @click="handleSearch">查询</el-button>
              <el-button icon="el-icon-refresh" @click="handleReset">重置</el-button>
              <el-button type="danger" icon="el-icon-delete" @click="handleBatchDelete" :disabled="selectedLogs.length===0">批量删除</el-button>
            </el-form-item>
          </el-form>
          <el-button type="primary" icon="el-icon-download" @click="handleExport">导出 CSV</el-button>
        </div>
      </div>
      <el-table v-loading="loading" :data="listData.list" style="width:100%" size="small" class="mt14" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="order_id" label="订单号" min-width="160" :show-overflow-tooltip="true" />
        <el-table-column prop="product_id" label="商品ID" width="80" />
        <el-table-column label="发货类型" width="100">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.delivery_type === 1" type="success" size="small">自动发卡密</el-tag>
            <el-tag v-else-if="scope.row.delivery_type === 2" type="warning" size="small">手动处理</el-tag>
            <el-tag v-else type="info" size="small">固定内容</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="delivery_content" label="发货内容" min-width="200" :show-overflow-tooltip="true" />
        <el-table-column label="状态" width="80">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.status === 1" type="success" size="small">成功</el-tag>
            <el-tag v-else type="danger" size="small">失败</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="retry_count" label="重试次数" width="80" />
        <el-table-column prop="error_msg" label="错误信息" min-width="180" :show-overflow-tooltip="true" />
        <el-table-column prop="platform_source" label="来源平台" width="100" />
        <el-table-column label="创建时间" width="160">
          <template slot-scope="scope">
            <span>{{ scope.row.create_time }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80" fixed="right">
          <template slot-scope="scope">
            <el-button v-if="scope.row.status === 2 && scope.row.retry_count < 3" type="text" size="mini" @click="handleRetry(scope.row)">重试</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="mt14"
        :current-page="listPram.page"
        :page-sizes="constants.page.limit"
        :layout="constants.page.layout"
        :total="listData.total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        background
      />
    </el-card>
  </div>
</template>

<script>
import { deliveryLogsApi, deliveryRetryApi, deleteDeliveryLogsApi, exportDeliveryLogsApi } from '@/api/virtualGoods';

export default {
  data() {
    return {
      constants: this.$constants,
      listPram: {
        order_id: '',
        status: null,
        page: 1,
        limit: this.$constants.page.limit[0],
      },
      listData: { list: [], total: 0 },
      loading: false,
      selectedLogs: [],
    };
  },
  mounted() {
    this.getList();
  },
  methods: {
    handleSelectionChange(val) {
      this.selectedLogs = val;
    },
    getList() {
      this.loading = true;
      const params = {
        page: this.listPram.page,
        limit: this.listPram.limit,
      };
      if (this.listPram.order_id) params.order_id = this.listPram.order_id;
      if (this.listPram.status) params.status = this.listPram.status;
      deliveryLogsApi(params)
        .then((data) => {
          this.listData = data;
          this.loading = false;
        })
        .catch(() => {
          this.listData = { list: [], total: 0 };
          this.loading = false;
        });
    },
    handleSearch() {
      this.listPram.page = 1;
      this.getList();
    },
    handleReset() {
      this.listPram.order_id = '';
      this.listPram.status = null;
      this.listPram.page = 1;
      this.getList();
    },
    handleBatchDelete() {
      if (this.selectedLogs.length === 0) {
        this.$message.warning('请选择要删除的日志');
        return;
      }
      this.$confirm(`确认删除选中的 ${this.selectedLogs.length} 条发货日志吗？`, '提示', { type: 'warning' }).then(() => {
        const ids = this.selectedLogs.map(item => item.id);
        deleteDeliveryLogsApi(ids)
          .then((res) => {
            this.$message.success(`成功删除 ${res.deleted} 条日志`);
            this.getList();
          })
          .catch(() => {});
      }).catch(() => {});
    },
    handleExport() {
      const params = { page: 1, limit: 1000 };
      if (this.listPram.order_id) params.order_id = this.listPram.order_id;
      if (this.listPram.status) params.status = this.listPram.status;
      exportDeliveryLogsApi(params).then((blob) => {
        const link = document.createElement('a');
        link.href = window.URL.createObjectURL(new Blob([blob]));
        link.download = `delivery_logs_${new Date().getTime()}.csv`;
        link.click();
        window.URL.revokeObjectURL(link.href);
        this.$message.success('导出成功');
      }).catch(() => {});
    },
    handleSizeChange(val) {
      this.listPram.limit = val;
      this.getList();
    },
    handleCurrentChange(val) {
      this.listPram.page = val;
      this.getList();
    },
    handleRetry(row) {
      this.$confirm('确认重新执行发货？', '提示', { type: 'warning' }).then(() => {
        deliveryRetryApi({ order_id: row.order_id })
          .then(() => {
            this.$message.success('重试成功');
            this.getList();
          })
          .catch(() => {});
      }).catch(() => {});
    },
  },
};
</script>

<style scoped>
.filter-container {
  display: flex;
  align-items: center;
}
.flex-between {
  justify-content: space-between;
}
.mt14 {
  margin-top: 14px;
}
</style>
