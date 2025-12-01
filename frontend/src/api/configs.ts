import request from '@/utils/request';
import { APIResponse } from './rts';

// 系统配置类型
export interface SystemConfig {
  proxy_list: string; // JSON 字符串
  client_id_list: string; // JSON 字符串
  auto_refresh_enabled: boolean;
  auto_refresh_interval: number; // 分钟
}

// 环境变量配置（只读）
export interface EnvConfig {
  API_PREFIX: string;
  API_SECRET: string;
  ADMIN_PREFIX: string;
  ADMIN_USERNAME: string;
  ADMIN_PASSWORD: string;
}

// 系统配置响应
export interface SystemConfigsResponse {
  configs: Record<string, string>;
  env_configs: EnvConfig;
}

// 配置管理 API（全部改成 POST + JSON Body）
export const configsApi = {
  // 获取系统配置
  getSystemConfigs: (): Promise<APIResponse<SystemConfigsResponse>> => {
    return request.post('/configs/get-system', {});
  },

  // 保存系统配置
  saveSystemConfigs: (configs: Record<string, string>): Promise<APIResponse<null>> => {
    return request.post('/configs/save-system', { configs });
  },

  // 获取代理列表
  getProxyList: (): Promise<APIResponse<string[]>> => {
    return request.post('/configs/get-proxy-list', {});
  },

  // 获取 Client ID 列表
  getClientIDList: (): Promise<APIResponse<string[]>> => {
    return request.post('/configs/get-clientid-list', {});
  },
};
