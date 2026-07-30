// +----------------------------------------------------------------------
// +----------------------------------------------------------------------
// +----------------------------------------------------------------------
// +----------------------------------------------------------------------
// +----------------------------------------------------------------------

import request from '@/utils/request';

/**
 * 获取版权信息
 */
export function copyrightInfoApi() {
  return request({
    url: '/admin/copyright/get/info',
    method: 'get',
  });
}

/**
 * 保存版权信息
 */
export function saveCrmebCopyRight(data) {
  return request({
    url: '/admin/copyright/update/company/info',
    method: 'post',
    data,
  });
}

/**
 * @description 账号登录检测
 */
export function accountDetectionApi(data) {
  return request({
    url: '/admin/login/account/detection',
    method: 'post',
    data,
  });
}
