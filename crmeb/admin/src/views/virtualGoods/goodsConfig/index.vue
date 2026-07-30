<template>
  <div class="divBox relative">
    <el-card class="box-card" shadow="never">
      <div class="clearfix">
        <div class="filter-container">
          <el-form inline size="small">
            <el-form-item label="商品ID：">
              <el-input v-model="listPram.product_id" placeholder="请输入商品ID" style="width:200px" clearable />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="el-icon-search" @click="handleSearch">查询</el-button>
              <el-button icon="el-icon-refresh" @click="handleReset">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
        <div class="mt10 flex-between">
          <el-button type="primary" size="small" icon="el-icon-plus" @click="handleCreate">添加配置</el-button>
          <el-button type="success" size="small" icon="el-icon-download" @click="handleExport">导出 CSV</el-button>
        </div>
      </div>
      <el-table v-loading="loading" :data="listData.list" style="width:100%" size="small" class="mt14">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="product_id" label="商品ID" width="80" />
        <el-table-column label="发货类型" width="120">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.delivery_type === 1" type="success" size="small">自动发卡密</el-tag>
            <el-tag v-else-if="scope.row.delivery_type === 2" type="warning" size="small">手动处理</el-tag>
            <el-tag v-else type="info" size="small">固定内容</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="delivery_content" label="发货内容" min-width="150" :show-overflow-tooltip="true">
          <template slot-scope="scope">
            <span v-if="scope.row.delivery_content && scope.row.delivery_content.length > 60">{{ scope.row.delivery_content.substring(0, 60) }}...</span>
            <span v-else>{{ scope.row.delivery_content }}</span>
          </template>
        </el-table-column>
        <el-table-column label="选卡规则" width="90">
          <template slot-scope="scope">
            <span>{{ scope.row.pick_rule === 1 ? '顺序' : '随机' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="stock_warn_threshold" label="告警阈值" width="80" />
        <el-table-column label="显示库存" width="80">
          <template slot-scope="scope">
            <span>{{ scope.row.is_show_stock ? '是' : '否' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="160">
          <template slot-scope="scope">
            <span>{{ scope.row.create_time }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template slot-scope="scope">
            <a @click="handleEdit(scope.row)">编辑</a>
            <el-divider direction="vertical" />
            <a @click="handleDelete(scope.row)" style="color:#F56C6C">删除</a>
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

    <!-- 编辑/新增弹窗 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="560px" :close-on-click-modal="false">
      <el-form ref="configForm" :model="formData" :rules="formRules" label-width="110px" size="small">
        <el-form-item label="商品ID" prop="product_id">
          <el-input-number v-model="formData.product_id" :min="1" :disabled="isEdit" style="width:100%" />
        </el-form-item>
        <el-form-item label="发货类型" prop="delivery_type">
          <el-select v-model="formData.delivery_type" style="width:100%">
            <el-option :value="1" label="自动发卡密" />
            <el-option :value="2" label="手动处理" />
            <el-option :value="3" label="固定内容" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="formData.delivery_type === 3" label="发货固定内容" prop="delivery_content">
          <el-input v-model="formData.delivery_content" type="textarea" :rows="3" placeholder="固定发货内容" />
        </el-form-item>
        <el-form-item label="选卡规则" prop="pick_rule">
          <el-select v-model="formData.pick_rule" style="width:100%">
            <el-option :value="1" label="顺序选取" />
            <el-option :value="2" label="随机选取" />
          </el-select>
        </el-form-item>
        <el-form-item label="库存告警阈值" prop="stock_warn_threshold">
          <el-input-number v-model="formData.stock_warn_threshold" :min="0" style="width:100%" />
        </el-form-item>
        <el-form-item label="是否显示库存" prop="is_show_stock">
          <el-switch v-model="formData.is_show_stock" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false" size="small">取消</el-button>
        <el-button type="primary" @click="handleSubmit" size="small" :loading="submitLoading">确定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import {
  virtualGoodsConfigListApi,
  virtualGoodsConfigGetApi,
  virtualGoodsConfigSaveApi,
  virtualGoodsConfigUpdateApi,
  virtualGoodsConfigDeleteApi,
  exportConfigsApi,
} from '@/api/virtualGoods';

export default {
  data() {
    return {
      constants: this.$constants,
      listPram: {
        product_id: null,
        page: 1,
        limit: this.$constants.page.limit[0],
      },
      listData: { list: [], total: 0 },
      loading: false,
      submitLoading: false,
      dialogVisible: false,
      isEdit: false,
      formData: {
        id: 0,
        product_id: undefined,
        delivery_type: 1,
        delivery_content: '',
        pick_rule: 1,
        stock_warn_threshold: 10,
        is_show_stock: 1,
      },
      formRules: {
        product_id: [{ required: true, message: '请输入商品ID', trigger: 'blur' }],
        delivery_type: [{ required: true, message: '请选择发货类型', trigger: 'change' }],
        pick_rule: [{ required: true, message: '请选择选卡规则', trigger: 'change' }],
      },
    };
  },
  computed: {
    dialogTitle() {
      return this.isEdit ? '编辑虚拟商品配置' : '添加虚拟商品配置';
    },
  },
  mounted() {
    this.getList();
  },
  methods: {
    getList() {
      this.loading = true;
      const params = {
        page: this.listPram.page,
        limit: this.listPram.limit,
      };
      if (this.listPram.product_id) {
        params.product_id = this.listPram.product_id;
      }
      virtualGoodsConfigListApi(params)
        .then((data) => {
          this.listData = data;
          this.loading = false;
        })
        .catch(() => {
          this.loading = false;
        });
    },
    handleSearch() {
      this.listPram.page = 1;
      this.getList();
    },
    handleReset() {
      this.listPram.product_id = null;
      this.listPram.page = 1;
      this.getList();
    },
    handleSizeChange(val) {
      this.listPram.limit = val;
      this.getList();
    },
    handleCurrentChange(val) {
      this.listPram.page = val;
      this.getList();
    },
    resetForm() {
      this.formData = {
        id: 0,
        product_id: undefined,
        delivery_type: 1,
        delivery_content: '',
        pick_rule: 1,
        stock_warn_threshold: 10,
        is_show_stock: 1,
      };
      if (this.$refs.configForm) {
        this.$refs.configForm.resetFields();
      }
    },
    handleCreate() {
      this.isEdit = false;
      this.resetForm();
      this.dialogVisible = true;
    },
    handleEdit(row) {
      this.isEdit = true;
      this.formData = {
        id: row.id,
        product_id: row.product_id,
        delivery_type: row.delivery_type,
        delivery_content: row.delivery_content || '',
        pick_rule: row.pick_rule,
        stock_warn_threshold: row.stock_warn_threshold,
        is_show_stock: row.is_show_stock,
      };
      this.dialogVisible = true;
    },
    handleDelete(row) {
      this.$modalSure('确认删除该商品配置？').then(() => {
        virtualGoodsConfigDeleteApi(row.product_id)
          .then(() => {
            this.$message.success('删除成功');
            this.getList();
          })
          .catch(() => {
            this.$message.error('删除失败');
          });
      });
    },
    handleSubmit() {
      this.$refs.configForm.validate((valid) => {
        if (!valid) return;
        this.submitLoading = true;
        const api = this.isEdit ? virtualGoodsConfigUpdateApi : virtualGoodsConfigSaveApi;
        const data = { ...this.formData };
        if (data.delivery_type !== 3) {
          data.delivery_content = '';
        }
        api(data)
          .then(() => {
            this.$message.success('操作成功');
            this.dialogVisible = false;
            this.submitLoading = false;
            this.getList();
          })
          .catch(() => {
            this.submitLoading = false;
          });
      });
    },
    handleExport() {
      exportConfigsApi().then((blob) => {
        const link = document.createElement('a');
        link.href = window.URL.createObjectURL(new Blob([blob]));
        link.download = `virtual_goods_config_${new Date().getTime()}.csv`;
        link.click();
        window.URL.revokeObjectURL(link.href);
        this.$message.success('导出成功');
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
.mt10 {
  margin-top: 10px;
}
.mt14 {
  margin-top: 14px;
}
.flex-between {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
