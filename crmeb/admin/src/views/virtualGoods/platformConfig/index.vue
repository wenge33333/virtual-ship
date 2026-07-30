<template>
  <div class="divBox relative">
    <el-card class="box-card" shadow="never">
      <div class="clearfix flex-between">
        <el-button type="primary" size="small" icon="el-icon-plus" @click="handleCreate">添加平台配置</el-button>
        <el-button type="primary" size="small" icon="el-icon-download" @click="handleExport">导出 CSV</el-button>
      </div>
      <el-table v-loading="loading" :data="configList" style="width:100%" size="small" class="mt14">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="platform_code" label="平台代码" width="120" />
        <el-table-column prop="platform_name" label="平台名称" width="130" />
        <el-table-column prop="app_key" label="AppKey" min-width="160" :show-overflow-tooltip="true" />
        <el-table-column label="拉取间隔" width="100">
          <template slot-scope="scope">
            <span>{{ scope.row.pull_interval }}秒</span>
          </template>
        </el-table-column>
        <el-table-column label="上次拉取" width="160">
          <template slot-scope="scope">
            <span v-if="scope.row.last_pull_time">{{ scope.row.last_pull_time }}</span>
            <span v-else class="text-gray">未拉取</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template slot-scope="scope">
            <el-switch
              v-model="scope.row.status"
              :active-value="1"
              :inactive-value="0"
              active-text="启用"
              inactive-text="停用"
              @change="handleToggle(scope.row)"
            />
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
    </el-card>

    <!-- 编辑/新增弹窗 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="520px" :close-on-click-modal="false">
      <el-form ref="configForm" :model="formData" :rules="formRules" label-width="100px" size="small">
        <el-form-item label="平台代码" prop="platform_code">
          <el-select v-model="formData.platform_code" :disabled="isEdit" style="width:100%" placeholder="选择平台">
            <el-option value="taobao" label="淘宝/天猫" />
            <el-option value="xianyu" label="闲鱼" />
            <el-option value="douyin" label="抖音" />
            <el-option value="pinduoduo" label="拼多多" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台名称" prop="platform_name">
          <el-input v-model="formData.platform_name" placeholder="自定义显示名称" />
        </el-form-item>
        <el-form-item label="AppKey" prop="app_key">
          <el-input v-model="formData.app_key" placeholder="平台分配的AppKey" />
        </el-form-item>
        <el-form-item label="AppSecret" prop="app_secret">
          <el-input v-model="formData.app_secret" type="password" placeholder="平台分配的AppSecret" show-password />
        </el-form-item>
        <el-form-item label="SessionKey">
          <el-input v-model="formData.session_key" placeholder="OAuth授权后的SessionKey（可选）" />
        </el-form-item>
        <el-form-item label="拉取间隔(秒)" prop="pull_interval">
          <el-input-number v-model="formData.pull_interval" :min="10" :max="3600" style="width:100%" />
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
import { platformConfigListApi, platformConfigSaveApi, platformConfigToggleApi, platformConfigDeleteApi, exportPlatformConfigsApi } from '@/api/virtualGoods';

export default {
  data() {
    return {
      configList: [],
      loading: false,
      submitLoading: false,
      dialogVisible: false,
      isEdit: false,
      formData: {
        platform_code: '',
        platform_name: '',
        app_key: '',
        app_secret: '',
        session_key: '',
        pull_interval: 60,
      },
      formRules: {
        platform_code: [{ required: true, message: '请选择平台', trigger: 'change' }],
        platform_name: [{ required: true, message: '请输入平台名称', trigger: 'blur' }],
        app_key: [{ required: true, message: '请输入AppKey', trigger: 'blur' }],
        app_secret: [{ required: true, message: '请输入AppSecret', trigger: 'blur' }],
        pull_interval: [{ required: true, message: '请设置拉取间隔', trigger: 'blur' }],
      },
    };
  },
  mounted() {
    this.getList();
  },
  methods: {
    getList() {
      this.loading = true;
      platformConfigListApi()
        .then((res) => {
          this.configList = Array.isArray(res) ? res : [];
          this.loading = false;
        })
        .catch(() => {
          this.loading = false;
        });
    },
    resetForm() {
      this.formData = {
        platform_code: '',
        platform_name: '',
        app_key: '',
        app_secret: '',
        session_key: '',
        pull_interval: 60,
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
        platform_code: row.platform_code,
        platform_name: row.platform_name,
        app_key: row.app_key,
        app_secret: '',
        session_key: '',
        pull_interval: row.pull_interval,
      };
      this.dialogVisible = true;
    },
    handleSubmit() {
      this.$refs.configForm.validate((valid) => {
        if (!valid) return;
        this.submitLoading = true;
        platformConfigSaveApi({
          ...this.formData,
          status: this.isEdit ? undefined : 1,
        })
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
    handleToggle(row) {
      platformConfigToggleApi({
        platform_code: row.platform_code,
        status: row.status,
      })
        .then(() => {
          this.$message.success(row.status === 1 ? '已启用' : '已停用');
        })
        .catch(() => {
          this.getList();
        });
    },
    handleExport() {
      exportPlatformConfigsApi().then((blob) => {
        const link = document.createElement('a');
        link.href = window.URL.createObjectURL(new Blob([blob]));
        link.download = `platform_configs_${new Date().getTime()}.csv`;
        link.click();
        window.URL.revokeObjectURL(link.href);
        this.$message.success('导出成功');
      }).catch(() => {});
    },
  },
};
</script>

<style scoped>
.mt14 {
  margin-top: 14px;
}
.text-gray {
  color: #c0c4cc;
}
.flex-between {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
