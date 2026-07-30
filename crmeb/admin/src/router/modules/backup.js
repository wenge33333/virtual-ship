// 系统管理 - 数据备份路由
export default {
  path: '/system',
  component: () => import('@/layout'),
  redirect: '/system/backup',
  name: 'System',
  meta: {
    title: '系统管理',
    icon: 'el-icon-setting',
    roles: ['admin'],
  },
  children: [
    {
      path: 'backup',
      component: () => import('@/views/system/backup/index'),
      name: 'BackupManagement',
      meta: {
        title: '数据备份',
        icon: 'el-icon-connection',
        roles: ['admin'],
      },
    },
  ],
};
