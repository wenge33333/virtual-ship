# 需求实施计划

- [ ] 1. 项目初始化与品牌替换
  - [x] 1.1 将 CRMEB_Java 源码克隆到工作区并重命名项目
    - 克隆源码到 /workspace/crmeb
    - 修改 Maven 父 POM 的 groupId 和 artifactId
    - 修改子模块 artifactId

  - [x] 1.2 替换后端 Java 源码中的品牌标识
    - 批量移除所有 Java 文件的 CRMEB 版权注释头（参考需求1-3）
    - 将 CrmebConfig.java 重命名为 AppConfig.java，属性前缀 crmeb 改为 app（参考需求1-1）
    - 修改 application.yml 中的 version、water-mark 等品牌配置（参考需求1-1）
    - 修改 Constants.java 中品牌相关常量（参考需求1-5）

  - [x] 1.3 替换前端 Admin (Vue) 中的品牌标识
    - 替换 admin/public/static/ 下的 favicon、logo 图片资源（参考需求1-1）
    - 搜索替换 admin/src/ 下所有 Vue/JS 文件中的 crmeb 字符串（参考需求1-1）
    - 修改 admin/package.json 中的 name、description 字段（参考需求1-1）

  - [x] 1.4 替换前端 App (UniApp) 中的品牌标识
    - 修改 app/manifest.json 中的 name、description（参考需求1-1）
    - 修改 app/pages.json 中的 navigationBarTitleText（参考需求1-1）
    - 搜索替换 app/pages/ 下所有 vue 文件中的 crmeb 字符串（参考需求1-1）

  - [x] 1.5 移除升级/授权检测模块
    - 删除 admin/src/layout/upgrade/ 目录（参考需求1-6）
    - 删除 app/components/update/ 目录（参考需求1-6）
    - 删除 app/pages/users/app_update/ 目录（参考需求1-6）
    - 删除 CopyrightServiceImpl.java（参考需求1-5）
    - 移除 IndexController.java 中的 /index/get/version 接口（参考需求1-5）
    - 移除 ThemeController.java 中的 /version 接口（参考需求1-5）

- [x] 2. 数据库改造
  - [x] 2.1 修改数据库初始化脚本
    - 替换 Crmeb_v3.0.sql 中品牌相关的 INSERT 数据（参考需求1-4）
    - 移除 SQL 文件头部的 CRMEB 版权注释（参考需求1-3）

  - [x] 2.2 创建虚拟商品发货系统相关数据表
    - 创建 eb_virtual_goods_config（虚拟商品发货配置表）
    - 创建 eb_card_secret（卡密库存表）
    - 创建 eb_card_import_log（卡密导入记录表）
    - 创建 eb_delivery_log（发货日志表）
    - 创建 eb_platform_config（外部平台授权配置表）
    - 创建 eb_platform_order（外部平台订单映射表）
    - 创建 eb_api_credential（API开放接口认证表）
    - ALTER TABLE eb_store_product 新增 product_type 字段

- [x] 3. 检查点 - 确保项目结构和数据表就绪
  - 确保 Maven 项目可正常编译
  - 确保 SQL 脚本可正常执行
  - 如有问题请询问用户

- [ ] 4. 实现虚拟商品管理模块
  - [ ] 4.1 创建 VirtualGoodsConfig 实体和 DTO
    - 在 crmeb-common 创建对应 model、request、response 类（参考需求2-1~2-5）
    - 创建 VirtualGoodsConfigRequest（保存/更新请求）
    - 创建 VirtualGoodsConfigResponse（查询响应）

  - [ ] 4.2 实现 VirtualGoodsService 接口与实现类
    - 实现 saveConfig 方法：为商品配置虚拟发货规则（参考需求2-1、2-3）
    - 实现 getConfig 方法：按商品 ID 查询配置（参考需求2-2）
    - 实现 updateConfig 方法：更新发货配置（参考需求2-3）

  - [ ] 4.3 创建 Admin 端虚拟商品配置控制器
    - 创建 VirtualGoodsController，提供配置的 CRUD 接口
    - 接口纳入 Spring Security 权限控制

  - [ ] 4.4 编写虚拟商品管理单元测试
    - 测试 saveConfig 创建配置
    - 测试 updateConfig 更新配置
    - 测试 getConfig 查询配置

