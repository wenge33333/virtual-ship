// +----------------------------------------------------------------------
// +----------------------------------------------------------------------
// +----------------------------------------------------------------------
// +----------------------------------------------------------------------
// +----------------------------------------------------------------------

/**
 * 再次确定弹窗组件
 * @param title 标题
 * @returns {Promise<any>}
 */
export default function modalSure(title) {
  return new Promise((resolve, reject) => {
    this.$confirm(`确定${title || '永久删除该数据'}`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
      customClass: 'deleteConfirm',
    })
      .then(() => {
        resolve();
      })
      .catch(() => {
        reject();
        this.$message({
          type: 'info',
          message: '已取消',
        });
      });
  });
}
