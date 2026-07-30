import request from '@/utils/request';

// ========== 虚拟商品配置 ==========
export function virtualGoodsConfigListApi(params) {
  return request({ url: '/admin/virtual/goods/config/list', method: 'get', params });
}

export function virtualGoodsConfigGetApi(productId) {
  return request({ url: `/admin/virtual/goods/config/${productId}`, method: 'get' });
}

export function virtualGoodsConfigSaveApi(data) {
  return request({ url: '/admin/virtual/goods/config/save', method: 'post', data });
}

export function virtualGoodsConfigUpdateApi(data) {
  return request({ url: '/admin/virtual/goods/config/update', method: 'post', data });
}

export function virtualGoodsConfigDeleteApi(productId) {
  return request({ url: `/admin/virtual/goods/config/${productId}`, method: 'delete' });
}

export function exportConfigsApi() {
  return request({ url: '/admin/virtual/goods/config/export', method: 'get', responseType: 'blob' });
}

// ========== 卡密管理 ==========
export function cardListApi(params) {
  return request({ url: '/admin/virtual/card/list', method: 'get', params });
}

export function cardStockApi(productId) {
  return request({ url: `/admin/virtual/card/stock/${productId}`, method: 'get' });
}

export function cardImportApi(data) {
  return request({
    url: '/admin/virtual/card/import',
    method: 'post',
    data,
    headers: { 'Content-Type': 'multipart/form-data' },
  });
}

export function deleteCardsApi(ids) {
  return request({ url: '/admin/virtual/cards/batch', method: 'delete', data: { ids } });
}

export function exportCardsApi(params) {
  return request({ url: '/admin/virtual/cards/export', method: 'get', params, responseType: 'blob' });
}

// ========== 发货管理 ==========
export function deliveryAutoApi(data) {
  return request({ url: '/admin/virtual/delivery/auto', method: 'post', data });
}

export function deliveryManualApi(data) {
  return request({ url: '/admin/virtual/delivery/manual', method: 'post', data });
}

export function deliveryLogsApi(params) {
  return request({ url: '/admin/virtual/delivery/logs', method: 'get', params });
}

export function deliveryRetryApi(data) {
  return request({ url: '/admin/virtual/delivery/retry', method: 'post', data });
}

export function deleteDeliveryLogsApi(ids) {
  return request({ url: '/admin/virtual/delivery/logs/batch', method: 'delete', data: { ids } });
}

export function exportDeliveryLogsApi(params) {
  return request({ url: '/admin/virtual/delivery/logs/export', method: 'get', params, responseType: 'blob' });
}

// ========== 平台配置 ==========
export function platformConfigListApi() {
  return request({ url: '/admin/virtual/platform/config/list', method: 'get' });
}

export function platformConfigSaveApi(data) {
  return request({ url: '/admin/virtual/platform/config/save', method: 'post', data });
}

export function platformConfigToggleApi(data) {
  return request({ url: '/admin/virtual/platform/config/toggle', method: 'post', data });
}

export function platformConfigDeleteApi(platformCode) {
  return request({ url: `/admin/virtual/platform/config/${platformCode}`, method: 'delete' });
}

export function exportPlatformConfigsApi() {
  return request({ url: '/admin/virtual/platform/config/export', method: 'get', responseType: 'blob' });
}

// ========== 数据备份管理 ==========
export function backupDatabaseApi() {
  return request({ url: '/admin/backup/database', method: 'post' });
}

export function listBackupsApi() {
  return request({ url: '/admin/backup/list', method: 'get' });
}

export function restoreDatabaseApi(data) {
  return request({ url: '/admin/backup/restore', method: 'post', data });
}

export function downloadBackupApi(filename) {
  return request({ url: `/admin/backup/download/${filename}`, method: 'get', responseType: 'blob' });
}
