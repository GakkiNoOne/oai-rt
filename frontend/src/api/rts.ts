import request from '@/utils/request';

// RT 数据类型
export interface RT {
  id: number;
  biz_id: string;
  user_name?: string;
  email?: string;
  type?: string;
  rt: string;
  at?: string;
  proxy?: string;
  client_id?: string;
  tag?: string;
  enabled: boolean;
  last_rt?: string;
  refresh_result?: string;
  user_info?: string;
  account_info?: string;
  last_refresh_time?: string;
  memo?: string;
  create_time: string;
  update_time: string;
}

// 创建 RT 请求
export interface CreateRTRequest {
  biz_id: string;
  rt_token: string;
  proxy?: string;
  tag?: string;
  enabled: boolean;
  memo?: string;
}

// 更新 RT 请求
export interface UpdateRTRequest {
  id: number;
  updates: {
    biz_id?: string;
    proxy?: string;
    tag?: string;
    enabled?: boolean;
    memo?: string;
  };
}

// 列表查询参数
export interface ListRTParams {
  page: number;
  page_size: number;
  biz_id?: string;
  tag?: string;
  email?: string;
  type?: string;
  enabled?: boolean;
  create_date?: string;
}

// 列表响应
export interface ListRTResponse {
  items: RT[];
  total: number;
  page: number;
  page_size: number;
}

// 批量操作结果
export interface BatchResult {
  total_count: number;
  success_count: number;
  fail_count: number;
  results?: Array<{
    rt_name: string;
    success: boolean;
    message: string;
  }>;
}

// API 响应格式
export interface APIResponse<T = any> {
  success: boolean;
  msg: string;
  data: T;
}

// RT 管理 API（全部改成 POST + JSON Body）
export const rtsApi = {
  // 获取 RT 列表
  list: (params: ListRTParams): Promise<APIResponse<ListRTResponse>> => {
    return request.post('/rts/list', params);
  },

  // 创建 RT
  create: (data: CreateRTRequest): Promise<APIResponse<RT>> => {
    return request.post('/rts/create', data);
  },

  // 更新 RT
  update: (id: number, updates: any): Promise<APIResponse<RT>> => {
    return request.post('/rts/update', { id, updates });
  },

  // 删除 RT
  delete: (id: number): Promise<APIResponse<null>> => {
    return request.post('/rts/delete', { id });
  },

  // 批量删除
  batchDelete: (ids: number[]): Promise<APIResponse<BatchResult>> => {
    return request.post('/rts/batch-delete', { ids });
  },

  // 批量刷新（测活）- 超时时间3600秒
  batchRefresh: (ids: number[]): Promise<APIResponse<BatchResult>> => {
    return request.post('/rts/batch-refresh', { ids }, { timeout: 3600000 });
  },

  // 单个刷新
  refresh: (id: number, refreshUserInfo: boolean = false, refreshAccountInfo: boolean = false): Promise<APIResponse<RT>> => {
    return request.post('/rts/refresh', { 
      id, 
      refresh_user_info: refreshUserInfo,
      refresh_account_info: refreshAccountInfo 
    });
  },

  // 批量导入
  batchCreate: (batchName: string, tag: string, rtTokens: string[]): Promise<APIResponse<BatchResult>> => {
    return request.post('/rts/batch-import', {
      batch_name: batchName,
      tag: tag,
      rt_tokens: rtTokens,
    });
  },

  // 刷新用户信息
  refreshUserInfo: (id: number): Promise<APIResponse<RT>> => {
    return request.post('/rts/refresh-user-info', { id });
  },

  // 刷新账号信息
  refreshAccountInfo: (id: number): Promise<APIResponse<RT>> => {
    return request.post('/rts/refresh-account-info', { id });
  },
};