- [ ] 5. 实现卡密库存管理模块
  - [ ] 5.1 创建 CardSecret 实体和 DTO
    - 在 crmeb-common 创建 CardSecret 模型、CardImportResult、CardStockStat 等类
    - CardSecret 模型包含 cardPassword AES-256 加密/解密逻辑（参考需求9-1）

  - [ ] 5.2 实现 CardSecretService 接口与实现类
    - 实现 importCards 方法：解析 Excel/CSV 文件，批量导入卡密（参考需求3-1、3-2）
    - 实现 getNextAvailable 方法：使用 SELECT FOR UPDATE 行锁获取可用卡密（参考需求3-4）
    - 实现 getStockStat 方法：统计库存数据（参考需求3-3）
    - 实现 recycleCard 方法：退款后回收卡密（参考需求3-6）

  - [ ] 5.3 创建 Admin 端卡密管理控制器
    - 创建 CardSecretController：卡密导入、列表查询、库存统计接口
    - 卡密导入接口接收 MultipartFile，返回 CardImportResult

  - [ ] 5.4 编写卡密库存管理单元测试
    - 测试 importCards 批量导入与校验
    - 测试 getNextAvailable 获取可用卡密
    - 测试 recycleCard 卡密回收

- [ ] 6. 实现自动发货模块
  - [ ] 6.1 创建 DeliveryLog 实体和 DTO
    - 创建 DeliveryLog 模型、DeliveryResult、DeliveryLogSearchRequest 等类

  - [ ] 6.2 实现 DeliveryService 接口与实现类
    - 实现 autoDelivery 方法：支付后自动发货核心逻辑（参考需求4-1~4-5）
    - 发货流程包含：查询配置 → 获取卡密 → 更新订单 → 写入日志 → 发送站内通知
    - 使用 Redis 分布式锁 lock:delivery:{orderId} 确保幂等性（参考需求9-4）
    - 实现 manualDelivery 方法：管理后台手动补发货（参考需求8-2）
    - 实现 retryDelivery 方法：失败重试，最多 3 次（参考需求4-5）

  - [ ] 6.3 改造支付回调流程
    - 修改 CallbackServiceImpl.java，支付成功后调用 DeliveryService.autoDelivery（参考需求4-1）
    - 添加商品类型判断：productType > 0 时触发虚拟发货逻辑

  - [ ] 6.4 实现站内消息通知
    - 发货完成后向用户发送站内消息通知（参考需求4-3）
    - 库存不足时向运营发送站内告警（参考需求4-4）

  - [ ] 6.5 编写自动发货模块单元测试
    - 测试 autoDelivery 卡密类型发货
    - 测试 autoDelivery 固定内容类型发货
    - 测试库存不足场景处理
    - 测试发货幂等性

  - [ ] 6.6 创建 DeliveryRetryTask 定时任务
    - 每 5 分钟扫描状态为失败的发货日志，执行重试（参考需求4-5）

  - [ ] 6.7 创建 StockWarnTask 定时任务
    - 每 5 分钟扫描库存低于阈值的商品，发送告警（参考需求3-5）

- [ ] 7. 检查点 - 确保自动发货核心流程可运行
  - 确保支付回调 → 自动发货 → 站内通知链路完整
  - 确保并发发货幂等性
  - 如有问题请询问用户

- [ ] 8. 实现管理后台订单与发货管理
  - [ ] 8.1 创建 Admin 端发货管理控制器
    - 创建 DeliveryController：手动补发货、发货日志查询、重试发货接口（参考需求8-2、8-3）
    - 订单列表增加虚拟商品类型筛选和平台来源筛选（参考需求8-1）

  - [ ] 8.2 创建虚拟商品销售统计接口
    - 按商品统计销量、销售额、退款率（参考需求8-4）
    - 支持时间范围筛选

  - [ ] 8.3 创建 Admin 前端虚拟商品配置页面
    - 在 admin/src/views/store/ 下创建 virtualGoodsConfig 页面（参考需求2）
    - 在 admin/src/views/ 下创建 cardSecret 卡密管理页面（参考需求3）
    - 在 admin/src/views/ 下创建 deliveryLog 发货日志页面（参考需求8-3）
    - 注册路由和左侧菜单

  - [ ] 8.4 编写后台管理模块单元测试
    - 测试手动补发货接口
    - 测试发货日志查询

