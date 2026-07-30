<template>
  <div class="divBox">
    <el-card shadow="never" class="mt14">
      <div slot="header" class="clearfix flex-between">
        <span>数据库备份</span>
        <el-button type="primary" icon="el-icon-circle-plus-outline" @click="handleBackup">立即备份</el-button>
      </div>

      <el-table v-loading="loading" :data="backupList" style="width:100%" size="small">
        <el-table-column prop="filename" label="文件名" min-width="200" />
        <el-table-column prop="size" label="文件大小" width="120">
          <template slot-scope="scope">
            <span>{{ formatSize(scope.row.size) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="time" label="备份时间" width="180" />
        <el-table-column label="操作" width="220" fixed="right">
          <template slot-scope="scope">
            <el-button type="success" size="mini" icon="el-icon-download" @click="handleDownload(scope.row)">下载</el-button>
            <el-button type="warning" size="mini" icon="el-icon-refresh" @click="handleRestore(scope.row)">恢复</el-button>
            <el-button type="danger" size="mini" icon="el-icon-delete" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="backupList.length === 0 && !loading" class="empty-tip text-center">
        <p>暂无备份记录</p>
      </div>
    </el-card>

    <!-- 恢复确认对话框 -->
    <el-dialog title="确认恢复" :visible.sync="restoreDialogVisible" width="450px" :close-on-click-modal="false">
      <el-alert
        title="警告：恢复操作将覆盖当前数据库"
        type="warning"
        :closable="false"
        show-icon
        class="mb20"
      >
        <p>此操作不可逆，请谨慎操作！</p>
        <p>确定要恢复备份文件 <b>{{ selectedBackup?.filename }}</b> 吗？</p>
      </el-alert>
      <span slot="footer" class="dialog-footer">
        <el-button @click="restoreDialogVisible = false" size="small">取消</el-button>
        <el-button type="primary" @click="confirmRestore" size="small" :loading="restoreLoading">确定恢复</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { listBackupsApi, backupDatabaseApi, restoreDatabaseApi, downloadBackupApi } from '@/api/virtualGoods';

export default {
  data() {
    return {
      loading: false,
      restoreLoading: false,
      backupList: [],
      restoreDialogVisible: false,
      selectedBackup: null,
    };
  },
  mounted() {
    this.getList();
  },
  methods: {
    getList() {
      this.loading = true;
      listBackupsApi()
        .then((data) => {
          this.backupList = data || [];
          this.loading = false;
        })
        .catch(() => {
          this.loading = false;
        });
    },
    handleBackup() {
      this.$confirm('确定要备份数据库吗？', '提示', { type: 'info' }).then(() => {
        backupDatabaseApi()
          .then((res) => {
            this.$message.success(`备份成功：${res.filename} (${this.formatSize(res.size)})`);
            this.getList();
          })
          .catch(() => {});
      }).catch(() => {});
    },
    handleDownload(row) {
      const link = document.createElement('a');
      link.href = window.URL.createObjectURL(new Blob([row]));
      link.download = row.filename;
      link.click();
      window.URL.revokeObjectURL(link.href);
      
      downloadBackupApi(row.filename).then((blob) => {
        const downloadLink = document.createElement('a');
        downloadLink.href = window.URL.createObjectURL(blob);
        downloadLink.download = row.filename;
        downloadLink.click();
        window.URL.revokeObjectURL(downloadLink.href);
        this.$message.success('下载已开始');
      }).catch(() => {});
    },
    handleRestore(row) {
      this.selectedBackup = row;
      this.restoreDialogVisible = true;
    },
    confirmRestore() {
      this.restoreLoading = true;
      restoreDatabaseApi({ filename: this.selectedBackup.filename })
        .then(() => {
          this.$message.success('恢复成功');
          this.restoreDialogVisible = false;
          this.restoreLoading = false;
          this.getList();
        })
        .catch(() => {
          this.restoreLoading = false;
        });
    },
    handleDelete(row) {
      this.$confirm(`确定要删除备份文件 ${row.filename} 吗？`, '提示', {
        type: 'warning',
      }).then(() => {
        this.$message.info('删除功能开发中');
      }).catch(() => {});
    },
    formatSize(bytes) {
      if (!bytes || bytes === 0) return '0 B';
      const k = 1024;
      const sizes = ['B', 'KB', 'MB', 'GB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i];
    },
  },
};
</script>

<style scoped>
.flex-between {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.mt14 {
  margin-top: 14px;
}
.mb20 {
  margin-bottom: 20px;
}
.empty-tip {
  padding: 40px;
  color: #909399;
}
.text-center {
  text-align: center;
}
</style>
