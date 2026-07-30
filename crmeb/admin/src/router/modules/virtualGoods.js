import Layout from '@/layout';

const virtualGoodsRouter = {
  path: '/virtualGoods',
  component: Layout,
  redirect: '/virtualGoods/goodsConfig',
  name: 'VirtualGoods',
  meta: { title: '虚拟商品', icon: 'el-icon-goods' },
  children: [
    {
      path: 'goodsConfig',
      component: () => import('@/views/virtualGoods/goodsConfig/index'),
      name: 'VirtualGoodsConfig',
      meta: { title: '商品配置' },
    },
    {
      path: 'cardManage',
      component: () => import('@/views/virtualGoods/cardManage/index'),
      name: 'VirtualCardManage',
      meta: { title: '卡密管理' },
    },
    {
      path: 'deliveryLog',
      component: () => import('@/views/virtualGoods/deliveryLog/index'),
      name: 'VirtualDeliveryLog',
      meta: { title: '发货日志' },
    },
    {
      path: 'platformConfig',
      component: () => import('@/views/virtualGoods/platformConfig/index'),
      name: 'VirtualPlatformConfig',
      meta: { title: '平台配置' },
    },
  ],
};

export default virtualGoodsRouter;