- [ ] 9. 实现用户端虚拟商品订单功能
  - [ ] 9.1 创建 Front 端虚拟商品订单控制器
    - 订单列表接口：区分虚拟商品与实物商品，增加虚拟标识（参考需求7-1）
    - 订单详情接口：展示已发放的卡密信息（参考需求7-2）
    - 卡密复制接口：提供完整的卡密解密查看（参考需求7-2）
    - 退款申请接口：判断卡密状态决定自动退款或人工审核（参考需求7-3~7-5）

  - [ ] 9.2 修改 Front 前端 App 页面
    - 改造 app/pages/order/order_details/ 展示卡密信息（参考需求7-2）
    - 增加"复制卡密"功能按钮（参考需求7-2）
    - 订单列表增加虚拟商品标识（参考需求7-1）

  - [ ] 9.3 编写用户端模块单元测试
    - 测试订单详情卡密展示
    - 测试退款自动回库

- [ ] 10. 实现外部平台订单对接
  - [ ] 10.1 创建 PlatformConfig 实体和 DTO
    - 创建 PlatformConfig 模型（含 AppSecret/SessionKey 的 AES 加解密）

  - [ ] 10.2 实现 PlatformSyncService 接口与实现类
    - 实现 syncTaobaoOrders：对接淘宝开放平台 TOP SDK（参考需求5-1）
    - 实现 syncXianyuOrders：对接闲鱼平台（通过淘宝开放平台）（参考需求5-2）
    - 实现 syncDouyinOrders：对接抖音开放平台（参考需求5-3）
    - 实现 syncPinduoduoOrders：对接拼多多开放平台（参考需求5-4）
    - 实现 bindProductMapping：外部商品与本地商品的映射绑定（参考需求5-6）
    - 每个平台同步方法包含：鉴权 → 拉取订单 → 创建本地订单 → 触发发货的完整链路

  - [ ] 10.3 创建 Admin 端平台配置控制器
    - 创建 PlatformController：平台授权配置 CRUD 接口（参考需求5-5）
    - 平台授权信息（AppKey/AppSecret/SessionKey）加密存储
    - 支持启停用和拉取间隔设置

  - [ ] 10.4 创建平台同步定时任务
    - 创建 TaobaoOrderSyncTask（参考需求5-1）
    - 创建 XianyuOrderSyncTask（参考需求5-2）
    - 创建 DouyinOrderSyncTask（参考需求5-3）
    - 创建 PinduoduoOrderSyncTask（参考需求5-4）

  - [ ] 10.5 创建 Admin 前端平台配置页面
    - 在 admin/src/views/ 下创建 platformConfig 平台授权配置页面
    - 在 admin/src/views/ 下创建 platformOrder 平台订单管理页面
    - 注册路由和左侧菜单

  - [ ] 10.6 编写平台对接模块单元测试
    - 测试平台配置的加密存储
    - 测试订单拉取和映射逻辑

- [ ] 11. 实现开放 API 接口
  - [ ] 11.1 创建 ApiCredential 实体和 HMAC 签名工具类
    - 创建 ApiCredential 模型（参考需求6-1）
    - 创建 HmacSignatureUtil：实现 HMAC-SHA256 签名生成与验证（参考需求6-1）

  - [ ] 11.2 实现 API 认证拦截器
    - 创建 ApiAuthInterceptor：从请求头解析 X-Api-Key、X-Timestamp、X-Signature
    - 验证时间戳有效期（5 分钟内）
    - 验证 HMAC 签名
    - 实现速率限制（参考需求6-8）

  - [ ] 11.3 创建 OpenApi 控制器
    - POST /api/open/v1/delivery/create：创建发货订单（参考需求6-3）
    - GET /api/open/v1/delivery/status/{orderNo}：查询发货状态（参考需求6-4）
    - GET /api/open/v1/goods/list：查询虚拟商品列表（参考需求6-2）
    - GET /api/open/v1/goods/{productId}：查询商品详情（参考需求6-2）
    - GET /api/open/v1/goods/{productId}/stock：查询库存余量（参考需求6-5）
    - 接口文档通过 Swagger/OpenAPI 自动生成（参考需求6-6）
    - 所有接口记录请求日志（参考需求6-7）

  - [ ] 11.4 编写开放 API 模块单元测试
    - 测试 HMAC 签名验证
    - 测试速率限制
    - 测试各接口的请求/响应

- [ ] 12. 检查点 - 最终验证
  - 确保所有模块可正常编译
  - 确保所有接口逻辑完整
  - 如有问题请询问用户
