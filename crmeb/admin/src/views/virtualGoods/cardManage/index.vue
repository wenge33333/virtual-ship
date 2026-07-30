<template>
  <div class="divBox relative">
    <el-card class="box-card" shadow="never">
      <div class="clearfix">
        <div class="filter-container flex-between">
          <el-form inline size="small">
            <el-form-item label="商品ID：">
              <el-input v-model="listPram.product_id" placeholder="请输入商品ID" style="width:200px" clearable @keyup.enter.native="handleSearch" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="el-icon-search" @click="handleSearch">查询库存</el-button>
              <el-button type="success" icon="el-icon-upload2" @click="handleImportDialog">导入卡密</el-button>
              <el-button type="danger" icon="el-icon-delete" @click="handleBatchDelete" :disabled="selectedCards.length===0">批量删除</el-button>
            </el-form-item>
          </el-form>
          <el-button type="primary" icon="el-icon-download" @click="handleExport">导出 CSV</el-button>
        </div>
      </div>

      <!-- 库存概览 -->
      <div v-if="stockStat" class="stock-overview mt14">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-label">总库存</div>
              <div class="stat-value">{{ stockStat.total }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card stat-remaining">
              <div class="stat-label">剩余可用</div>
              <div class="stat-value">{{ stockStat.remaining }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card stat-sold">
              <div class="stat-label">已售出</div>
              <div class="stat-value">{{ stockStat.sold }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card stat-warn">
              <div class="stat-label">告警阈值</div>
              <div class="stat-value">{{ stockStat.warn_threshold || '-' }}</div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 卡密列表 -->
      <el-table v-loading="loading" :data="cardData.list" style="width:100%" size="small" class="mt14" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="product_id" label="商品ID" width="80" />
        <el-table-column prop="card_number" label="卡号" min-width="150" :show-overflow-tooltip="true" />
        <el-table-column label="卡密" min-width="180">
          <template slot-scope="scope">
            <span v-if="!scope.row.isRevealed">
              {{ scope.row.card_password.substring(0, 4) }}****
              <el-button type="text" size="mini" @click="toggleReveal(scope.row)">查看</el-button>
            </span>
            <span v-else>
              {{ scope.row.card_password }}
              <el-button type="text" size="mini" @click="toggleReveal(scope.row)">隐藏</el-button>
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="face_value" label="面值" width="80" />
        <el-table-column label="状态" width="80">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.status === 0" type="success" size="small">未售</el-tag>
            <el-tag v-else-if="scope.row.status === 1" type="info" size="small">已售</el-tag>
            <el-tag v-else type="warning" size="small">已过期</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="过期时间" width="160">
          <template slot-scope="scope">
            <span v-if="scope.row.expire_time">{{ scope.row.expire_time }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="order_id" label="关联订单" width="150" :show-overflow-tooltip="true" />
        <el-table-column label="创建时间" width="160">
          <template slot-scope="scope">
            <span>{{ scope.row.create_time }}</span>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-if="cardData.total > 0"
        class="mt14"
        :current-page="listPram.page"
        :page-sizes="constants.page.limit"
        :layout="constants.page.layout"
        :total="cardData.total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        background
      />
    </el-card>

    <!-- 导入弹窗 -->
    <el-dialog title="导入卡密" :visible.sync="importDialogVisible" width="480px" :close-on-click-modal="false">
      <el-form ref="importForm" :model="importForm" :rules="importRules" label-width="100px" size="small">
        <el-form-item label="商品ID" prop="product_id">
          <el-input-number v-model="importForm.product_id" :min="1" style="width:100%" />
        </el-form-item>
        <el-form-item label="卡密文件" prop="file">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :limit="1"
            accept=".csv"
            action="#"
            :on-change="handleFileChange"
            :file-list="importForm.fileList"
          >
            <el-button size="small" type="primary" icon="el-icon-upload2">选择CSV文件</el-button>
            <div slot="tip" class="el-upload__tip">CSV格式：card_number,card_password,face_value,expire_time</div>
          </el-upload>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="importDialogVisible = false" size="small">取消</el-button>
        <el-button type="primary" @click="handleImport" size="small" :loading="importLoading">开始导入</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { cardListApi, cardStockApi, cardImportApi, deleteCardsApi, exportCardsApi } from '@/api/virtualGoods';

export default {
  data() {
    return {
      constants: this.$constants,
      listPram: {
        product_id: null,
        page: 1,
        limit: this.$constants.page.limit[0],
      },
      cardData: { list: [], total: 0 },
      stockStat: null,
      loading: false,
      importDialogVisible: false,
      importLoading: false,
      selectedCards: [],
      importForm: {
        product_id: undefined,
        fileList: [],
      },
      importRules: {
        product_id: [{ required: true, message: '请输入商品 ID', trigger: 'blur' }],
      },
    };
  },
  methods: {
    handleSelectionChange(val) {
      this.selectedCards = val;
    },
    handleSearch() {
      const pid = parseInt(this.listPram.product_id);
      if (!pid) {
        this.$message.warning('请输入商品 ID');
        return;
      }
      this.listPram.page = 1;
      this.loading = true;
      Promise.all([
        cardStockApi(pid).catch(() => null),
        cardListApi({ product_id: pid, page: 1, limit: this.listPram.limit }).catch(() => ({ list: [], total: 0 })),
      ]).then(([stockRes, listRes]) => {
        this.stockStat = stockRes;
        this.cardData = {
          list: (listRes.list || []).map((item) => ({ ...item, isRevealed: false })),
          total: listRes.total || 0,
        };
        this.loading = false;
      });
    },
    handleBatchDelete() {
      if (this.selectedCards.length === 0) {
        this.$message.warning('请选择要删除的卡密');
        return;
      }
      this.$confirm(`确认删除选中的 ${this.selectedCards.length} 条卡密吗？此操作不可恢复`, '提示', { type: 'warning' }).then(() => {
        const ids = this.selectedCards.map(item => item.id);
        deleteCardsApi(ids)
          .then((res) => {
            this.$message.success(`成功删除 ${res.deleted} 条卡密`);
            this.handleSearch();
          })
          .catch(() => {});
      }).catch(() => {});
    },
    handleExport() {
      const pid = parseInt(this.listPram.product_id);
      if (!pid) {
        this.$message.warning('请先输入商品 ID');
        return;
      }
      exportCardsApi({ product_id: pid }).then((blob) => {
        const link = document.createElement('a');
        link.href = window.URL.createObjectURL(new Blob([blob]));
        link.download = `cards_${pid}_${new Date().getTime()}.csv`;
        link.click();
        window.URL.revokeObjectURL(link.href);
        this.$message.success('导出成功');
      }).catch(() => {});
    },
    getList() {
      const pid = parseInt(this.listPram.product_id);
      if (!pid) return;
      this.loading = true;
      cardListApi({ product_id: pid, page: this.listPram.page, limit: this.listPram.limit })
        .then((data) => {
          this.cardData = {
            list: (data.list || []).map((item) => ({ ...item, isRevealed: false })),
            total: data.total || 0,
          };
          this.loading = false;
        })
        .catch(() => {
          this.cardData = { list: [], total: 0 };
          this.loading = false;
        });
    },
    handleSizeChange(val) {
      this.listPram.limit = val;
      this.getList();
    },
    handleCurrentChange(val) {
      this.listPram.page = val;
      this.getList();
    },
    toggleReveal(row) {
      this.$set(row, 'isRevealed', !row.isRevealed);
    },
    handleImportDialog() {
      this.importForm = { product_id: undefined, fileList: [] };
      this.importDialogVisible = true;
    },
    handleFileChange(file, fileList) {
      this.importForm.fileList = fileList;
    },
    handleImport() {
      this.$refs.importForm.validate((valid) => {
        if (!valid) return;
        if (this.importForm.fileList.length === 0) {
          this.$message.warning('请选择CSV文件');
          return;
        }
        this.importLoading = true;
        const formData = new FormData();
        formData.append('product_id', this.importForm.product_id);
        formData.append('operator_id', 1);
        formData.append('file', this.importForm.fileList[0].raw);
        cardImportApi(formData)
          .then((res) => {
            this.$message.success(`导入完成: 成功${res.success_count}条, 失败${res.fail_count}条`);
            this.importDialogVisible = false;
            this.importLoading = false;
            if (this.listPram.product_id) this.handleSearch();
          })
          .catch(() => {
            this.importLoading = false;
          });
      });
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
.stock-overview {
  margin-bottom: 10px;
}
.stat-card {
  text-align: center;
}
.stat-label {
  font-size: 13px;
  color: #909399;
}
.stat-value {
  font-size: 28px;
  font-weight: bold;
  margin-top: 8px;
  color: #303133;
}
.stat-remaining .stat-value {
  color: #67c23a;
}
.stat-sold .stat-value {
  color: #e6a23c;
}
.stat-warn .stat-value {
  color: #f56c6c;
}
</style>
